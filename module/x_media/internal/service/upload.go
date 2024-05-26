package service

import (
	"bytes"
	"encoding/json"
	"github.com/EdlinOrg/prominentcolor"

	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"mime/multipart"

	"github.com/daidr/doulog-core/lib/conf"
	"github.com/daidr/doulog-core/lib/daos"
	"github.com/daidr/doulog-core/lib/etag"
	"github.com/daidr/doulog-core/lib/localstorage"
	"github.com/daidr/doulog-core/lib/models"
	"github.com/daidr/doulog-core/lib/utils"
	"github.com/disintegration/imaging"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var AllowedMime = []string{
	"image/jpeg", "image/png", "image/gif", "image/webp",
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
	if !utils.InStringSlice(mime, AllowedMime) {
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
	img, err = imaging.Decode(bytes.NewReader(buf.Bytes()), imaging.AutoOrientation(true))
	if err != nil {
		return 0, err
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

	// 转换缩略图
	thumbnailBuf := &bytes.Buffer{}
	if width > 256 || height > 256 {
		var thumbnailImg image.Image
		if width > height {
			thumbnailImg = imaging.Resize(img, 256, 0, imaging.Lanczos)
		} else {
			thumbnailImg = imaging.Resize(img, 0, 256, imaging.Lanczos)
		}
		// 质量50%
		if err = imaging.Encode(thumbnailBuf, thumbnailImg, imaging.JPEG, imaging.JPEGQuality(50)); err != nil {
			return 0, err
		}
	} else {
		// 质量50%
		if err = imaging.Encode(thumbnailBuf, img, imaging.JPEG, imaging.JPEGQuality(50)); err != nil {
			return 0, err
		}
	}

	if err = localstorage.MediaStore.SaveMedia(tag, format, buf.Bytes()); err != nil {
		return 0, err
	}

	if err = localstorage.MediaStore.SaveThumbnail(tag, "jpg", thumbnailBuf.Bytes()); err != nil {
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
		Original:  uint64(buf.Len()),
		Thumbnail: uint64(thumbnailBuf.Len()),
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
