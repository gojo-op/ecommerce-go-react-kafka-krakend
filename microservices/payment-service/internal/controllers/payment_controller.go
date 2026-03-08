package controllers

import (
    "github.com/gin-gonic/gin"
    "payment-service/internal/models"
    "payment-service/internal/services"
    utils "payment-service/internal/utils"
)

type Controller struct { svc *services.Service }
func New(svc *services.Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) CreateIntent(ctx *gin.Context) {
    var req models.CreateIntentRequest
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.UnprocessableEntity(ctx, "Invalid request", err.Error()); return }
    intent, err := c.svc.CreateIntent(ctx, &req)
    if err != nil { utils.InternalServerError(ctx, "Failed to create intent"); return }
    utils.Created(ctx, "Intent created", intent)
}

func (c *Controller) StripeWebhook(ctx *gin.Context) {
    var ev models.WebhookEvent
    if err := ctx.ShouldBindJSON(&ev); err != nil { utils.UnprocessableEntity(ctx, "Invalid payload", err.Error()); return }
    ev.Provider = "stripe"
    if err := c.svc.HandleWebhook(ctx, &ev); err != nil { utils.InternalServerError(ctx, "Webhook error"); return }
    utils.Success(ctx, "ok", nil)
}

func (c *Controller) RazorpayWebhook(ctx *gin.Context) {
    var ev models.WebhookEvent
    if err := ctx.ShouldBindJSON(&ev); err != nil { utils.UnprocessableEntity(ctx, "Invalid payload", err.Error()); return }
    ev.Provider = "razorpay"
    if err := c.svc.HandleWebhook(ctx, &ev); err != nil { utils.InternalServerError(ctx, "Webhook error"); return }
    utils.Success(ctx, "ok", nil)
}