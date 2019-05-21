package rate

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLimit(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("NewLimit paniked")
		}
	}()
	NewLimit(time.Second)
}

func TestLimitWait(t *testing.T) {
	l := NewLimit(10 * time.Millisecond)
	defer l.Close()

	start := time.Now()
	l.Wait()
	if time.Since(start) < 10*time.Millisecond {
		t.Errorf("LimitWait was incorrect, should have run for 10ms")
	}
}

func TestLimitTry(t *testing.T) {
	l := NewLimit(10 * time.Millisecond)
	defer l.Close()

	var err error
	<-time.After(11 * time.Millisecond)
	err = l.Try()
	if err != nil {
		t.Errorf("LimitTry was incorrect, got: %v, want: %v", err, nil)
	}
	err = l.Try()
	if err != errTooFast {
		t.Errorf("LimitTry was incorrect, got: %v, want: %v", err, errTooFast)
	}
}

func ExampleNewLimit() {
	// initialize a new rate-limier
	limiter := NewLimit(100 * time.Millisecond)

	// remember to free resources
	defer limiter.Close()

	for i := 0; i < 10; i++ {
		// block until the next free token becomes available
		limiter.Wait()

		// calling Try() immediately afterwards will always fail, because the
		// above Wait() just consumed the token
		err := limiter.Try()
		fmt.Printf("%d, %s\n", i, err)
	}

	// if we wait for the regeneration of a new free token
	<-time.After(time.Second)

	// calling Try() now will not retrun an error, because there is a free
	// token available in the bucket.
	err := limiter.Try()
	fmt.Println(err)
	// Output:
	// 0, too fast
	// 1, too fast
	// 2, too fast
	// 3, too fast
	// 4, too fast
	// 5, too fast
	// 6, too fast
	// 7, too fast
	// 8, too fast
	// 9, too fast
	// <nil>

}
