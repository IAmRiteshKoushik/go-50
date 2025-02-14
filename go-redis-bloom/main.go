package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type BloomFilterConfig struct {
	Key       string
	Capacity  uint
	ErrorRate float64
}

type UsernameCache struct {
	rdb    *redis.Client
	config BloomFilterConfig
}

func NewUsernameCache(ctx context.Context, rdb *redis.Client,
	config BloomFilterConfig) (*UsernameCache, error) {

	exists, err := rdb.Do(ctx, "BF.EXISTS", config.Key).Bool()
	if err != nil {
		return nil, fmt.Errorf("Error checking Bloom filter: %w", err)
	}

	if !exists {
		err = rdb.Do(ctx, "BF.RESERVE", config.Key, config.ErrorRate, config.Capacity).Err()
		if err != nil {
			return nil, fmt.Errorf("Error creating Bloom filter: %w", err)
		}
		// TODO: Change this into a logging statement
		fmt.Println("Bloom filter created successfully")
	} else {
		fmt.Print("Bloom filter already exists")
	}
	return &UsernameCache{rdb: rdb, config: config}, nil
}

func (uc *UsernameCache) CheckUsername(ctx context.Context, rdb *redis.Client,
	username string) (bool, error) {
	exists, err := rdb.Do(ctx, "BF.CHECK", uc.config.Key, username).Bool()
	if err != nil {
		// TODO: Problems with bloom filters must be reported immediately
		return false, fmt.Errorf("Error checking Bloom Filter: %w", err)
	}

	// if exists is false, then it definitely does not exists
	// if exists is true, then there is chance of false-positive
	return exists, nil
}

func (uc *UsernameCache) AddUsername(ctx context.Context, rdb *redis.Client,
	username string) error {
	ok, err := rdb.Do(ctx, "BF.ADD", uc.config.Key, username).Bool()
	if err != nil {
		// TODO: Problems with bloom filters must be reported immediately via logs
		return fmt.Errorf("Error adding to Bloom filter: %w", err)
	}
	if !ok {
		// TODO: Problems with bloom filters must be reported immediately via logs
		return fmt.Errorf("Element already in set")
	}
	return nil
}
