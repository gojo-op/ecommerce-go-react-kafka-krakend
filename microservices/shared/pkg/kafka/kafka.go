package kafka

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/IBM/sarama"
    cfgpkg "github.com/your-org/microservices/shared/config"
)

const (
    EventProductCreated = "product.created"
    EventProductUpdated = "product.updated"
    EventProductDeleted = "product.deleted"
    EventCartItemAdded = "cart.item_added"
    EventCartItemRemoved = "cart.item_removed"
    EventCartCleared = "cart.cleared"
    EventOrderCreated = "order.created"
    EventOrderStatusChanged = "order.status_changed"
    EventPaymentProcessed = "payment.processed"
    EventPaymentFailed = "payment.failed"
    EventPaymentRefunded = "payment.refunded"
    EventChatMessageSent = "chat.message_sent"
)

type KafkaProducer struct{
    p sarama.SyncProducer
}

func NewKafkaProducer(cfg *cfgpkg.Config) (*KafkaProducer, error) {
    return &KafkaProducer{p: nil}, nil
}

type wire struct{
    Type string `json:"type"`
    Data map[string]interface{} `json:"data"`
    Metadata map[string]interface{} `json:"metadata"`
    Ts int64 `json:"ts"`
}

func (p *KafkaProducer) PublishEvent(ctx context.Context, eventType string, data map[string]interface{}, meta map[string]interface{}) error {
    if meta == nil { meta = map[string]interface{}{} }
    w := wire{ Type: eventType, Data: data, Metadata: meta, Ts: time.Now().UnixMilli() }
    b, err := json.Marshal(w)
    if err != nil { return err }
    if p.p == nil { return nil }
    msg := &sarama.ProducerMessage{ Topic: eventType, Value: sarama.ByteEncoder(b) }
    _, _, err = p.p.SendMessage(msg)
    return err
}

func (p *KafkaProducer) Close() error { if p.p != nil { return p.p.Close() }; return nil }

type EventMessage struct{
    Type string
    Data map[string]interface{}
    Metadata map[string]interface{}
}

type MessageHandler interface{
    HandleMessage(msg *EventMessage) error
}

type consumerGroupHandler struct{ h MessageHandler }
func (cg consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (cg consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (cg consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var w wire
        _ = json.Unmarshal(msg.Value, &w)
        em := &EventMessage{ Type: msg.Topic, Data: w.Data, Metadata: w.Metadata }
        _ = cg.h.HandleMessage(em)
        sess.MarkMessage(msg, "")
    }
    return nil
}

type KafkaConsumer struct{
    g sarama.ConsumerGroup
    cancel context.CancelFunc
}

func NewKafkaConsumer(cfg *cfgpkg.Config, groupID string, topics []string, handler MessageHandler) (*KafkaConsumer, error) {
    c := sarama.NewConfig()
    c.Version = sarama.V3_6_0_0
    g, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, groupID, c)
    if err != nil { return nil, err }
    ctx, cancel := context.WithCancel(context.Background())
    h := consumerGroupHandler{ h: handler }
    go func() {
        for {
            if err := g.Consume(ctx, topics, h); err != nil {
                time.Sleep(time.Second)
            }
            if ctx.Err() != nil {
                return
            }
        }
    }()
    return &KafkaConsumer{ g: g, cancel: cancel }, nil
}

func (c *KafkaConsumer) Close() error { if c.cancel != nil { c.cancel() }; if c.g != nil { return c.g.Close() }; return nil }
