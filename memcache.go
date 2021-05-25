package memcache

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrCacheMiss = errors.New("memcache: cache miss")
	ErrNotStored = errors.New("memcache: item not stored")
)

// Item is the unit of memcache gets and sets.
type Item struct {
	// Key is the Item's key (250 bytes maximum).
	Key string
	// Value is the Item's value.
	Value []byte
	// Object is the Item's value for use with a Codec.
	Object interface{}
	// Expiration is the maximum duration that the item will stay
	// in the cache.
	// The zero value means the Item has no expiration time.
	// Subsecond precision is ignored.
	// This is not set when getting items.
	Expiration     time.Duration
	expirationTime time.Time
}

var m = sync.Map{}

// Get gets the item for the given key. ErrCacheMiss is returned for a memcache
// cache miss.
func Get(ctx context.Context, key string) (*Item, error) {
	v, ok := m.Load(key)
	if !ok {
		return nil, ErrCacheMiss
	}

	item := v.(*Item)

	if item.Expiration != 0 && time.Now().After(item.expirationTime) {
		m.Delete(key)
		return nil, ErrCacheMiss
	}
	return item, nil
}

// Set writes the given item, unconditionally.
func Set(ctx context.Context, item *Item) error {
	if item.Expiration < 0 {
		return ErrNotStored
	}
	if item.Expiration != 0 {
		item.expirationTime = time.Now().Add(item.Expiration)
	}
	m.Store(item.Key, item)
	return nil
}
