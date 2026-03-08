package utils

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func Success(c *gin.Context, msg string, data interface{}) { c.JSON(http.StatusOK, gin.H{"message": msg, "data": data}) }
func SuccessWithMeta(c *gin.Context, msg string, data interface{}, meta map[string]interface{}) { c.JSON(http.StatusOK, gin.H{"message": msg, "data": data, "meta": meta}) }
func SuccessResponse(c *gin.Context, status int, msg string, data interface{}) { c.JSON(status, gin.H{"message": msg, "data": data}) }
func ErrorResponse(c *gin.Context, status int, msg string, detail interface{}) { c.JSON(status, gin.H{"error": msg, "detail": detail}) }
func Created(c *gin.Context, msg string, data interface{}) { c.JSON(http.StatusCreated, gin.H{"message": msg, "data": data}) }
func InternalServerError(c *gin.Context, msg string) { c.JSON(http.StatusInternalServerError, gin.H{"error": msg}) }
func UnprocessableEntity(c *gin.Context, msg string, detail interface{}) { c.JSON(http.StatusUnprocessableEntity, gin.H{"error": msg, "detail": detail}) }
func BadRequest(c *gin.Context, msg string) { c.JSON(http.StatusBadRequest, gin.H{"error": msg}) }
func NotFound(c *gin.Context, msg string) { c.JSON(http.StatusNotFound, gin.H{"error": msg}) }

type PaginationQuery struct { Page int `form:"page"` ; PerPage int `form:"per_page"` ; SortBy string `form:"sort_by"` ; Order string `form:"order"` }
func (p *PaginationQuery) GetLimit() int { if p.PerPage <= 0 { return 20 }; return p.PerPage }
func (p *PaginationQuery) GetOffset() int { if p.Page <= 1 { return 0 }; return (p.Page-1) * p.GetLimit() }
func BuildMeta(total, page, perPage int, sortBy, order string) map[string]interface{} { return map[string]interface{}{"total": total, "page": page, "per_page": perPage, "sort_by": sortBy, "order": order} }