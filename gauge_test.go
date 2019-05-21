package rate

import (
	"fmt"
	"testing"
	"time"
)

func TestNewGauge(t *testing.T) {
	var (
		start = time.Now()
		g1    = NewGauge(1, 1, time.Second)
		g2    = NewGauge(2, 1, time.Second)
	)
	defer g1.Close()
	defer g2.Close()
	g1.Wait()
	if time.Since(start) >= time.Second {
		t.Errorf("NewGauge was incorrect, should have run less than a second")
	} else {
		t.Logf("NewGauge ran: %s", time.Since(start))
	}
	g2.Wait()
	g2.Wait()
	if time.Since(start) >= time.Second {
		t.Errorf("NewGauge was incorrect, should have run less than a second")
	} else {
		t.Logf("NewGauge ran: %s", time.Since(start))
	}
	<-time.After(2 * time.Second)
	start = time.Now()
	g2.Wait()
	g2.Wait()
	if time.Since(start) >= time.Second {
		t.Errorf("NewGauge was incorrect, should have run less than a second")
	} else {
		t.Logf("NewGauge ran: %s", time.Since(start))
	}
}

func ExampleNewGauge() {
	// imitialize a new gauge
	gauge := NewGauge(10, 1, time.Second)

	// remember to free resources
	defer gauge.Close()

	// burst through the first 10 tokens
	for i := 0; i < 10; i++ {
		// block until the next free token becomes available
		gauge.Wait()
	}

	// calling Try() now will not return an error, because there is no free
	// token available in the bucket.
	err := gauge.Try()
	fmt.Println(err)

	// we should now measure 10/s ticks
	<-time.After(time.Second)
	fmt.Println("Ticks per second:", gauge.Measure())

	// Output:
	// too fast
	// Ticks per second: 10
}
