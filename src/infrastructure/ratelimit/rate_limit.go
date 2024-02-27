package ratelimit

import (
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/domain"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/env"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type RateLimiterParams struct {
	Limiter *rate.Limiter
	Visitor domain.VisitorEntity
}

type RateLimiterService struct {
	VisitorControl map[string]*RateLimiterParams
	Mtx            sync.Mutex
	Cfg            env.EnvConfig
	repository     domain.VisitorRepository
}

func NewRateLimiterService(cfg env.EnvConfig, repository domain.VisitorRepository) *RateLimiterService {
	return &RateLimiterService{
		VisitorControl: make(map[string]*RateLimiterParams),
		Cfg:            cfg,
		repository:     repository,
	}
}

func (rls *RateLimiterService) GetRateLimiterParams(apiKey, ip string) *RateLimiterParams {
	rls.Mtx.Lock()
	defer rls.Mtx.Unlock()

	limiterParams, exists := rls.VisitorControl[ip]
	if !exists {
		return rls.NewRateLimiter(apiKey, ip)
	}
	if rls.isBlockedTimeExpired(limiterParams) {
		limiterParams.Visitor.Unlock()
	}

	limiterParams.Visitor.UpdateLastSee()
	rls.repository.SaveOrUpdate(limiterParams.Visitor)

	return limiterParams
}

func (rls *RateLimiterService) isBlockedTimeExpired(limiterParams *RateLimiterParams) bool {
	return limiterParams.Visitor.Blocked && time.Now().After(limiterParams.Visitor.Unblock)
}

func (rls *RateLimiterService) NewRateLimiter(apiKey string, ip string) *RateLimiterParams {
	visitor := domain.NewVisitorEntity(rls.Cfg, apiKey, ip)
	rls.repository.SaveOrUpdate(*visitor)
	var limiterParams *RateLimiterParams

	if len(apiKey) > 0 {
		limiterParams = &RateLimiterParams{
			Limiter: rate.NewLimiter(rate.Every(time.Second), visitor.ApiKeyMaxRequestPerSecond),
			Visitor: *visitor,
		}
	} else {
		limiterParams = &RateLimiterParams{
			Limiter: rate.NewLimiter(rate.Every(time.Second), visitor.IpMaxRequestPerSecond),
			Visitor: *visitor,
		}
	}
	rls.VisitorControl[ip] = limiterParams
	return limiterParams
}

func (rls *RateLimiterService) CleanupVisitors() {
	rls.Mtx.Lock()
	defer rls.Mtx.Unlock()

	for ip, limiterParams := range rls.VisitorControl {
		if rls.isBlockedTimeExpired(limiterParams) {
			delete(rls.VisitorControl, ip)
		}
	}
}

func (rls *RateLimiterService) LockVisitor(ip string) {
	limiterParams, _ := rls.VisitorControl[ip]
	limiterParams.Visitor.Lock()
	limiterParams.Visitor.UpdateUnblock()
	rls.repository.SaveOrUpdate(limiterParams.Visitor)
}
