package article

import (
	"context"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"time"
	"webook/_internal/interactive/_internal/repository"
	"webook/pkg/saramax"
)

func NewInteractiveConsumer(client sarama.Client, repo repository.InteractiveRepository) *InteractiveConsumer {
	return &InteractiveConsumer{client: client, repo: repo}
}

type InteractiveConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
}

func (c InteractiveConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		return err
	}
	//应为要在main 循环启动多个start 所以开goroutine
	go func() {
		er := cg.Consume(context.Background(), []string{"read_article"},
			saramax.ConsumerHandle[ReadEvent](c.Consume))
		if er != nil {
			zap.L().Error("退出了消费者", zap.Error(er))
		}
	}()
	return err
}

func (c InteractiveConsumer) Consume(message *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()
	return c.repo.IncrReadCnt(ctx, t.Aid)
}
