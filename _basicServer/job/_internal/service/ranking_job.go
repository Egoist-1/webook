package service

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"go.uber.org/zap"
	"sync"
	"time"
	"webook/_internal/_internal/service"
)

func NewRankingJob() Job {
	return &RankingJob{}
}

type RankingJob struct {
	localLock  sync.Mutex
	client     *rlock.Client
	rankingSvc service.RankingService
	key        string
	//分布式任务锁的过期时间
	timeout time.Duration
	lock    *rlock.Lock
}

func (r RankingJob) Name() string {
	return "ranking_job"
}

func (r RankingJob) Run() error {
	//没有锁 开始抢占
	if r.client == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			Interval: time.Second,
			Max:      3,
		}, time.Second)
		//没抢到锁
		if err != nil {
			return err
		}
		r.lock = lock
		//开始续约
		go func() {
			c, cancel2 := context.WithTimeout(context.Background(), time.Second)
			defer cancel2()
			er := r.lock.Refresh(c)
			if er != nil {
				zap.L().Error("续约失败", zap.Error(er))
			}
		}()
	}

	return r.rankingSvc.TobN(context.Background())
}
