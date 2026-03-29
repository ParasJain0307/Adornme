package otp

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Store struct {
	Client *redis.Client
}

func NewStore(addr string) *Store {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Store{Client: rdb}
}

// Save OTP
func (s *Store) Save(userID, otpType, otp string) error {
	key := "otp:" + otpType + ":" + userID
	return s.Client.Set(ctx, key, otp, 5*time.Minute).Err()
}

// Get OTP
func (s *Store) Get(userID, otpType string) (string, error) {
	key := "otp:" + otpType + ":" + userID
	return s.Client.Get(ctx, key).Result()
}

// Delete OTP
func (s *Store) Delete(userID, otpType string) {
	key := "otp:" + otpType + ":" + userID
	s.Client.Del(ctx, key)
}
