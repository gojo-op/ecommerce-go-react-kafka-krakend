package utils

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

type Meta struct{ Total, Page, PerPage int; SortBy, Order string }
func Success(ctx *gin.Context, msg string, data interface{}) { ctx.JSON(http.StatusOK, gin.H{"message": msg, "data": data}) }
func Created(ctx *gin.Context, msg string, data interface{}) { ctx.JSON(http.StatusCreated, gin.H{"message": msg, "data": data}) }
func BadRequest(ctx *gin.Context, msg string) { ctx.JSON(http.StatusBadRequest, gin.H{"error": msg}) }
func UnprocessableEntity(ctx *gin.Context, msg string, details interface{}) { ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": msg, "details": details}) }
func InternalServerError(ctx *gin.Context, msg string) { ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg}) }
func NotFound(ctx *gin.Context, msg string) { ctx.JSON(http.StatusNotFound, gin.H{"error": msg}) }
func SuccessWithMeta(ctx *gin.Context, msg string, data interface{}, meta Meta) { ctx.JSON(http.StatusOK, gin.H{"message": msg, "data": data, "meta": meta}) }