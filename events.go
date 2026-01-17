package article

import "github.com/KOMKZ/go-yogan-framework/event"

// 事件名称常量
const (
	EventArticleCreated        = "article.created"
	EventArticleDeleted        = "article.deleted"
	EventArticleMoved          = "article.moved"
	EventArticleContentUpdated = "article.content.updated"
)

// ArticleCreatedEvent 文章创建事件
type ArticleCreatedEvent struct {
	event.BaseEvent
	ArticleID uint
	FolderID  *uint
}

// ArticleDeletedEvent 文章删除事件
type ArticleDeletedEvent struct {
	event.BaseEvent
	ArticleID uint
	FolderID  *uint // 删除前所属的文件夹
}

// ArticleMovedEvent 文章移动事件
type ArticleMovedEvent struct {
	event.BaseEvent
	ArticleID   uint
	OldFolderID *uint
	NewFolderID *uint
}

// NewArticleCreatedEvent 创建文章创建事件
func NewArticleCreatedEvent(articleID uint, folderID *uint) *ArticleCreatedEvent {
	return &ArticleCreatedEvent{
		BaseEvent: event.NewEvent(EventArticleCreated),
		ArticleID: articleID,
		FolderID:  folderID,
	}
}

// NewArticleDeletedEvent 创建文章删除事件
func NewArticleDeletedEvent(articleID uint, folderID *uint) *ArticleDeletedEvent {
	return &ArticleDeletedEvent{
		BaseEvent: event.NewEvent(EventArticleDeleted),
		ArticleID: articleID,
		FolderID:  folderID,
	}
}

// NewArticleMovedEvent 创建文章移动事件
func NewArticleMovedEvent(articleID uint, oldFolderID, newFolderID *uint) *ArticleMovedEvent {
	return &ArticleMovedEvent{
		BaseEvent:   event.NewEvent(EventArticleMoved),
		ArticleID:   articleID,
		OldFolderID: oldFolderID,
		NewFolderID: newFolderID,
	}
}

// ArticleContentUpdatedEvent 文章内容更新事件
type ArticleContentUpdatedEvent struct {
	event.BaseEvent
	ArticleID   uint
	ContentType string // markdown, rich_text, table
}

// NewArticleContentUpdatedEvent 创建文章内容更新事件
func NewArticleContentUpdatedEvent(articleID uint, contentType string) *ArticleContentUpdatedEvent {
	return &ArticleContentUpdatedEvent{
		BaseEvent:   event.NewEvent(EventArticleContentUpdated),
		ArticleID:   articleID,
		ContentType: contentType,
	}
}

// ============== CacheInvalidator 接口实现 ==============

// CacheArgs 返回缓存失效参数（ArticleDeletedEvent）
func (e *ArticleDeletedEvent) CacheArgs() []any {
	return []any{e.ArticleID}
}

// CacheArgs 返回缓存失效参数（ArticleContentUpdatedEvent）
func (e *ArticleContentUpdatedEvent) CacheArgs() []any {
	return []any{e.ArticleID}
}
