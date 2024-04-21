package webp

/*
#include <stdlib.h>
#include <string.h>
#include <webp/encode.h>
#include <webp/mux.h>

int writeWebP(uint8_t*, size_t, struct WebPPicture*);

static WebPPicture *calloc_WebPPicture(void) {
	return calloc(sizeof(WebPPicture), 1);
}

static void free_WebPPicture(WebPPicture* webpPicture) {
	free(webpPicture);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"image"
	"io"
	"time"
	"unsafe"
)

// AnimationEncoder encodes multiple images into an animated WebP.
type AnimationEncoder struct {
	opts         C.WebPAnimEncoderOptions
	c            *C.WebPAnimEncoder
	duration     time.Duration
	WebPData     *C.WebPData
	WebPMux      *C.WebPMux
	WebPPictures []*C.WebPPicture
}

// NewAnimationEncoder initializes a new encoder.
func NewAnimationEncoder(width, height, kmin, kmax int, loopCount int) (*AnimationEncoder, error) {
	ae := &AnimationEncoder{}

	if C.WebPAnimEncoderOptionsInit(&ae.opts) == 0 {
		return nil, errors.New("failed to initialize animation encoder config")
	}
	ae.opts.kmin = C.int(kmin)
	ae.opts.kmax = C.int(kmax)
	ae.opts.anim_params.loop_count = C.int(loopCount)

	ae.c = C.WebPAnimEncoderNew(C.int(width), C.int(height), &ae.opts)
	if ae.c == nil {
		return nil, errors.New("failed to initialize animation encoder")
	}

	return ae, nil
}

// AddFrame adds a frame to the encoder.
func (ae *AnimationEncoder) AddFrame(img image.Image, duration time.Duration) error {
	pic := C.calloc_WebPPicture()
	if pic == nil {
		return errors.New("could not allocate webp picture")
	}
	//defer C.free_WebPPicture(pic)
	ae.WebPPictures = append(ae.WebPPictures, pic)
	if C.WebPPictureInit(pic) == 0 {
		return errors.New("could not initialize webp picture")
	}
	//defer C.WebPPictureFree(pic)

	pic.use_argb = 1

	pic.width = C.int(img.Bounds().Dx())
	pic.height = C.int(img.Bounds().Dy())

	pic.writer = C.WebPWriterFunction(C.writeWebP)

	switch p := img.(type) {
	case *RGBImage:
		C.WebPPictureImportRGB(pic, (*C.uint8_t)(&p.Pix[0]), C.int(p.Stride))
	case *image.RGBA:
		C.WebPPictureImportRGBA(pic, (*C.uint8_t)(&p.Pix[0]), C.int(p.Stride))
	case *image.NRGBA:
		C.WebPPictureImportRGBA(pic, (*C.uint8_t)(&p.Pix[0]), C.int(p.Stride))
	default:
		return errors.New("unsupported image type")
	}

	timestamp := C.int(ae.duration / time.Millisecond)
	ae.duration += duration

	if C.WebPAnimEncoderAdd(ae.c, pic, timestamp, nil) == 0 {
		return fmt.Errorf(
			"encoding error: %d - %s",
			int(pic.error_code),
			C.GoString(C.WebPAnimEncoderGetError(ae.c)),
		)
	}

	return nil
}

// Assemble assembles all frames into animated WebP.
func (ae *AnimationEncoder) Assemble(w io.Writer) error {
	// add final empty frame
	if C.WebPAnimEncoderAdd(ae.c, nil, C.int(ae.duration/time.Millisecond), nil) == 0 {
		return errors.New("couldn't add final empty frame")
	}

	ae.WebPData = &C.WebPData{}
	C.WebPDataInit(ae.WebPData)
	//defer C.WebPDataClear(data)

	if C.WebPAnimEncoderAssemble(ae.c, ae.WebPData) == 0 {
		return errors.New("error assembling animation")
	}

	size := int(ae.WebPData.size)
	out := make([]byte, size)
	n := copy(
		out, C.GoBytes(
			unsafe.Pointer(ae.WebPData.bytes),
			C.int(size),
		),
	)

	if n != size {
		return errors.New("error copying animation from C to Go")
	}

	_, err := w.Write(out)
	return err
}

// Close deletes the encoder and frees resources.
func (ae *AnimationEncoder) Close() {
	C.WebPDataClear(ae.WebPData)
	C.WebPMuxDelete(ae.WebPMux)
	for _, webpPicture := range ae.WebPPictures {
		C.free_WebPPicture(webpPicture)
		C.WebPPictureFree(webpPicture)
	}
	C.WebPAnimEncoderDelete(ae.c)
}
