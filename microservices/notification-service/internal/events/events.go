package events

import (
  "context"
  "encoding/json"
  "os"
  "strings"
  "github.com/IBM/sarama"
)

type EventMessage struct{
  Type string
  Data map[string]interface{}
  Metadata map[string]interface{}
}

type MessageHandler interface{
  HandleMessage(msg *EventMessage) error
}

type KafkaConsumer struct{
  g sarama.ConsumerGroup
  cancel context.CancelFunc
}

func Start(ctx context.Context, topics []string, h MessageHandler) (*KafkaConsumer, error) {
  brokers := os.Getenv("KAFKA_BROKERS")
  if brokers == "" { brokers = "kafka:9092" }
  cfg := sarama.NewConfig()
  cfg.Version = sarama.V3_6_0_0
  group, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), "notification-service", cfg)
  if err != nil { return nil, err }
  cctx, cancel := context.WithCancel(ctx)
  handler := &cgHandler{ h: h }
  go func(){ for { if err := group.Consume(cctx, topics, handler); err != nil { if cctx.Err() != nil { return } } } }()
  return &KafkaConsumer{ g: group, cancel: cancel }, nil
}

func (c *KafkaConsumer) Close() error { if c.cancel != nil { c.cancel() }; if c.g != nil { return c.g.Close() }; return nil }

type cgHandler struct{ h MessageHandler }
func (cg *cgHandler) Setup(sarama.ConsumerGroupSession) error { return nil }
func (cg *cgHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (cg *cgHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
  for msg := range claim.Messages() {
    var payload map[string]interface{}
    _ = json.Unmarshal(msg.Value, &payload)
    em := &EventMessage{ Type: msg.Topic, Data: payload, Metadata: map[string]interface{}{"partition": msg.Partition, "offset": msg.Offset} }
    _ = cg.h.HandleMessage(em)
    sess.MarkMessage(msg, "")
  }
  return nil
}