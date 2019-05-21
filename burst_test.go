package rate

import (
	"fmt"
	"testing"
	"time"
)

func TestNewBurstLimit(t *testing.T) {
	var (
		start = time.Now()
		b1    = NewBurstLimit(1, 1, time.Second)
		b2    = NewBurstLimit(2, 1, time.Second)
	)
	defer b1.Close()
	defer b2.Close()
	b1.Wait()
	if time.Since(start) >= time.Second {
		t.Errorf("NewLimitBurst was incorrect, should have run less than a second")
	} else {
		t.Logf("NewBurstLimit ran: %s", time.Since(start))
	}
	b2.Wait()
	b2.Wait()
	if time.Since(start) >= time.Second {
		t.Errorf("NewLimitBurst was incorrect, should have run less than a second")
	} else {
		t.Logf("NewBurstLimit ran: %s", time.Since(start))
	}
	<-time.After(2 * time.Second)
	start = time.Now()
	b2.Wait()
	b2.Wait()
	if time.Since(start) >= time.Second {
		t.Errorf("NewLimitBurst was incorrect, should have run less than a second")
	} else {
		t.Logf("NewBurstLimit ran: %s", time.Since(start))
	}
}

func ExampleNewBurstLimit() {
	// initialize a new burst-rate-limiter with a burst limit of 10 and a
	// regeneration intervall of 1/s.
	burst := NewBurstLimit(10, 1, time.Second)

	// remember to free resources
	defer burst.Close()

	// burst through the first 10 tokens
	for i := 0; i < 10; i++ {
		// block until the next free token becomes available
		burst.Wait()
	}

	// calling Try() now will not return an error, because there is no free
	// token available in the bucket.
	err := burst.Try()
	fmt.Println(err)
	// Output:
	// too fast
}
