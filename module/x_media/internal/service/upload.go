package service

import (
	"bytes"
	"encoding/json"
	"github.com/EdlinOrg/prominentcolor"
	"image/color"
	"image/draw"
	"time"

	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/etag"
	"github.com/daidr/doulog-core/lib/localstorage"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/daidr/doulog-core/lib/webp"
	"github.com/disintegration/imaging"
	_ "github.com/jdeng/goheif"
	"github.com/jdeng/goheif/heif"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"mime/multipart"
)

func ExtractHeifRotation(ra io.ReaderAt) (int, error) {
	f := heif.Open(ra)
	primary, err := f.PrimaryItem()
	if err != nil {
		return 0, err
	}
	if primary == nil {
		return 0, errors.New("no primary item")
	}
	rotations := primary.Rotations()
	return rotations, nil
}

func Upload(sp *models.Scope, uid uint64, file *multipart.FileHeader) (uint64, error) {
	// start := time.Now()
	// 大小验证
	// fmt.Println("size:", file.Size)
	if int64(conf.C.Limit.Media.FileSize*1024) < file.Size {
		return 0, errors.New("wrong size")
	}

	// fmt.Println("t size:", time.Since(start).Milliseconds())

	f, err := file.Open()
	if err != nil {
		return 0, err
	}

	buf := &bytes.Buffer{}
	if _, err = io.Copy(buf, f); err != nil {
		return 0, err
	}
	// fmt.Println("t buf:", time.Since(start).Milliseconds())

	// 校验mime
	mime, err := utils.GetMPFDContentType(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return 0, err
	}
	if !utils.InStringSlice(mime, conf.C.Limit.Media.MIME) {
		return 0, errors.New("wrong mime")
	}
	// fmt.Println("mime:", mime)
	// fmt.Println("t mime:", time.Since(start).Milliseconds())

	// 获取 width 和 height
	_, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return 0, err
	}

	var img image.Image
	var gifImg *gif.GIF
	img, err = imaging.Decode(bytes.NewReader(buf.Bytes()), imaging.AutoOrientation(true))
	if err != nil {
		return 0, err
	}

	if format == "gif" {
		gifImg, err = gif.DecodeAll(bytes.NewReader(buf.Bytes()))
		if err != nil {
			return 0, err
		}
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// fmt.Println(width, height)
	// 限制尺寸
	if width > int(conf.C.Limit.Media.ImageSize) || height > int(conf.C.Limit.Media.ImageSize) {
		return 0, errors.New("excessive width or height")
	}

	// 获取etag
	tag, err := etag.New(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return 0, err
	}

	// fmt.Println("etag:", tag)
	// fmt.Println("t etag:", time.Since(start).Milliseconds())

	// 如果数据库里已经有该文件，直接返回
	media := &models.TMedia{}
	// 不发生not found错误即找到
	if err = sp.DB.PgSQL.Select("id").Where("etag = ?", tag).First(&media).Error; err == nil {
		return media.Id, nil
	}

	user, err := daos.NewUser(sp.DB).Get(uid)
	if err != nil {
		return 0, err
	}

	var title = file.Filename
	if title == "" {
		title = tag + "." + format
	}

	webpOptions, err := webp.ConfigPreset(webp.PresetDefault, 75)
	if err != nil {
		panic(err)
	}

	thumbnailWebpOptions, err := webp.ConfigPreset(webp.PresetDefault, 50)
	if err != nil {
		panic(err)
	}

	webpBuf := &bytes.Buffer{}

	if format == "heic" {
		// heic格式需要旋转
		ra := bytes.NewReader(buf.Bytes())
		rotations, err := ExtractHeifRotation(ra)
		if err != nil {
			return 0, err
		}

		if rotations != 0 {
			if rotations == 1 {
				img = imaging.Rotate90(img)
				width, height = height, width
			} else if rotations == 2 {
				img = imaging.Rotate180(img)
			} else if rotations == 3 {
				img = imaging.Rotate270(img)
				width, height = height, width
			}

			err = webp.EncodeRGBA(webpBuf, img, webpOptions)
			if err != nil {
				return 0, err
			}
		}
	} else if format == "gif" {
		// gif格式转换为webp
		loopCount := 0
		if gifImg.LoopCount == -1 {
			loopCount = 1
		} else if gifImg.LoopCount > 0 {
			loopCount = gifImg.LoopCount + 1
		}
		webpAnimImg, err := webp.NewAnimationEncoder(gifImg.Config.Width, gifImg.Config.Height, 7, 17, loopCount)
		if err != nil {
			return 0, err
		}
		defer webpAnimImg.Close()

		paintImage := image.NewRGBA(image.Rect(0, 0, width, height))
		bufferImage := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(paintImage, paintImage.Bounds(), gifImg.Image[0], image.Point{}, draw.Src)
		// 备份当前paintImage
		draw.Draw(bufferImage, bufferImage.Bounds(), paintImage, image.Point{}, draw.Src)
		err = webpAnimImg.AddFrame(paintImage, time.Duration(gifImg.Delay[0]*10)*time.Millisecond)
		if err != nil {
			return 0, err
		}

		for i := range gifImg.Image {
			if i == 0 {
				continue
			}
			// 判断使用 Disposal
			if gifImg.Disposal[i-1] == gif.DisposalNone {
				// 备份当前paintImage
				draw.Draw(bufferImage, bufferImage.Bounds(), paintImage, image.Point{}, draw.Src)
				// 直接使用当前帧
				draw.Draw(paintImage, paintImage.Bounds(), gifImg.Image[i], image.Point{}, draw.Over)
			} else if gifImg.Disposal[i-1] == gif.DisposalBackground {
				// 使用背景颜色
				pa, ok := gifImg.Image[i-1].ColorModel().(color.Palette)
				if !ok {
					return 0, errors.New("wrong color model")
				}
				draw.Draw(paintImage, paintImage.Bounds(), image.NewUniform(pa[gifImg.BackgroundIndex]), image.Point{}, draw.Src)
				// 绘制当前帧
				draw.Draw(paintImage, paintImage.Bounds(), gifImg.Image[i], image.Point{}, draw.Over)
			} else if gifImg.Disposal[i-1] == gif.DisposalPrevious {
				// 使用上一帧Buffer
				draw.Draw(paintImage, paintImage.Bounds(), bufferImage, image.Point{}, draw.Src)
				// 绘制当前帧
				draw.Draw(paintImage, paintImage.Bounds(), gifImg.Image[i], image.Point{}, draw.Over)
			} else {
				// 备份当前paintImage
				draw.Draw(bufferImage, bufferImage.Bounds(), paintImage, image.Point{}, draw.Src)
				// 完整覆盖
				draw.Draw(paintImage, paintImage.Bounds(), gifImg.Image[i], image.Point{}, draw.Over)
			}
			err = webpAnimImg.AddFrame(paintImage, time.Duration(gifImg.Delay[i]*10)*time.Millisecond)

			if err != nil {
				return 0, err
			}
		}

		err = webpAnimImg.Assemble(webpBuf)
		if err != nil {
			return 0, err
		}
	} else {
		// 直接转换webp
		err = webp.EncodeRGBA(webpBuf, img, webpOptions)
		if err != nil {
			return 0, err
		}
	}

	// 转换缩略图
	thumbnailBuf := &bytes.Buffer{}
	if width > 256 || height > 256 {
		var thumbnailImg image.Image
		if width > height {
			thumbnailImg = imaging.Resize(img, 256, 0, imaging.Lanczos)
		} else {
			thumbnailImg = imaging.Resize(img, 0, 256, imaging.Lanczos)
		}
		err = webp.Encode(thumbnailBuf, thumbnailImg, thumbnailWebpOptions)
		if err != nil {
			return 0, err
		}
	} else {
		err = webp.Encode(thumbnailBuf, img, thumbnailWebpOptions)
		if err != nil {
			return 0, err
		}
	}

	if err = localstorage.OriginMediaStore.SaveMedia(tag, format, buf.Bytes()); err != nil {
		return 0, err
	}

	if err = localstorage.MediaStore.SaveMedia(tag, "webp", webpBuf.Bytes()); err != nil {
		return 0, err
	}

	if err = localstorage.MediaStore.SaveThumbnail(tag, "webp", thumbnailBuf.Bytes()); err != nil {
		return 0, err
	}

	cols, err := prominentcolor.KmeansWithAll(3, img, prominentcolor.ArgumentDefault, prominentcolor.DefaultSize, []prominentcolor.ColorBackgroundMask{})
	if err != nil {
		return 0, err
	}

	// 最多五个颜色
	ProminentColors := pq.StringArray{}
	for i := 0; i < int(math.Min(5, float64(len(cols)))); i++ {
		ProminentColors = append(ProminentColors, cols[i].AsString())
	}

	fileSizes, err := json.Marshal(&models.MediaFileSizes{
		Original:   uint64(buf.Len()),
		Compressed: uint64(webpBuf.Len()),
		Thumbnail:  uint64(thumbnailBuf.Len()),
	})

	if err != nil {
		return 0, err
	}

	media = &models.TMedia{
		Etag:           tag,
		Ext:            format,
		Src:            conf.ImageSrcLocal,
		Width:          width,
		Height:         height,
		MIME:           mime,
		Owner:          *user,
		Title:          title,
		Alt:            "",
		ProminentColor: ProminentColors,
		FileSizes:      fileSizes,
	}

	// TODO 同tag的并发处理
	if err = sp.DB.PgSQL.Create(media).Error; err != nil {
		return 0, err
	}
	// fmt.Println("t create:", time.Since(start).Milliseconds())

	return media.Id, err
}
