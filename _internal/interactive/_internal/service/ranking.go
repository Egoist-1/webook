package service

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"math"
	"time"
	service2 "webook/_internal/article/_internal/service"
	"webook/_internal/interactive/_internal/domain"
)

type RankingService interface {
	TobN(ctx context.Context) error
}

func NewRankingService(asvc service2.ArticleService, isvc InteractiveService) RankingService {
	return &rankingService{
		artSvc: asvc,
		intr:   isvc,
		batch:  3,
		score: func(readCnt int64, ctime int64) float64 {
			ti := time.UnixMilli(ctime)
			time := time.Since(ti).Seconds()
			return float64(readCnt) / math.Pow(float64(time), 1.5)
		},
	}
}

type rankingService struct {
	artSvc service2.ArticleService
	intr   InteractiveService
	//每次处理多少
	capacity int
	batch    int
	score    func(readCnt int64, ctime int64) float64
	biz      string
}

func (r rankingService) TobN(ctx context.Context) error {
	topN := queue.NewConcurrentPriorityQueue[domain.Interactive](r.capacity, func(src domain.Interactive, dst domain.Interactive) int {
		srcScore := r.score(src.ReadCnt, src.Ctime)
		dstScore := r.score(dst.ReadCnt, dst.Ctime)
		if srcScore > dstScore {
			return -1
		} else if srcScore == dstScore {
			return 0
		}
		return -1
	})
	now := time.Now()
	offset := 0

	for {
		intrs, err := r.intr.TopN(ctx, now, r.batch)
		if err != nil {
			return err
		}
		if len(intrs) < r.batch {
			break
		}
		for _, val := range intrs {
			topN.Enqueue(val)
		}
		offset += r.batch
	}
	ids := make([]int64, r.capacity)
	for _ = range r.capacity {
		intr, _ := topN.Dequeue()
		ids = append(ids, intr.BizId)
	}
	//TODO 存储 redis
	_, err := r.artSvc.TopN(ctx, ids)
	return err
}
