package search

import "github.com/blevesearch/bleve/v2"

// TagWithPage page 从0开始
func TagWithPage(s string, pn int, ps int) []TagSearch {
	// 先尝试QueryString模式，如果失败再尝试性能较低但效果比较好的WildCard模式
	search := bleve.NewSearchRequestOptions(bleve.NewQueryStringQuery("tag:"+s), ps, pn*ps, false)
	results := tagEnd(search)
	if len(results) > 0 {
		return results
	}

	search = bleve.NewSearchRequestOptions(bleve.NewWildcardQuery("*"+s+"*"), ps, pn*ps, false)
	return tagEnd(search)
}

func tagEnd(search *bleve.SearchRequest) []TagSearch {
	search.Fields = []string{"*"}
	result, err := tagIndex.Search(search)
	if err != nil {
		return nil
	}
	tags := make([]TagSearch, 0)
	for _, r := range result.Hits {
		tags = append(tags, TagSearch{
			ID:   uint64(r.Fields["id"].(float64)),
			Tag:  r.Fields["tag"].(string),
			Slug: r.Fields["slug"].(string),
		})
	}
	return tags
}
