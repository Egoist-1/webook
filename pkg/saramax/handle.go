package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ConsumerHandle[T any] func(message *sarama.ConsumerMessage, msg T) error

func (c ConsumerHandle[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c ConsumerHandle[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c ConsumerHandle[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			zap.L().Error("反序列化失败", zap.Error(err))
		}
		for _ = range 3 {
			err = c(msg, t)
		}
		if err != nil {
			zap.L().Error("处理消息失败-重试次数上限",
				zap.Error(err),
				zap.String("topic", msg.Topic),
				zap.Int64("partition", int64(msg.Partition)),
				zap.Int64("offset", msg.Offset))
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
