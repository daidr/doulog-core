package filter

import (
	_ "embed"
	"github.com/importcjj/sensitive"
	"strings"
)

type text struct {
	filter *sensitive.Filter
}

//go:embed keywords.txt
var keywords string

var Text = &text{}

func init() {
	Text.filter = sensitive.New()
	Text.filter.AddWord(strings.Split(keywords, "|")...)
}

// Replace 把词语中的字符替换成*
func (t *text) Replace(text string) string {
	return t.filter.Replace(text, '*')
}

// FindIn 查找并返回第一个敏感词，如果没有则返回false
func (t *text) FindIn(text string) (bool, string) {
	return t.filter.FindIn(text)
}

// FindAll 查找内容中的全部敏感词，以数组返回。
func (t *text) FindAll(text string) []string {
	return t.filter.FindAll(text)
}

// Validate 验证内容是否ok，如果含有敏感词，则返回false和第一个敏感词。
func (t *text) Validate(text string) (bool, string) {
	return t.filter.Validate(text)
}

// Remove 直接移除词语
func (t *text) Remove(text string) string {
	return t.filter.Filter(text)
}
