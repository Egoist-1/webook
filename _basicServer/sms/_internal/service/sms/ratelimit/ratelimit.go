package ratelimit

import (
	"context"
	"sync"
	"time"
	"webook/_basicServer/sms/_internal/service/sms"
)

type ratelimit struct {
	sms.SMS
	ch          chan time.Time
	limit       int
	windowsSize time.Duration
	first       time.Time
	lock        *sync.Mutex
}

func NewRatelimit(SMS sms.SMS) sms.SMS {
	l := 100
	ch := make(chan time.Time, l-1)
	now := time.Now()
	ch <- now
	return &ratelimit{
		SMS:         SMS,
		ch:          ch,
		limit:       l,
		windowsSize: time.Second * 1,
		first:       now,
	}
}

func (m *ratelimit) Send(ctx context.Context, phone string, biz string, templateParam []string) error {
	m.lock.Lock()
	//长度不满不需要限流
	if len(m.ch) < m.limit {
		m.ch <- time.Now()
		return m.Send(ctx, phone, biz, templateParam)
	}
	now := time.Now()
	//如果第一个请求加上窗口大小在now之前需要遍历chan里的数据是否限流
	if m.first.Add(m.windowsSize).Before(time.Now()) {
		for _ = range m.limit {
			t := <-m.ch
			if t.Add(m.windowsSize).Before(now) {
				continue
			}
			m.first = t
			break
		}
	}
	m.lock.Unlock()
	return m.Send(ctx, phone, biz, templateParam)
}
