package search

type ArticleSearch struct {
	ID      uint64 `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UserSearch struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TagSearch struct {
	ID   uint64 `json:"id"`
	Tag  string `json:"tag"`
	Slug string `json:"slug"`
}
