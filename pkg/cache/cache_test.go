package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {

	c, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if ok := c.Set("1", 1); !ok {
		t.Fatal("failed to set 1")
	}

	if ok := c.SetTime("2", 2, 1*time.Second); !ok {
		t.Fatal("failed to set 2")
	}

	c.Wait()

	if _, ok := c.Get("1"); !ok {
		t.Fatal("failed to get 1")
	}

	time.Sleep(1 * time.Second)

	if _, ok := c.Get("2"); ok {
		t.Fatal("set time err")
	}
}
