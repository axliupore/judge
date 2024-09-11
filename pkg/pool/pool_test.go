package pool

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	p, err := New(3)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Submit(func() {
		fmt.Println("hello pool")
	})
	if err != nil {
		t.Fatal(err)
	}

	p.Wait()
}
