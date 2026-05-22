package main

import (
	"errors"
	"fmt"
)

// Add returns the sum of a and b.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns the result of a minus b.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of a and b.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns the result of a divided by b or an error when b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	a := Add(2.3, 4.5)
	b := Subtract(5.0, 1.2)
	c := Multiply(3.0, 4.0)
	d, _ := Divide(10.0, 2.0)

	fmt.Println(a, b, c, d)
}
