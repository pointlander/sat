// Copyright 2018 The SAT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestTwoPeak(t *testing.T) {
	size, problem := TwoPeak(8)
	vars := make([]int32, size)
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
	size, problem := FalsePeak(8)
	vars := make([]int32, size)
	for i := range vars {
		vars[i] = One
	}
	if value := problem.Eval(vars); value != One {
		t.Fatal("should be one", value)
	}
}
