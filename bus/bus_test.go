//
// bus_test.go
// Copyright (C) 2025 veypi <i@veypi.com>
// 2025-03-05 17:46
// Distributed under terms of the MIT license.
//

package bus

import (
	"testing"
)

func TestBusOnAndEmit(t *testing.T) {
	testData := "test message"
	var bus = New[string]()
	called := 0
	cancelFn := bus.On(func(data string) {
		if data != testData {
			t.Errorf("Expected data to be %s, got %s", testData, data)
		}
		called++
	})
	bus.Emit(testData)
	cancelFn()
	bus.Emit(testData)
	if called != 1 {
		t.Errorf("Expected the callback function to be called: %d", called)
	}
}

func TestBusMultipleCallbacks(t *testing.T) {
	testData := 42
	var bus = New[int]()

	callCount := 0
	for i := 0; i < 3; i++ {
		bus.On(func(data int) {
			if data != testData {
				t.Errorf("Expected data to be %d, got %d", testData, data)
			}
			callCount++
		})
	}
	bus.Emit(testData)
	if callCount != 3 {
		t.Errorf("Expected 3 callbacks to be called, got %d", callCount)
	}
}
