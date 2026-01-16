package article

import (
	"context"
	"errors"
	"time"

	"github.com/KOMKZ/go-yogan-domain-article/model"
	"github.com/KOMKZ/go-yogan-framework/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service 文章服务
type Service struct {
	articleRepo  ArticleRepository
	markdownRepo MarkdownArticleRepository
	richTextRepo RichTextArticleRepository
	tableRepo    TableArticleRepository
	tableRowRepo TableArticleRowRepository
	logger       *logger.CtxZapLogger
}

// NewService 创建文章服务
func NewService(
	articleRepo ArticleRepository,
	markdownRepo MarkdownArticleRepository,
	richTextRepo RichTextArticleRepository,
	tableRepo TableArticleRepository,
	tableRowRepo TableArticleRowRepository,
	log *logger.CtxZapLogger,
) *Service {
	return &Service{
		articleRepo:  articleRepo,
		markdownRepo: markdownRepo,
		richTextRepo: richTextRepo,
		tableRepo:    tableRepo,
		tableRowRepo: tableRowRepo,
		logger:       log,
	}
}

// ==================== 通用文章操作 ====================

// CreateArticleInput 创建文章输入
type CreateArticleInput struct {
	Title       string
	ArticleType string
	FolderID    *uint
	OwnerID     uint
	OwnerType   string
}

// CreateArticle 创建文章
func (s *Service) CreateArticle(ctx context.Context, input *CreateArticleInput) (*model.Article, error) {
	// 验证文章类型
	if !isValidArticleType(input.ArticleType) {
		return nil, ErrBadRequest.WithMsgf("不支持的文章类型: %s", input.ArticleType)
	}

	article := &model.Article{
		Title:       input.Title,
		ArticleType: input.ArticleType,
		FolderID:    input.FolderID,
		OwnerID:     input.OwnerID,
		OwnerType:   input.OwnerType,
		Status:      model.StatusPublished,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.articleRepo.Create(ctx, article); err != nil {
		s.logger.ErrorCtx(ctx, "创建文章失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "文章创建成功", zap.Uint("article_id", article.ID))
	return article, nil
}

// GetArticle 获取文章详情
func (s *Service) GetArticle(ctx context.Context, id uint) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound.WithMsg("文章不存在")
		}
		return nil, ErrDatabaseError.Wrap(err)
	}

	if article.IsDeleted() {
		return nil, ErrDeleted.WithMsg("文章已删除")
	}

	return article, nil
}

// UpdateArticleInput 更新文章输入
type UpdateArticleInput struct {
	Title    *string
	Status   *int
	FolderID **uint // 二级指针：nil=不更新, *nil=清空, *value=设置新值
}

// UpdateArticle 更新文章
func (s *Service) UpdateArticle(ctx context.Context, id uint, input *UpdateArticleInput) error {
	article, err := s.GetArticle(ctx, id)
	if err != nil {
		return err
	}

	if input.Title != nil {
		article.Title = *input.Title
	}
	if input.Status != nil {
		article.Status = *input.Status
	}
	if input.FolderID != nil {
		article.FolderID = *input.FolderID
	}
	article.UpdatedAt = time.Now()

	if err := s.articleRepo.Update(ctx, article); err != nil {
		s.logger.ErrorCtx(ctx, "更新文章失败", zap.Uint("article_id", id), zap.Error(err))
		return ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "文章更新成功", zap.Uint("article_id", id))
	return nil
}

// DeleteArticle 删除文章（软删除）
func (s *Service) DeleteArticle(ctx context.Context, id uint) error {
	_, err := s.GetArticle(ctx, id)
	if err != nil {
		return err
	}

	if err := s.articleRepo.Delete(ctx, id); err != nil {
		s.logger.ErrorCtx(ctx, "删除文章失败", zap.Uint("article_id", id), zap.Error(err))
		return ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "文章删除成功", zap.Uint("article_id", id))
	return nil
}

// PageResult 分页结果
type PageResult struct {
	Records     []model.Article `json:"records"`
	Total       int64           `json:"total"`
	Size        int             `json:"size"`
	Current     int             `json:"current"`
	Pages       int             `json:"pages"`
	HasPrevious bool            `json:"hasPrevious"`
	HasNext     bool            `json:"hasNext"`
	IsFirst     bool            `json:"isFirst"`
	IsLast      bool            `json:"isLast"`
}

// ListArticles 分页查询文章
func (s *Service) ListArticles(ctx context.Context, page, size int, ownerId *uint, ownerType, articleType, title string, folderID *uint) (*PageResult, error) {
	articles, total, err := s.articleRepo.Paginate(ctx, page, size, ownerId, ownerType, articleType, title, folderID)
	if err != nil {
		s.logger.ErrorCtx(ctx, "查询文章列表失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	pages := int(total) / size
	if int(total)%size > 0 {
		pages++
	}

	return &PageResult{
		Records:     articles,
		Total:       total,
		Size:        size,
		Current:     page,
		Pages:       pages,
		HasPrevious: page > 1,
		HasNext:     page < pages,
		IsFirst:     page == 1,
		IsLast:      page >= pages,
	}, nil
}

// ListArticlesByFolderIDs 分页查询文章（支持多个文件夹ID，用于树形筛选）
func (s *Service) ListArticlesByFolderIDs(ctx context.Context, page, size int, ownerId *uint, ownerType, articleType, title string, folderIDs []uint) (*PageResult, error) {
	articles, total, err := s.articleRepo.PaginateByFolderIDs(ctx, page, size, ownerId, ownerType, articleType, title, folderIDs)
	if err != nil {
		s.logger.ErrorCtx(ctx, "查询文章列表失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	pages := int(total) / size
	if int(total)%size > 0 {
		pages++
	}

	return &PageResult{
		Records:     articles,
		Total:       total,
		Size:        size,
		Current:     page,
		Pages:       pages,
		HasPrevious: page > 1,
		HasNext:     page < pages,
		IsFirst:     page == 1,
		IsLast:      page >= pages,
	}, nil
}

// ==================== 富文本文章操作 ====================

// CreateRichTextArticleInput 创建富文本文章输入
type CreateRichTextArticleInput struct {
	Title     string
	FolderID  *uint
	OwnerID   uint
	OwnerType string
	Content   string
}

// CreateRichTextArticle 创建富文本文章
func (s *Service) CreateRichTextArticle(ctx context.Context, input *CreateRichTextArticleInput) (*model.Article, error) {
	// 1. 创建主表
	article, err := s.CreateArticle(ctx, &CreateArticleInput{
		Title:       input.Title,
		ArticleType: model.ArticleTypeRichText,
		FolderID:    input.FolderID,
		OwnerID:     input.OwnerID,
		OwnerType:   input.OwnerType,
	})
	if err != nil {
		return nil, err
	}

	// 2. 创建富文本内容
	richText := &model.RichTextArticle{
		ArticleID:     article.ID,
		Content:       input.Content,
		FormatVersion: "1.0",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.richTextRepo.Create(ctx, richText); err != nil {
		s.logger.ErrorCtx(ctx, "创建富文本内容失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	return article, nil
}

// RichTextArticleContent 富文本文章完整内容
type RichTextArticleContent struct {
	Article *model.Article `json:"article"`
	Content string         `json:"content"`
}

// GetRichTextArticleContent 获取富文本文章内容
func (s *Service) GetRichTextArticleContent(ctx context.Context, articleID uint) (*RichTextArticleContent, error) {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}

	if article.ArticleType != model.ArticleTypeRichText {
		return nil, ErrBadRequest.WithMsg("该文章不是富文本类型")
	}

	richText, err := s.richTextRepo.FindByArticleID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &RichTextArticleContent{Article: article, Content: ""}, nil
		}
		return nil, ErrDatabaseError.Wrap(err)
	}

	return &RichTextArticleContent{
		Article: article,
		Content: richText.Content,
	}, nil
}

// UpdateRichTextContent 更新富文本内容
func (s *Service) UpdateRichTextContent(ctx context.Context, articleID uint, content string) error {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return err
	}

	if article.ArticleType != model.ArticleTypeRichText {
		return ErrBadRequest.WithMsg("该文章不是富文本类型")
	}

	richText, err := s.richTextRepo.FindByArticleID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			richText = &model.RichTextArticle{
				ArticleID:     articleID,
				Content:       content,
				FormatVersion: "1.0",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			return s.richTextRepo.Create(ctx, richText)
		}
		return ErrDatabaseError.Wrap(err)
	}

	richText.Content = content
	richText.UpdatedAt = time.Now()
	return s.richTextRepo.Update(ctx, richText)
}

// ==================== 表格文章操作 ====================

// CreateTableArticleInput 创建表格文章输入
type CreateTableArticleInput struct {
	Title       string
	TableID     string
	FolderID    *uint
	OwnerID     uint
	OwnerType   string
	Structure   []map[string]interface{}
	ColumnOrder []string
	Filters     []map[string]interface{}
	Data        []map[string]interface{}
}

// CreateTableArticle 创建表格文章
func (s *Service) CreateTableArticle(ctx context.Context, input *CreateTableArticleInput) (*model.Article, error) {
	// 1. 创建主表
	article, err := s.CreateArticle(ctx, &CreateArticleInput{
		Title:       input.Title,
		ArticleType: model.ArticleTypeTable,
		FolderID:    input.FolderID,
		OwnerID:     input.OwnerID,
		OwnerType:   input.OwnerType,
	})
	if err != nil {
		return nil, err
	}

	// 2. 创建表格结构
	structure := model.JSONArray(input.Structure)
	var columnOrder model.JSONArray
	if len(input.ColumnOrder) > 0 {
		for _, col := range input.ColumnOrder {
			columnOrder = append(columnOrder, map[string]interface{}{"field": col})
		}
	}
	var filters model.JSONArray
	if len(input.Filters) > 0 {
		filters = model.JSONArray(input.Filters)
	}

	tableArticle := &model.TableArticle{
		ArticleID:   article.ID,
		TableID:     input.TableID,
		Structure:   structure,
		ColumnOrder: columnOrder,
		Filters:     filters,
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.tableRepo.Create(ctx, tableArticle); err != nil {
		s.logger.ErrorCtx(ctx, "创建表格结构失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	// 3. 创建行数据
	if len(input.Data) > 0 {
		rows := make([]model.TableArticleRow, len(input.Data))
		for i, rowData := range input.Data {
			idx := i
			rows[i] = model.TableArticleRow{
				ArticleID: article.ID,
				RowData:   model.JSONMap(rowData),
				RowIndex:  &idx,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
		if err := s.tableRowRepo.BatchCreate(ctx, rows); err != nil {
			s.logger.ErrorCtx(ctx, "创建表格行数据失败", zap.Error(err))
			return nil, ErrDatabaseError.Wrap(err)
		}
	}

	return article, nil
}

// TableArticleContent 表格文章完整内容
type TableArticleContent struct {
	Article     *model.Article           `json:"article"`
	Structure   []map[string]interface{} `json:"structure"`
	Data        []map[string]interface{} `json:"data"`
	ColumnOrder []string                 `json:"columnOrder"`
	Filters     []map[string]interface{} `json:"filters"`
}

// GetTableArticleContent 获取表格文章内容
func (s *Service) GetTableArticleContent(ctx context.Context, articleID uint) (*TableArticleContent, error) {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}

	if article.ArticleType != model.ArticleTypeTable {
		return nil, ErrBadRequest.WithMsg("该文章不是表格类型")
	}

	tableArticle, err := s.tableRepo.FindByArticleID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &TableArticleContent{Article: article}, nil
		}
		return nil, ErrDatabaseError.Wrap(err)
	}

	rows, err := s.tableRowRepo.FindByArticleID(ctx, articleID)
	if err != nil {
		return nil, ErrDatabaseError.Wrap(err)
	}

	// 转换行数据
	data := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		data[i] = map[string]interface{}(row.RowData)
	}

	// 转换列顺序
	var columnOrder []string
	for _, col := range tableArticle.ColumnOrder {
		if field, ok := col["field"].(string); ok {
			columnOrder = append(columnOrder, field)
		}
	}

	return &TableArticleContent{
		Article:     article,
		Structure:   []map[string]interface{}(tableArticle.Structure),
		Data:        data,
		ColumnOrder: columnOrder,
		Filters:     []map[string]interface{}(tableArticle.Filters),
	}, nil
}

// GetArticleByTableID 根据tableId获取文章
func (s *Service) GetArticleByTableID(ctx context.Context, tableID string) (*model.Article, error) {
	tableArticle, err := s.tableRepo.FindByTableID(ctx, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound.WithMsg("表格不存在")
		}
		return nil, ErrDatabaseError.Wrap(err)
	}

	return s.GetArticle(ctx, tableArticle.ArticleID)
}

// UpdateTableStructure 更新表格结构
func (s *Service) UpdateTableStructure(ctx context.Context, articleID uint, structure []map[string]interface{}) error {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return err
	}

	if article.ArticleType != model.ArticleTypeTable {
		return ErrBadRequest.WithMsg("该文章不是表格类型")
	}

	tableArticle, err := s.tableRepo.FindByArticleID(ctx, articleID)
	if err != nil {
		return ErrDatabaseError.Wrap(err)
	}

	tableArticle.Structure = model.JSONArray(structure)
	tableArticle.Version++
	tableArticle.UpdatedAt = time.Now()

	return s.tableRepo.Update(ctx, tableArticle)
}

// SaveTableRows 保存表格行数据
func (s *Service) SaveTableRows(ctx context.Context, articleID uint, rowsData []map[string]interface{}) error {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return err
	}

	if article.ArticleType != model.ArticleTypeTable {
		return ErrBadRequest.WithMsg("该文章不是表格类型")
	}

	rows := make([]model.TableArticleRow, len(rowsData))
	for i, rowData := range rowsData {
		idx := i
		rows[i] = model.TableArticleRow{
			ArticleID: articleID,
			RowData:   model.JSONMap(rowData),
			RowIndex:  &idx,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}

	return s.tableRowRepo.ReplaceAll(ctx, articleID, rows)
}

// ==================== Markdown文章操作 ====================

// CreateMarkdownArticleInput 创建Markdown文章输入
type CreateMarkdownArticleInput struct {
	Title     string
	FolderID  *uint
	OwnerID   uint
	OwnerType string
	Content   string
}

// CreateMarkdownArticle 创建Markdown文章
func (s *Service) CreateMarkdownArticle(ctx context.Context, input *CreateMarkdownArticleInput) (*model.Article, error) {
	// 1. 创建主表
	article, err := s.CreateArticle(ctx, &CreateArticleInput{
		Title:       input.Title,
		ArticleType: model.ArticleTypeMarkdown,
		FolderID:    input.FolderID,
		OwnerID:     input.OwnerID,
		OwnerType:   input.OwnerType,
	})
	if err != nil {
		return nil, err
	}

	// 2. 创建Markdown内容
	markdown := &model.MarkdownArticle{
		ArticleID:     article.ID,
		Content:       input.Content,
		HTMLContent:   "", // 可选：在此处或前端渲染HTML
		FormatVersion: "1.0",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.markdownRepo.Create(ctx, markdown); err != nil {
		s.logger.ErrorCtx(ctx, "创建Markdown内容失败", zap.Error(err))
		return nil, ErrDatabaseError.Wrap(err)
	}

	return article, nil
}

// MarkdownArticleContent Markdown文章完整内容
type MarkdownArticleContent struct {
	Article *model.Article `json:"article"`
	Content string         `json:"content"`
}

// GetMarkdownArticleContent 获取Markdown文章内容
func (s *Service) GetMarkdownArticleContent(ctx context.Context, articleID uint) (*MarkdownArticleContent, error) {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}

	if article.ArticleType != model.ArticleTypeMarkdown {
		return nil, ErrBadRequest.WithMsg("该文章不是Markdown类型")
	}

	markdown, err := s.markdownRepo.FindByArticleID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &MarkdownArticleContent{Article: article, Content: ""}, nil
		}
		return nil, ErrDatabaseError.Wrap(err)
	}

	return &MarkdownArticleContent{
		Article: article,
		Content: markdown.Content,
	}, nil
}

// UpdateMarkdownContent 更新Markdown内容
func (s *Service) UpdateMarkdownContent(ctx context.Context, articleID uint, content string) error {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return err
	}

	if article.ArticleType != model.ArticleTypeMarkdown {
		return ErrBadRequest.WithMsg("该文章不是Markdown类型")
	}

	markdown, err := s.markdownRepo.FindByArticleID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			markdown = &model.MarkdownArticle{
				ArticleID:     articleID,
				Content:       content,
				HTMLContent:   "",
				FormatVersion: "1.0",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			return s.markdownRepo.Create(ctx, markdown)
		}
		return ErrDatabaseError.Wrap(err)
	}

	markdown.Content = content
	markdown.UpdatedAt = time.Now()
	return s.markdownRepo.Update(ctx, markdown)
}

// ==================== 文件夹相关操作 ====================

// MoveToFolder 移动文章到指定文件夹
// 注意：folder 有效性验证由应用层/聚合层负责
func (s *Service) MoveToFolder(ctx context.Context, articleID uint, folderID *uint) error {
	article, err := s.GetArticle(ctx, articleID)
	if err != nil {
		return err
	}

	article.FolderID = folderID
	article.UpdatedAt = time.Now()

	if err := s.articleRepo.Update(ctx, article); err != nil {
		s.logger.ErrorCtx(ctx, "移动文章到文件夹失败", zap.Uint("article_id", articleID), zap.Error(err))
		return ErrDatabaseError.Wrap(err)
	}

	s.logger.InfoCtx(ctx, "文章移动成功", zap.Uint("article_id", articleID), zap.Uintp("folder_id", folderID))
	return nil
}

// CountByFolder 统计指定文件夹下的文章数量
func (s *Service) CountByFolder(ctx context.Context, folderID uint) (int64, error) {
	return s.articleRepo.CountByFolderID(ctx, folderID)
}

// ListByFolder 获取指定文件夹下的文章
func (s *Service) ListByFolder(ctx context.Context, folderID uint) ([]model.Article, error) {
	return s.articleRepo.FindByFolderID(ctx, folderID)
}

// ==================== 辅助函数 ====================

func isValidArticleType(t string) bool {
	return t == model.ArticleTypeTable || t == model.ArticleTypeMarkdown || t == model.ArticleTypeRichText
}
