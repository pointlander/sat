package main

import (
	"fmt"
	"math"
)

// Statistic tracks the statistics of some variable
type Statistic struct {
	sum, squared float64
	n            uint64
}

// Add adds the value of a variable to the statistics
func (s *Statistic) Add(x float64) {
	s.sum += x
	s.squared += x * x
	s.n++
}

// Average is the average of the variable
func (s *Statistic) Average() float64 {
	return s.sum / float64(s.n)
}

// Variance is the variance of a variable
func (s *Statistic) Variance() float64 {
	average := s.Average()
	return s.squared/float64(s.n) - average*average
}

// String generates a string for the statistic
func (s *Statistic) String() string {
	return fmt.Sprintf("%f +- %f", s.Average(), math.Sqrt(s.Variance()))
}

// Statistics are statistics gathered from optimize
type Statistics struct {
	Scores   Statistic
	Restarts Statistic
}

// String generates a string for the Statistics
func (s *Statistics) String() string {
	return fmt.Sprintf("Scores: %v\nRestarts: %v",
		&s.Scores, &s.Restarts)
}
