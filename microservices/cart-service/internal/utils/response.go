package utils

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, msg string, data interface{}) { ctx.JSON(http.StatusOK, gin.H{"message": msg, "data": data}) }
func UnprocessableEntity(ctx *gin.Context, msg string, details interface{}) { ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": msg, "details": details}) }
func InternalServerError(ctx *gin.Context, msg string) { ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg}) }