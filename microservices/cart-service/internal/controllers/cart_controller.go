package controllers

import (
    "github.com/gin-gonic/gin"
    "cart-service/internal/models"
    "cart-service/internal/services"
    utils "cart-service/internal/utils"
)

type Controller struct { svc *services.Service }

func New(svc *services.Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Get(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    cart, err := c.svc.Get(ctx, userID)
    if err != nil { utils.InternalServerError(ctx, "Failed to get cart"); return }
    utils.Success(ctx, "Cart", cart)
}

func (c *Controller) AddItem(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    var item models.CartItem
    if err := ctx.ShouldBindJSON(&item); err != nil { utils.UnprocessableEntity(ctx, "Invalid request", err.Error()); return }
    cart, err := c.svc.AddItem(ctx, userID, item)
    if err != nil { utils.InternalServerError(ctx, "Failed to add item"); return }
    utils.Success(ctx, "Item added", cart)
}

func (c *Controller) RemoveItem(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    sku := ctx.Param("sku")
    cart, err := c.svc.RemoveItem(ctx, userID, sku)
    if err != nil { utils.InternalServerError(ctx, "Failed to remove item"); return }
    utils.Success(ctx, "Item removed", cart)
}

func (c *Controller) Clear(ctx *gin.Context) {
    userID := ctx.Param("user_id")
    if err := c.svc.Clear(ctx, userID); err != nil { utils.InternalServerError(ctx, "Failed to clear cart"); return }
    utils.Success(ctx, "Cart cleared", nil)
}