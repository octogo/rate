package meter

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Errorf("New was incorrect, got: %v, want: %v", err, nil)
		}
	}()
	New(time.Second)
}

func TestSpeed(t *testing.T) {
	m := New(time.Second)
	defer m.Close()

	for i := 0; i < 10; i++ {
		m.Tick()
	}
	<-time.After(time.Second)
	v := m.Measure()
	if v != 10 {
		t.Errorf("Speed was incorrect, got: %v, want: %v", v, 10)
	}
}

func ExampleNew() {
	// initialize new speed-o-meter
	m := New(time.Second)

	// remember to free resoures
	defer m.Close()

	// tick 10 times
	for i := 0; i < 10; i++ {
		m.Tick()
	}

	// after one interval the rate should thus read: 10
	<-time.After(time.Second)
	fmt.Println(m.Measure())

	// Output: 10
}
