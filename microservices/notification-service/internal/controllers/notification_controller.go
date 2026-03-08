package controllers

import (
    "github.com/gin-gonic/gin"
    "notification-service/internal/services"
    utils "notification-service/internal/utils"
)

type Controller struct { svc *services.Service }

func New(svc *services.Service) *Controller { return &Controller{ svc: svc } }

func (c *Controller) List(ctx *gin.Context) {
    userID := ctx.Query("user_id")
    utils.Success(ctx, "Notifications", c.svc.List(userID))
}