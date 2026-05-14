package main

import "fmt"

func main() {
	celsius := 100.0
	kelvin := celsius + 273.15

	fmt.Printf("Temperatura de ebulição da água: %.2f °C = %.2f K\n", celsius, kelvin)
}
