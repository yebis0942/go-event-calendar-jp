package main

import (
	"github.com/syumai/workers/cloudflare/kv"
)

const (
	cacheNamespace = "CACHE"
	cacheKey       = "calendar"
	cacheTTL       = 300 // seconds
)

// cache represents a wrapper around Cloudflare KV storage for caching calendar data.
type cache struct {
	kv *kv.Namespace
}

// NewCache creates and initializes a new cache instance.
func NewCache() (*cache, error) {
	cacheKV, err := kv.NewNamespace(cacheNamespace)
	if err != nil {
		return nil, err
	}
	return &cache{kv: cacheKV}, nil
}

// Lookup retrieves the cached calendar data.
// It returns the cached string data (if any), and a boolean indicating if the cache was hit.
func (c *cache) Lookup() (string, bool, error) {
	body, err := c.kv.GetString(cacheKey, nil)
	// Cloudflare KV returns "<null>" when the key is not found
	if err != nil || body == "<null>" {
		return "", false, err
	}
	return body, true, nil
}

// Put stores the calendar data in the cache with a 5-minute (300 seconds) TTL.
func (c *cache) Put(body string) error {
	options := &kv.PutOptions{
		ExpirationTTL: cacheTTL,
	}
	return c.kv.PutString(cacheKey, body, options)
}
