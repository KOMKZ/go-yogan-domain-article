package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// TableArticle 表格文章实体
type TableArticle struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ArticleID   uint      `gorm:"uniqueIndex;not null" json:"articleId"`
	TableID     string    `gorm:"size:100;uniqueIndex;not null" json:"tableId"` // 前端tableId
	Structure   JSONArray `gorm:"type:json;not null" json:"structure"`          // 表结构定义
	ColumnOrder JSONArray `gorm:"type:json" json:"columnOrder"`                 // 列顺序
	Filters     JSONArray `gorm:"type:json" json:"filters"`                     // 过滤条件
	Version     int       `gorm:"not null;default:1" json:"version"`
	CreatedAt   time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"not null" json:"updatedAt"`
}

// TableName 指定表名
func (TableArticle) TableName() string {
	return "table_articles"
}

// TableArticleRow 表格文章行数据
type TableArticleRow struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `gorm:"index;not null" json:"articleId"`
	RowData   JSONMap   `gorm:"type:json;not null" json:"rowData"`
	RowIndex  *int      `gorm:"index" json:"rowIndex"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt"`
}

// TableName 指定表名
func (TableArticleRow) TableName() string {
	return "table_article_rows"
}

// TableArticleStructureHistory 表格结构变更历史
type TableArticleStructureHistory struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	ArticleID         uint      `gorm:"index;not null" json:"articleId"`
	Version           int       `gorm:"not null" json:"version"`
	Structure         JSONArray `gorm:"type:json;not null" json:"structure"`
	ChangeType        string    `gorm:"size:50" json:"changeType"`
	ChangeDescription string    `gorm:"type:text" json:"changeDescription"`
	CreatedAt         time.Time `gorm:"not null" json:"createdAt"`
}

// TableName 指定表名
func (TableArticleStructureHistory) TableName() string {
	return "table_article_structure_history"
}

// JSONArray JSON数组类型
type JSONArray []map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// JSONMap JSON对象类型
type JSONMap map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
