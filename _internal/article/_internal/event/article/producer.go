package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	IncrReadCnt(ctx context.Context, event ReadEvent) error
}

func NewProducerImpl() Producer {
	return &ProducerKafka{}
}

type ProducerKafka struct {
	producer sarama.SyncProducer
}

func (p *ProducerKafka) IncrReadCnt(ctx context.Context, event ReadEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_article",
		Value: sarama.ByteEncoder(bytes),
	})
	return err
}

type ReadEvent struct {
	Uid int64
	Aid int64
}
