package models

type TTag struct {
	ID    uint64 `gorm:"primaryKey;autoIncrement;column:id;index"` // ID
	Slug  string `gorm:"not null;unique;column:slug"`              // 标签slug，小写字母/数字/中划线
	Label string `gorm:"not null;unique;column:label"`             // 标签名
	TimeHook
}

func (t *TTag) B() *BTag {
	return &BTag{
		ID:    t.ID,
		Slug:  t.Slug,
		Label: t.Label,
	}
}

type BTag struct {
	ID    uint64 `json:"id"`    // ID
	Slug  string `json:"slug"`  // 标签slug
	Label string `json:"label"` // 标签名
}
