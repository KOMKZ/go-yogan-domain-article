package article

import (
	"context"

	"github.com/KOMKZ/go-yogan-domain-article/model"
	"gorm.io/gorm"
)

// ArticleGORMRepository GORM 文章仓储实现
type ArticleGORMRepository struct {
	db *gorm.DB
}

func NewArticleGORMRepository(db *gorm.DB) *ArticleGORMRepository {
	return &ArticleGORMRepository{db: db}
}

func (r *ArticleGORMRepository) Create(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *ArticleGORMRepository) Update(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *ArticleGORMRepository) FindByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleGORMRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).Update("status", model.StatusDeleted).Error
}

func (r *ArticleGORMRepository) Paginate(ctx context.Context, page, pageSize int, ownerId *uint, ownerType, articleType, title string, folderID *uint) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Article{}).Where("status != ?", model.StatusDeleted)

	if ownerId != nil {
		query = query.Where("owner_id = ?", *ownerId)
	}
	if ownerType != "" {
		query = query.Where("owner_type = ?", ownerType)
	}
	if articleType != "" {
		query = query.Where("article_type = ?", articleType)
	}
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if folderID != nil {
		query = query.Where("folder_id = ?", *folderID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *ArticleGORMRepository) CountByFolderID(ctx context.Context, folderID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Article{}).
		Where("folder_id = ? AND status != ?", folderID, model.StatusDeleted).
		Count(&count).Error
	return count, err
}

func (r *ArticleGORMRepository) FindByFolderID(ctx context.Context, folderID uint) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.WithContext(ctx).
		Where("folder_id = ? AND status != ?", folderID, model.StatusDeleted).
		Order("created_at DESC").
		Find(&articles).Error
	return articles, err
}

// MarkdownArticleGORMRepository GORM Markdown文章仓储实现
type MarkdownArticleGORMRepository struct {
	db *gorm.DB
}

func NewMarkdownArticleGORMRepository(db *gorm.DB) *MarkdownArticleGORMRepository {
	return &MarkdownArticleGORMRepository{db: db}
}

func (r *MarkdownArticleGORMRepository) Create(ctx context.Context, article *model.MarkdownArticle) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *MarkdownArticleGORMRepository) Update(ctx context.Context, article *model.MarkdownArticle) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *MarkdownArticleGORMRepository) FindByArticleID(ctx context.Context, articleID uint) (*model.MarkdownArticle, error) {
	var article model.MarkdownArticle
	err := r.db.WithContext(ctx).Where("article_id = ?", articleID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *MarkdownArticleGORMRepository) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.MarkdownArticle{}).Error
}

// RichTextArticleGORMRepository GORM 富文本文章仓储实现
type RichTextArticleGORMRepository struct {
	db *gorm.DB
}

func NewRichTextArticleGORMRepository(db *gorm.DB) *RichTextArticleGORMRepository {
	return &RichTextArticleGORMRepository{db: db}
}

func (r *RichTextArticleGORMRepository) Create(ctx context.Context, article *model.RichTextArticle) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *RichTextArticleGORMRepository) Update(ctx context.Context, article *model.RichTextArticle) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *RichTextArticleGORMRepository) FindByArticleID(ctx context.Context, articleID uint) (*model.RichTextArticle, error) {
	var article model.RichTextArticle
	err := r.db.WithContext(ctx).Where("article_id = ?", articleID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *RichTextArticleGORMRepository) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.RichTextArticle{}).Error
}

// TableArticleGORMRepository GORM 表格文章仓储实现
type TableArticleGORMRepository struct {
	db *gorm.DB
}

func NewTableArticleGORMRepository(db *gorm.DB) *TableArticleGORMRepository {
	return &TableArticleGORMRepository{db: db}
}

func (r *TableArticleGORMRepository) Create(ctx context.Context, article *model.TableArticle) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *TableArticleGORMRepository) Update(ctx context.Context, article *model.TableArticle) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *TableArticleGORMRepository) FindByArticleID(ctx context.Context, articleID uint) (*model.TableArticle, error) {
	var article model.TableArticle
	err := r.db.WithContext(ctx).Where("article_id = ?", articleID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *TableArticleGORMRepository) FindByTableID(ctx context.Context, tableID string) (*model.TableArticle, error) {
	var article model.TableArticle
	err := r.db.WithContext(ctx).Where("table_id = ?", tableID).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *TableArticleGORMRepository) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.TableArticle{}).Error
}

// TableArticleRowGORMRepository GORM 表格行仓储实现
type TableArticleRowGORMRepository struct {
	db *gorm.DB
}

func NewTableArticleRowGORMRepository(db *gorm.DB) *TableArticleRowGORMRepository {
	return &TableArticleRowGORMRepository{db: db}
}

func (r *TableArticleRowGORMRepository) Create(ctx context.Context, row *model.TableArticleRow) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *TableArticleRowGORMRepository) BatchCreate(ctx context.Context, rows []model.TableArticleRow) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(rows, 100).Error
}

func (r *TableArticleRowGORMRepository) FindByArticleID(ctx context.Context, articleID uint) ([]model.TableArticleRow, error) {
	var rows []model.TableArticleRow
	err := r.db.WithContext(ctx).Where("article_id = ?", articleID).Order("row_index ASC").Find(&rows).Error
	return rows, err
}

func (r *TableArticleRowGORMRepository) DeleteByArticleID(ctx context.Context, articleID uint) error {
	return r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.TableArticleRow{}).Error
}

func (r *TableArticleRowGORMRepository) ReplaceAll(ctx context.Context, articleID uint, rows []model.TableArticleRow) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除旧数据
		if err := tx.Where("article_id = ?", articleID).Delete(&model.TableArticleRow{}).Error; err != nil {
			return err
		}
		// 插入新数据
		if len(rows) > 0 {
			return tx.CreateInBatches(rows, 100).Error
		}
		return nil
	})
}
