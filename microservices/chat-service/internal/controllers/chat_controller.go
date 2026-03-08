package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "github.com/google/uuid"
    "chat-service/internal/events"
)

type Hub struct {
    clients map[string]*websocket.Conn
}

type Controller struct {
    hub   *Hub
    kafka *events.Publisher
}

func New(kafka *events.Publisher) *Controller { return &Controller{hub: &Hub{clients: map[string]*websocket.Conn{}}, kafka: kafka} }

var upgrader = websocket.Upgrader{ CheckOrigin: func(r *http.Request) bool { return true } }

func (c *Controller) WS(ctx *gin.Context) {
    conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
    if err != nil { return }
    id := uuid.New().String()
    c.hub.clients[id] = conn
    for {
        var msg map[string]interface{}
        if err := conn.ReadJSON(&msg); err != nil { break }
        _ = c.kafka.Publish(ctx, "chat.message_sent", msg)
        for _, cl := range c.hub.clients { _ = cl.WriteJSON(msg) }
    }
    conn.Close()
    delete(c.hub.clients, id)
}