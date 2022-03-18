package cache

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	N              = 1000
	uploadLifetime = tempScale * 60
)

var getTests = []struct {
	key   string
	value []byte
}{
	{"123456", []byte("test1")},
	{"12", []byte("test2")},
}

func TestTimedGetSet(t *testing.T) {
	c := New()
	for _, tt := range getTests {
		expire := time.Now().Add(1 * time.Second)
		c.Set(tt.key, &Item{tt.value, &expire})
		item, err := c.Get(tt.key)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if !bytes.Equal(tt.value, item.Value) {
			t.Fatalf("cache hit = %v; want %v", item.Value, tt.value)
		}
		time.Sleep(2 * time.Second)
		item, err = c.Get(tt.key)
		if err == nil {
			t.Fatalf("key %s should not be present", tt.key)
		}
	}
}

func TestGetSet(t *testing.T) {
	c := New()
	for _, tt := range getTests {
		expire := time.Now().Add(uploadLifetime)
		c.Set(tt.key, &Item{tt.value, &expire})
		item, err := c.Get(tt.key)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if !bytes.Equal(tt.value, item.Value) {
			t.Fatalf("cache hit = %v; want %v", item.Value, tt.value)
		}
	}
}

func TestGetSetConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	c := New()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			expire := time.Now().Add(uploadLifetime)
			c.Set(fmt.Sprintf("%d", i), &Item{[]byte(fmt.Sprintf("%d", i)), &expire})
			wg.Done()
		}(i)
	}

	wg.Wait()
	for _, i := range ints {
		if _, err := c.Get(fmt.Sprintf("%d", i)); err != nil {
			t.Errorf(err.Error())
		}
	}
}
