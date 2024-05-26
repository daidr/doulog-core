package localstorage

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/daidr/doulog-core/lib/utils"
	"os"
	"strings"
)

type LocalStorage struct {
	Biz  string
	Root string
}

var MediaStore *LocalStorage
var OriginMediaStore *LocalStorage

func Init() error {
	var err error
	if MediaStore, err = initStore("media"); err != nil {
		return err
	}

	return nil
}

func initStore(biz string) (*LocalStorage, error) {
	var err error

	path := "./data/store/" + biz
	if !utils.PathExist(path) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return &LocalStorage{
		Biz:  biz,
		Root: path,
	}, nil
}

func EtagToHex(etag string) (string, error) {
	rawEtag, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(strings.ReplaceAll(etag, "_", "/"), "-", "+"))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawEtag), nil
}

func (store *LocalStorage) IsMediaExist(etag string, ext string) (bool, error) {
	var (
		err     error
		hexEtag string
	)

	if hexEtag, err = EtagToHex(etag); err != nil {
		return false, err
	}

	prefix := hexEtag[:2]
	prefix2 := hexEtag[2:4]

	path := store.Root + "/" + prefix + "/" + prefix2 + "/" + hexEtag + "." + ext
	return utils.PathExist(path), nil
}

func (store *LocalStorage) SaveMedia(etag string, ext string, data []byte) error {
	var (
		err     error
		hexEtag string
	)

	if hexEtag, err = EtagToHex(etag); err != nil {
		return err
	}

	prefix := hexEtag[:2]
	prefix2 := hexEtag[2:4]

	path := store.Root + "/" + prefix + "/" + prefix2
	println("path:", path)
	if !utils.PathExist(path) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	path = path + "/" + hexEtag + "." + ext

	//if utils.PathExist(path) {
	//	return nil
	//}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.Write(data)
	return err
}

func (store *LocalStorage) SaveThumbnail(etag string, ext string, data []byte) error {
	var (
		err     error
		hexEtag string
	)

	if hexEtag, err = EtagToHex(etag); err != nil {
		return err
	}

	prefix := hexEtag[:2]
	prefix2 := hexEtag[2:4]

	path := store.Root + "/" + prefix + "/" + prefix2
	if !utils.PathExist(path) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	path = path + "/" + hexEtag + "_thumbnail." + ext

	//if utils.PathExist(path) {
	//	return nil
	//}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.Write(data)
	return err
}

func (store *LocalStorage) GetMedia(etag string, ext string) ([]byte, error) {
	var (
		err     error
		hexEtag string
	)

	if hexEtag, err = EtagToHex(etag); err != nil {
		return nil, err
	}

	path := store.Root + "/" + hexEtag[:2] + "/" + hexEtag[2:4] + "/" + hexEtag + "." + ext

	if !utils.PathExist(path) {
		return nil, nil
	}

	return os.ReadFile(path)
}
