package model

type Article struct {
	ID        uint32 `gorm:"primary_key" json:"id"`
	ArticleID string
	Title     string
	Content   string
}
