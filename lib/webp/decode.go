package webp

/*
#include <stdlib.h>
#include <webp/decode.h>

static VP8StatusCode CheckDecBuffer(const WebPDecBuffer* const buffer);

*/
import "C"

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"unsafe"
)

const webPHeader = "RIFF????WEBPVP8"

// DecoderOptions specifies decoding options of WebP.
type DecoderOptions struct {
	BypassFiltering        bool            // If true, bypass filtering process
	NoFancyUpsampling      bool            // If true, do not fancy upsampling
	Crop                   image.Rectangle // Do cropping if image.Rectangle is not empty.
	Scale                  image.Rectangle // Do scaling if image.Rectangle is not empty.
	UseThreads             bool            // If true, use multi threads
	DitheringStrength      int             // Specify dithering strength [0=Off .. 100=full]
	Flip                   bool            // If true, flip output vertically
	AlphaDitheringStrength int             // Specify alpha dithering strength in [0..100]
}

// BitstreamFeatures represents the image properties which are retrieved from
// data stream.
type BitstreamFeatures struct {
	Width        int  // Image width in pixels.
	Height       int  // Image height in pixels.
	HasAlpha     bool // True if data stream contains a alpha channel.
	HasAnimation bool // True if data stream is an animation.
	Format       int  // 0 = undefined (/mixed), 1 = lossy, 2 = lossless.
}

// GetDecoderVersion returns decoder's version number, packed in hexadecimal.
// e.g; v0.4.2 is 0x000402
func GetDecoderVersion() (v int) {
	return int(C.WebPGetDecoderVersion())
}

// GetInfo retrieves width/height from data bytes.
func GetInfo(data []byte) (width, height int) {
	var w, h C.int
	C.WebPGetInfo((*C.uint8_t)(&data[0]), (C.size_t)(len(data)), &w, &h)
	return int(w), int(h)
}

// GetFeatures returns features as BitstreamFeatures retrieved from data stream.
func GetFeatures(data []byte) (f *BitstreamFeatures, err error) {
	var cf C.WebPBitstreamFeatures
	status := C.WebPGetFeatures((*C.uint8_t)(&data[0]), (C.size_t)(len(data)), &cf)

	if status != C.VP8_STATUS_OK {
		return nil, fmt.Errorf("unexpected VP8 status: %s", statusString(status))
	}

	f = &BitstreamFeatures{
		Width:        int(cf.width), // TODO: use Rectangle instead?
		Height:       int(cf.height),
		HasAlpha:     cf.has_alpha > 0,
		HasAnimation: cf.has_animation > 0,
		Format:       int(cf.format),
	}
	return
}

// DecodeYUVA decodes WebP image into YUV image with alpha channel, and returns
// it as *YUVAImage.
func DecodeYUVA(data []byte, options *DecoderOptions) (img *YUVAImage, err error) {
	config, err := initDecoderConfig(options)
	if err != nil {
		return nil, err
	}

	cDataPtr := (*C.uint8_t)(&data[0])
	cDataSize := (C.size_t)(len(data))

	// Retrieve WebP features from data stream.
	if status := C.WebPGetFeatures(cDataPtr, cDataSize, &config.input); status != C.VP8_STATUS_OK {
		return nil, fmt.Errorf("could not get features from the data stream, return %s", statusString(status))
	}

	outWidth, outHeight := calcOutputSize(config)
	buf := (*C.WebPYUVABuffer)(unsafe.Pointer(&config.output.u[0]))

	// Set up output configurations
	colorSpace := YUV420
	config.output.colorspace = C.MODE_YUV
	if config.input.has_alpha > 0 {
		colorSpace = YUV420A
		config.output.colorspace = C.MODE_YUVA
	}
	config.output.is_external_memory = 1

	// Allocate image and fill into buffer
	img = NewYUVAImage(image.Rect(0, 0, outWidth, outHeight), colorSpace)
	buf.y = (*C.uint8_t)(&img.Y[0])
	buf.u = (*C.uint8_t)(&img.Cb[0])
	buf.v = (*C.uint8_t)(&img.Cr[0])
	buf.a = nil
	buf.y_stride = C.int(img.YStride)
	buf.u_stride = C.int(img.CStride)
	buf.v_stride = C.int(img.CStride)
	buf.a_stride = 0
	buf.y_size = C.size_t(len(img.Y))
	buf.u_size = C.size_t(len(img.Cb))
	buf.v_size = C.size_t(len(img.Cr))
	buf.a_size = 0

	if config.input.has_alpha > 0 {
		buf.a = (*C.uint8_t)(&img.A[0])
		buf.a_stride = C.int(img.AStride)
		buf.a_size = C.size_t(len(img.A))
	}

	if status := C.WebPDecode(cDataPtr, cDataSize, config); status != C.VP8_STATUS_OK {
		return nil, fmt.Errorf("could not decode data stream, return %s", statusString(status))
	}

	return
}

// DecodeRGBA decodes WebP image into RGBA image and returns it as an *image.RGBA.
func DecodeRGBA(data []byte, options *DecoderOptions) (img *image.RGBA, err error) {
	config, err := initDecoderConfig(options)
	if err != nil {
		return nil, err
	}

	cDataPtr := (*C.uint8_t)(&data[0])
	cDataSize := (C.size_t)(len(data))

	// Retrieve WebP features.
	if status := C.WebPGetFeatures(cDataPtr, cDataSize, &config.input); status != C.VP8_STATUS_OK {
		return nil, fmt.Errorf("could not get features from the data stream, return %s", statusString(status))
	}

	// Allocate output image
	outWidth, outHeight := calcOutputSize(config)
	img = image.NewRGBA(image.Rect(0, 0, outWidth, outHeight))

	// Set up output configurations
	config.output.colorspace = C.MODE_RGBA
	config.output.is_external_memory = 1

	// Allocate WebPRGBABuffer and fill in the pointers to output image
	buf := (*C.WebPRGBABuffer)(unsafe.Pointer(&config.output.u[0]))
	buf.rgba = (*C.uint8_t)(&img.Pix[0])
	buf.stride = C.int(img.Stride)
	buf.size = (C.size_t)(len(img.Pix))

	// Decode
	if status := C.WebPDecode(cDataPtr, cDataSize, config); status != C.VP8_STATUS_OK {
		return nil, fmt.Errorf("could not decode data stream, return %s", statusString(status))
	}

	return
}

// Decode reads a WebP image from r and returns it as an [image.Image].
// The type of Image returned depends on the WebP contents.
func Decode(r io.Reader) (image.Image, error) {
	// TODO: Find a way not to use ReadAll here
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return DecodeRGBA(data, nil)
}

// DecodeConfig returns the color model and dimensions of a WebP image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	// TODO: Find a way not to use ReadAll here
	data, err := io.ReadAll(r)
	if err != nil {
		return image.Config{}, err
	}

	features, err := GetFeatures(data)
	if err != nil {
		return image.Config{}, err
	}

	return image.Config{ColorModel: color.RGBAModel, Width: features.Width, Height: features.Height}, nil
}

// statusString converts the VP8StatsCode into a string.
func statusString(status C.VP8StatusCode) string {
	switch status {
	case C.VP8_STATUS_OK:
		return "VP8_STATUS_OK"
	case C.VP8_STATUS_OUT_OF_MEMORY:
		return "VP8_STATUS_OUT_OF_MEMORY"
	case C.VP8_STATUS_INVALID_PARAM:
		return "VP8_STATUS_INVALID_PARAM"
	case C.VP8_STATUS_BITSTREAM_ERROR:
		return "VP8_STATUS_BITSTREAM_ERROR"
	case C.VP8_STATUS_UNSUPPORTED_FEATURE:
		return "VP8_STATUS_UNSUPPORTED_FEATURE"
	case C.VP8_STATUS_SUSPENDED:
		return "VP8_STATUS_SUSPENDED"
	case C.VP8_STATUS_USER_ABORT:
		return "VP8_STATUS_USER_ABORT"
	case C.VP8_STATUS_NOT_ENOUGH_DATA:
		return "VP8_STATUS_NOT_ENOUGH_DATA"
	}
	return "Unexpected Status Code"
}

// initDecoderConfig initializes a decoder configuration and sets up the options.
func initDecoderConfig(options *DecoderOptions) (config *C.WebPDecoderConfig, err error) {
	// Initialize decoder config
	config = &C.WebPDecoderConfig{}
	if C.WebPInitDecoderConfig(config) == 0 {
		return nil, errors.New("could not initialize decoder config")
	}

	if options == nil {
		options = &DecoderOptions{}
	}

	// Set up decoder options
	if options.BypassFiltering {
		config.options.bypass_filtering = 1
	}
	if options.NoFancyUpsampling {
		config.options.no_fancy_upsampling = 1
	}
	if options.Crop.Max.X > 0 && options.Crop.Max.Y > 0 {
		config.options.use_cropping = 1
		config.options.crop_left = C.int(options.Crop.Min.X)
		config.options.crop_top = C.int(options.Crop.Min.Y)
		config.options.crop_width = C.int(options.Crop.Dx())
		config.options.crop_height = C.int(options.Crop.Dy())
	}
	if options.Scale.Max.X > 0 && options.Scale.Max.Y > 0 {
		config.options.use_scaling = 1
		config.options.scaled_width = C.int(options.Scale.Max.X)
		config.options.scaled_height = C.int(options.Scale.Max.Y)
	}
	if options.UseThreads {
		config.options.use_threads = 1
	}
	config.options.dithering_strength = C.int(options.DitheringStrength)

	return
}

// calcOutputSize retrieves width and height of output image from the decoder configuration.
func calcOutputSize(config *C.WebPDecoderConfig) (width, height int) {
	options := config.options
	if options.use_scaling > 0 {
		width = int(config.options.scaled_width)
		height = int(config.options.scaled_height)
		return
	}
	if config.options.use_cropping > 0 {
		width = int(config.options.crop_width)
		height = int(config.options.crop_height)
		return
	}

	width = int(config.input.width)
	height = int(config.input.height)
	return
}

func init() {
	image.RegisterFormat("webp", webPHeader, Decode, DecodeConfig)
}
