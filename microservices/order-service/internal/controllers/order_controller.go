package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "order-service/internal/models"
    "order-service/internal/services"
    utils "order-service/internal/utils"
)

type Controller struct { svc *services.Service }
func New(svc *services.Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Create(ctx *gin.Context) {
    var req models.CreateOrderRequest
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.UnprocessableEntity(ctx, "Invalid request", err.Error()); return }
    o, err := c.svc.Create(ctx, &req)
    if err != nil { utils.InternalServerError(ctx, "Failed to create order"); return }
    utils.Created(ctx, "Order created", o)
}

func (c *Controller) Checkout(ctx *gin.Context) {
    var req models.CheckoutRequest
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.UnprocessableEntity(ctx, "Invalid request", err.Error()); return }
    o, err := c.svc.Checkout(ctx, &req)
    if err != nil { utils.InternalServerError(ctx, "Failed to checkout"); return }
    utils.Created(ctx, "Order created", o)
}

func (c *Controller) Get(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil { utils.BadRequest(ctx, "Invalid id"); return }
    o, err := c.svc.Get(ctx, id)
    if err != nil { utils.NotFound(ctx, "Order not found"); return }
    utils.Success(ctx, "Order", o)
}

func (c *Controller) ListByUser(ctx *gin.Context) {
    userIDStr := ctx.Query("user_id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil { utils.BadRequest(ctx, "Invalid user id"); return }
    var q utils.PaginationQuery
    _ = ctx.ShouldBindQuery(&q)
    items, total, err := c.svc.ListByUser(ctx, userID, q.GetLimit(), q.GetOffset())
    if err != nil { utils.InternalServerError(ctx, "Failed to list orders"); return }
    meta := utils.BuildMeta(int(total), q.Page, q.PerPage, q.SortBy, q.Order)
    utils.SuccessWithMeta(ctx, "Orders", items, meta)
}

func (c *Controller) UpdateStatus(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil { utils.BadRequest(ctx, "Invalid id"); return }
    var req models.UpdateStatusRequest
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.UnprocessableEntity(ctx, "Invalid request", err.Error()); return }
    o, err := c.svc.UpdateStatus(ctx, id, req.Status)
    if err != nil { utils.InternalServerError(ctx, "Failed to update status"); return }
    utils.Success(ctx, "Status updated", o)
}