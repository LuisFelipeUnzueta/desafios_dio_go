package main

import (
	"math"
	"testing"
)

const eps = 1e-9

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) <= eps
}

func TestAdd(t *testing.T) {
	// Arrange
	tests := []struct {
		name string
		a, b float64
		exp  float64
	}{
		{"positive", 2, 3, 5},
		{"negative", -1, -2, -3},
		{"float", 1.5, 2.25, 3.75},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := Add(tc.a, tc.b)
			// Assert
			if !floatEquals(got, tc.exp) {
				t.Fatalf("Add(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.exp)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	// Arrange
	tests := []struct {
		name string
		a, b float64
		exp  float64
	}{
		{"positive", 5, 3, 2},
		{"negative result", 2, 5, -3},
		{"float", 2.5, 1.25, 1.25},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := Subtract(tc.a, tc.b)
			// Assert
			if !floatEquals(got, tc.exp) {
				t.Fatalf("Subtract(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.exp)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	// Arrange
	tests := []struct {
		name string
		a, b float64
		exp  float64
	}{
		{"positive", 3, 4, 12},
		{"zero", 5, 0, 0},
		{"float", 1.5, 2, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got := Multiply(tc.a, tc.b)
			// Assert
			if !floatEquals(got, tc.exp) {
				t.Fatalf("Multiply(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.exp)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	// Arrange
	tests := []struct {
		name    string
		a, b    float64
		exp     float64
		wantErr bool
	}{
		{"normal", 10, 2, 5, false},
		{"float", 7.5, 2.5, 3, false},
		{"zero divisor", 1, 0, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			got, err := Divide(tc.a, tc.b)
			// Assert
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Divide(%v, %v) expected error, got nil", tc.a, tc.b)
				}
				return
			}
			if err != nil {
				t.Fatalf("Divide(%v, %v) unexpected error: %v", tc.a, tc.b, err)
			}
			if !floatEquals(got, tc.exp) {
				t.Fatalf("Divide(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.exp)
			}
		})
	}
}
