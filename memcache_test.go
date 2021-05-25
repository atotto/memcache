package memcache

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestMemcache(t *testing.T) {
	ctx := context.Background()

	if err := Set(ctx, &Item{
		Key:        "hello",
		Value:      []byte("world"),
		Expiration: time.Second,
	}); err != nil {
		t.Fatal(err)
	}

	item, err := Get(ctx, "hello")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(item.Value, []byte("world")) {
		t.Fatalf("want world")
	}

	for {
		_, err := Get(ctx, "hello")
		if err == ErrCacheMiss {
			break
		}
	}
}
