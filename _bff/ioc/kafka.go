package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"webook/_internal/article/internal/event"
	"webook/_internal/article/internal/event/article"
)

func InitKafka() sarama.Client {
	type Cfg struct {
		Addrs []string `yaml:"addrs"`
	}
	var cfg Cfg
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}
func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func NewConsumers(ac article.InteractiveConsumer) []event.Consumer {
	return []event.Consumer{
		ac,
	}
}
