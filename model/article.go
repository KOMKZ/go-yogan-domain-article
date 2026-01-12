package model

import "time"

// Article 文章总表实体
type Article struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	ArticleType string    `gorm:"size:50;not null;index" json:"articleType"` // table, markdown, rich_text
	OwnerID     uint      `gorm:"not null" json:"ownerId"`
	OwnerType   string    `gorm:"size:50;not null;index" json:"ownerType"` // user, admin, team
	Status      int       `gorm:"not null;default:1;index" json:"status"`  // 0=草稿, 1=已发布, 2=已删除
	CreatedAt   time.Time `gorm:"not null;index" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"not null" json:"updatedAt"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}

// ArticleType 文章类型常量
const (
	ArticleTypeTable    = "table"
	ArticleTypeMarkdown = "markdown"
	ArticleTypeRichText = "rich_text"
)

// OwnerType 所有者类型常量
const (
	OwnerTypeUser  = "user"
	OwnerTypeAdmin = "admin"
	OwnerTypeTeam  = "team"
)

// Status 状态常量
const (
	StatusDraft     = 0
	StatusPublished = 1
	StatusDeleted   = 2
)

// IsDeleted 是否已删除
func (a *Article) IsDeleted() bool {
	return a.Status == StatusDeleted
}

// IsPublished 是否已发布
func (a *Article) IsPublished() bool {
	return a.Status == StatusPublished
}
