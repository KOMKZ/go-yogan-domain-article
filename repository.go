package article

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-article/model"
)

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	Update(ctx context.Context, article *model.Article) error
	FindByID(ctx context.Context, id uint) (*model.Article, error)
	Delete(ctx context.Context, id uint) error
	Paginate(ctx context.Context, page, pageSize int, ownerId *uint, ownerType, articleType, title string, folderID *uint) ([]model.Article, int64, error)
	// PaginateByFolderIDs 分页查询（支持多个文件夹ID，用于树形筛选）
	PaginateByFolderIDs(ctx context.Context, page, pageSize int, ownerId *uint, ownerType, articleType, title string, folderIDs []uint) ([]model.Article, int64, error)
	CountByFolderID(ctx context.Context, folderID uint) (int64, error)
	FindByFolderID(ctx context.Context, folderID uint) ([]model.Article, error)
}

// MarkdownArticleRepository Markdown文章仓储接口
type MarkdownArticleRepository interface {
	Create(ctx context.Context, article *model.MarkdownArticle) error
	Update(ctx context.Context, article *model.MarkdownArticle) error
	FindByArticleID(ctx context.Context, articleID uint) (*model.MarkdownArticle, error)
	DeleteByArticleID(ctx context.Context, articleID uint) error
}

// RichTextArticleRepository 富文本文章仓储接口
type RichTextArticleRepository interface {
	Create(ctx context.Context, article *model.RichTextArticle) error
	Update(ctx context.Context, article *model.RichTextArticle) error
	FindByArticleID(ctx context.Context, articleID uint) (*model.RichTextArticle, error)
	DeleteByArticleID(ctx context.Context, articleID uint) error
}

// TableArticleRepository 表格文章仓储接口
type TableArticleRepository interface {
	Create(ctx context.Context, article *model.TableArticle) error
	Update(ctx context.Context, article *model.TableArticle) error
	FindByArticleID(ctx context.Context, articleID uint) (*model.TableArticle, error)
	FindByTableID(ctx context.Context, tableID string) (*model.TableArticle, error)
	DeleteByArticleID(ctx context.Context, articleID uint) error
}

// TableArticleRowRepository 表格行数据仓储接口
type TableArticleRowRepository interface {
	Create(ctx context.Context, row *model.TableArticleRow) error
	BatchCreate(ctx context.Context, rows []model.TableArticleRow) error
	FindByArticleID(ctx context.Context, articleID uint) ([]model.TableArticleRow, error)
	DeleteByArticleID(ctx context.Context, articleID uint) error
	ReplaceAll(ctx context.Context, articleID uint, rows []model.TableArticleRow) error
}
