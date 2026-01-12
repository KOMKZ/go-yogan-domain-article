package article

import (
	"net/http"

	"github.com/KOMKZ/go-yogan-framework/errcode"
)

// 错误码模块
const ModuleArticle = 22

// 领域错误定义
var (
	// ErrDatabaseError 数据库错误
	ErrDatabaseError = errcode.Register(errcode.New(
		ModuleArticle, 1001,
		"article",
		"error.article.database_error",
		"数据库操作失败",
		http.StatusInternalServerError,
	))

	// ErrNotFound 文章不存在
	ErrNotFound = errcode.Register(errcode.New(
		ModuleArticle, 1002,
		"article",
		"error.article.not_found",
		"文章不存在",
		http.StatusNotFound,
	))

	// ErrBadRequest 请求参数错误
	ErrBadRequest = errcode.Register(errcode.New(
		ModuleArticle, 1003,
		"article",
		"error.article.bad_request",
		"请求参数错误",
		http.StatusBadRequest,
	))

	// ErrDeleted 文章已删除
	ErrDeleted = errcode.Register(errcode.New(
		ModuleArticle, 1004,
		"article",
		"error.article.deleted",
		"文章已删除",
		http.StatusNotFound,
	))
)
