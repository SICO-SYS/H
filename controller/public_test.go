package controller

import (
	"testing"

	"github.com/SiCo-Ops/dao/redis"
	"github.com/SiCo-Ops/public"
)

func TestPublicGenerateToken(t *testing.T) {
	key := public.GenerateHexString()
	err := redis.Set(publicPool, key, config.PublicTokenStatus, int64(public.StringToInt(config.PublicTokenExpire)))
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkPublicGenerateToken(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := public.GenerateHexString()
		err := redis.Set(publicPool, key, config.PublicTokenStatus, int64(public.StringToInt(config.PublicTokenExpire)))
		if err != nil {
			b.Error(err)
		}
	}
}

func TestPublicValidateToken(t *testing.T) {
	key := public.GenerateHexString()
	_, code := PublicValidateToken(key)
	if code != 0 {
		t.Error()
	}
}

func BenchmarkPublicValidateToken(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		key := public.GenerateHexString()
		ok, code := PublicValidateToken(key)
		if code != 0 {
			b.Error(code)
		}
		if ok {
			b.Log("duplicate token")
		}
	}
}
