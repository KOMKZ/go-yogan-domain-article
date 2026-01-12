package model

import "time"

// MarkdownArticle Markdown文章实体
type MarkdownArticle struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	ArticleID     uint      `gorm:"uniqueIndex;not null" json:"articleId"`
	Content       string    `gorm:"type:text;not null" json:"content"`         // Markdown内容
	HTMLContent   string    `gorm:"type:longtext" json:"htmlContent"`          // 渲染后的HTML（缓存）
	FormatVersion string    `gorm:"size:20;default:'1.0'" json:"formatVersion"`
	CreatedAt     time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"not null" json:"updatedAt"`
}

// TableName 指定表名
func (MarkdownArticle) TableName() string {
	return "markdown_articles"
}
