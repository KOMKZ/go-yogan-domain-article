# go-yogan-domain-article

Yogan Framework 的文章领域包，支持多种文章类型：Markdown、富文本、表格。

## 安装

```bash
go get github.com/KOMKZ/go-yogan-domain-article@latest
```

## 功能

- 多类型文章支持（Markdown、富文本、表格）
- 文章 CRUD 操作
- 表格文章结构管理
- 表格行数据批量操作
- 软删除支持

## 文章类型

| 类型 | 说明 |
|------|------|
| `markdown` | Markdown 文章 |
| `rich_text` | 富文本文章 |
| `table` | 表格文章 |

## 使用示例

```go
import (
    article "github.com/KOMKZ/go-yogan-domain-article"
    "github.com/KOMKZ/go-yogan-framework/logger"
    "gorm.io/gorm"
)

// 初始化
func InitArticleService(db *gorm.DB) *article.Service {
    log := logger.GetLogger("yogan")
    return article.NewService(
        article.NewArticleGORMRepository(db),
        article.NewMarkdownArticleGORMRepository(db),
        article.NewRichTextArticleGORMRepository(db),
        article.NewTableArticleGORMRepository(db),
        article.NewTableArticleRowGORMRepository(db),
        log,
    )
}

// 创建富文本文章
func CreateRichText(svc *article.Service) (*article.model.Article, error) {
    return svc.CreateRichTextArticle(ctx, &article.CreateRichTextArticleInput{
        Title:     "我的文章",
        OwnerID:   1,
        OwnerType: "admin",
        Content:   "<p>Hello World</p>",
    })
}

// 创建表格文章
func CreateTable(svc *article.Service) (*article.model.Article, error) {
    return svc.CreateTableArticle(ctx, &article.CreateTableArticleInput{
        Title:     "数据表格",
        TableID:   "table_001",
        OwnerID:   1,
        OwnerType: "admin",
        Structure: []map[string]interface{}{
            {"field": "name", "title": "名称"},
            {"field": "value", "title": "值"},
        },
        Data: []map[string]interface{}{
            {"name": "item1", "value": "100"},
        },
    })
}
```

## 数据模型

### Article (主表)
```go
type Article struct {
    ID          uint
    Title       string
    ArticleType string  // table, markdown, rich_text
    OwnerID     uint
    OwnerType   string  // user, admin, team
    Status      int     // 0=草稿, 1=已发布, 2=已删除
}
```

### TableArticle (表格结构)
```go
type TableArticle struct {
    ID          uint
    ArticleID   uint
    TableID     string
    Structure   JSONArray  // 表结构定义
    ColumnOrder JSONArray
    Filters     JSONArray
    Version     int
}
```

## 依赖

- [go-yogan-framework](https://github.com/KOMKZ/go-yogan-framework) - 核心框架
- [gorm](https://gorm.io) - ORM

## License

MIT
