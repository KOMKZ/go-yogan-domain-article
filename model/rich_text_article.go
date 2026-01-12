package model

import "time"

// RichTextArticle 富文本文章实体
type RichTextArticle struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	ArticleID     uint      `gorm:"uniqueIndex;not null" json:"articleId"`
	Content       string    `gorm:"type:longtext;not null" json:"content"` // HTML内容
	FormatVersion string    `gorm:"size:20;default:'1.0'" json:"formatVersion"`
	CreatedAt     time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"not null" json:"updatedAt"`
}

// TableName 指定表名
func (RichTextArticle) TableName() string {
	return "rich_text_articles"
}
