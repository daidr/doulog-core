package search

import (
	"strconv"
)

func IndexArticle(article *ArticleSearch) error {
	return articleIndex.Index(strconv.FormatUint(article.ID, 10), *article)
}

func IndexTag(tag TagSearch) error {
	return tagIndex.Index(strconv.FormatUint(tag.ID, 10), tag)
}

func IndexUser(user UserSearch) error {
	return userIndex.Index(strconv.FormatUint(user.ID, 10), user)
}
