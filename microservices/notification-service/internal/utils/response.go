package utils

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, msg string, data interface{}) { ctx.JSON(http.StatusOK, gin.H{"message": msg, "data": data}) }