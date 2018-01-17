package main

import (
	"testing"
)

func TestTwoPeak(t *testing.T) {
	problem, vars := TwoPeak(8), make([]int32, 8)
	if value := problem.Eval(vars); value != One {
		t.Fatal("should be One", value)
	}

	for i := range vars {
		vars[i] = One
	}
	if value := problem.Eval(vars); value != One {
		t.Fatal("should be one", value)
	}

	for i := 0; i < 4; i++ {
		vars[i] = 0
	}
	if value := problem.Eval(vars); value != One/2 {
		t.Fatal("should be one/2", value)
	}
}

func TestFalsePeak(t *testing.T) {
	problem, vars := FalsePeak(8), make([]int32, 8)
	for i := range vars {
		vars[i] = One
	}
	if value := problem.Eval(vars); value != One {
		t.Fatal("should be one", value)
	}
}
