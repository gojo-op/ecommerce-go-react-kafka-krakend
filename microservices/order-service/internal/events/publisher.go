package events

import (
  "context"
  "encoding/json"
  "os"
  "strings"
  "github.com/IBM/sarama"
)

type Publisher struct{ p sarama.SyncProducer }

func New() *Publisher {
  brokers := os.Getenv("KAFKA_BROKERS")
  if brokers == "" { brokers = "kafka:9092" }
  cfg := sarama.NewConfig()
  cfg.Producer.RequiredAcks = sarama.WaitForLocal
  cfg.Producer.Return.Successes = true
  prod, err := sarama.NewSyncProducer(strings.Split(brokers, ","), cfg)
  if err != nil { return &Publisher{} }
  return &Publisher{ p: prod }
}

func (p *Publisher) Publish(_ context.Context, topic string, payload map[string]interface{}) error {
  if p.p == nil { return nil }
  b, _ := json.Marshal(payload)
  _, _, err := p.p.SendMessage(&sarama.ProducerMessage{ Topic: topic, Value: sarama.ByteEncoder(b) })
  return err
}

func (p *Publisher) Close() error { if p.p != nil { return p.p.Close() }; return nil }