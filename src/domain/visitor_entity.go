package domain

import (
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/env"
	"time"
)

type VisitorEntity struct {
	ApiKey                    string
	ApiKeyMaxRequestPerSecond int
	IpAddress                 string
	IpMaxRequestPerSecond     int
	BlockedTimePerSecond      int
	LastSeen                  time.Time
	Blocked                   bool
	Unblock                   time.Time
}

func NewVisitorEntity(cfg env.EnvConfig, apiKey, ipAddress string) *VisitorEntity {
	return &VisitorEntity{
		ApiKey:                    apiKey,
		ApiKeyMaxRequestPerSecond: cfg.MaxReqPerSecondToken,
		IpAddress:                 ipAddress,
		IpMaxRequestPerSecond:     cfg.MaxReqPerSecondIp,
		BlockedTimePerSecond:      cfg.BlockedTimePerSecond,
		LastSeen:                  time.Now(),
	}
}

func (v *VisitorEntity) UpdateLastSee() {
	v.LastSeen = time.Now()
}

func (v *VisitorEntity) UpdateUnblock() {
	v.Unblock = time.Now().Add(time.Second * time.Duration(v.BlockedTimePerSecond))
}

func (v *VisitorEntity) Lock() {
	v.Blocked = true
}

func (v *VisitorEntity) Unlock() {
	v.Blocked = false
}
