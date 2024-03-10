package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/pkg/errors"
	"strconv"
)

// UserWithPage page 从0开始
func UserWithPage(name string, pn int, ps int) []UserSearch {
	// 先尝试QueryString模式，如果失败再尝试性能较低但效果比较好的WildCard模式
	search := bleve.NewSearchRequestOptions(bleve.NewQueryStringQuery("name:"+name), ps, pn*ps, false)
	results := userEnd(search)
	if len(results) > 0 {
		return results
	}

	search = bleve.NewSearchRequestOptions(bleve.NewWildcardQuery("*"+name+"*"), ps, pn*ps, false)
	return userEnd(search)
}

func UserGet(id uint64) (UserSearch, error) {
	query := bleve.NewDocIDQuery([]string{strconv.FormatUint(id, 10)})
	search := bleve.NewSearchRequest(query)
	results := userEnd(search)

	if len(results) != 1 {
		return UserSearch{}, errors.New("can not find user")
	}

	return results[0], nil
}

func UserEdit(search UserSearch) error {
	b := userIndex.NewBatch()

	b.Delete(strconv.FormatUint(search.ID, 10))
	if err := b.Index(strconv.FormatUint(search.ID, 10), search); err != nil {
		return err
	}

	return userIndex.Batch(b)
}

func userEnd(search *bleve.SearchRequest) []UserSearch {
	search.Fields = []string{"*"}
	result, err := userIndex.Search(search)
	if err != nil {
		return nil
	}

	users := make([]UserSearch, 0)
	for _, r := range result.Hits {
		users = append(users, UserSearch{
			ID:    uint64(r.Fields["id"].(float64)),
			Name:  r.Fields["name"].(string),
			Email: r.Fields["email"].(string),
		})
	}
	return users
}
