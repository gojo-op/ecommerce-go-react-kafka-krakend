package utils

import (
  "github.com/gin-gonic/gin"
)

func SuccessResponse(ctx *gin.Context, status int, msg string, data interface{}) { ctx.JSON(status, gin.H{"message": msg, "data": data}) }
func ErrorResponse(ctx *gin.Context, status int, msg string, details interface{}) { ctx.JSON(status, gin.H{"error": msg, "details": details}) }