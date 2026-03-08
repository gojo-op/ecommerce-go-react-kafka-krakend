package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "product-service/internal/models"
    "product-service/internal/services"
    utils "product-service/internal/utils"
)

type Controller struct { svc *services.Service }

func New(svc *services.Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Create(ctx *gin.Context) {
    var req models.CreateProductRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        utils.UnprocessableEntity(ctx, "Invalid request", err.Error())
        return
    }
    p, err := c.svc.Create(ctx, &req)
    if err != nil { utils.InternalServerError(ctx, "Failed to create product"); return }
    utils.Created(ctx, "Product created", p)
}

func (c *Controller) Update(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil { utils.BadRequest(ctx, "Invalid id"); return }
    var req models.UpdateProductRequest
    if err := ctx.ShouldBindJSON(&req); err != nil { utils.UnprocessableEntity(ctx, "Invalid request", err.Error()); return }
    p, err := c.svc.Update(ctx, id, &req)
    if err != nil { utils.InternalServerError(ctx, "Failed to update product"); return }
    utils.Success(ctx, "Product updated", p)
}

func (c *Controller) UpdateStock(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil { utils.BadRequest(ctx, "Invalid id"); return }
    var payload struct{ Quantity int `json:"quantity" binding:"required"` }
    if err := ctx.ShouldBindJSON(&payload); err != nil { utils.UnprocessableEntity(ctx, "Invalid request", err.Error()); return }
    p, err := c.svc.UpdateStock(ctx, id, payload.Quantity)
    if err != nil { utils.InternalServerError(ctx, "Failed to update stock"); return }
    utils.Success(ctx, "Stock updated", p)
}

func (c *Controller) Delete(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil { utils.BadRequest(ctx, "Invalid id"); return }
    if err := c.svc.Delete(ctx, id); err != nil { utils.InternalServerError(ctx, "Failed to delete product"); return }
    ctx.Status(http.StatusNoContent)
}

func (c *Controller) Get(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil { utils.BadRequest(ctx, "Invalid id"); return }
    p, err := c.svc.Get(ctx, id)
    if err != nil { utils.NotFound(ctx, "Product not found"); return }
    utils.Success(ctx, "Product", p)
}

func (c *Controller) GetBySKU(ctx *gin.Context) {
    sku := ctx.Param("sku")
    p, err := c.svc.GetBySKU(ctx, sku)
    if err != nil { utils.NotFound(ctx, "Product not found"); return }
    utils.Success(ctx, "Product", p)
}

func (c *Controller) List(ctx *gin.Context) {
    var q utils.PaginationQuery
    if err := ctx.ShouldBindQuery(&q); err != nil {}
    items, total, err := c.svc.List(ctx, q.GetLimit(), q.GetOffset())
    if err != nil { utils.InternalServerError(ctx, "Failed to list products"); return }
    meta := utils.BuildMeta(int(total), q.Page, q.PerPage, q.SortBy, q.Order)
    utils.SuccessWithMeta(ctx, "Products", items, meta)
}