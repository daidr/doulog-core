package search

import "github.com/blevesearch/bleve/v2"

// ArticleWithPage pn 从0开始,ps: page size
func ArticleWithPage(s string, pn int, ps int) []*ArticleSearch {
	query := bleve.NewQueryStringQuery("title:" + s + "^3 content:" + s)
	search := bleve.NewSearchRequestOptions(query, ps, pn*ps, false)
	return articleEnd(search)
}

func articleEnd(search *bleve.SearchRequest) []*ArticleSearch {
	search.Fields = []string{"*"}
	result, err := articleIndex.Search(search)
	if err != nil {
		return nil
	}
	articles := make([]*ArticleSearch, 0)
	for _, r := range result.Hits {
		articles = append(articles, &ArticleSearch{
			ID:      uint64(r.Fields["id"].(float64)),
			Title:   r.Fields["title"].(string),
			Content: r.Fields["content"].(string),
		})
	}
	return articles
}
