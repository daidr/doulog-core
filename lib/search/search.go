package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/daidr/doulog-core/lib/utils"
	"os"
	"path"
)

var (
	articleIndex bleve.Index
	tagIndex     bleve.Index
	userIndex    bleve.Index
)

func Init() error {
	var err error

	if tagIndex, err = initIndex("tag"); err != nil {
		return err
	}
	if userIndex, err = initIndex("user"); err != nil {
		return err
	}
	if articleIndex, err = initIndex("article"); err != nil {
		return err
	}
	return nil
}

func initIndex(biz string) (bleve.Index, error) {
	var (
		index bleve.Index
		err   error
	)

	if !utils.PathExist("./data") {
		if err = os.MkdirAll("./data", os.ModePerm); err != nil {
			return nil, err
		}
	}

	filePath := path.Join("./data", biz+".index")
	if !utils.PathExist(filePath) {
		index, err = bleve.New(filePath, bleve.NewIndexMapping())
		if err != nil {
			return nil, err
		}
		if err = index.Close(); err != nil {
			return nil, err
		}
	}

	if index, err = bleve.Open(filePath); err != nil {
		return nil, err
	}
	return index, err
}
