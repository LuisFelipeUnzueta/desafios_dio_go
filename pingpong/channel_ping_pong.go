package main

import (
	"fmt"
	"sync"
)

func main() {
	// Criar dois canais para coordenação entre as goroutines
	pingChan := make(chan bool)
	pongChan := make(chan bool)

	// WaitGroup para sincronizar o fim das goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine para enviar "ping"
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			fmt.Println("ping")
			pingChan <- true
			<-pongChan
		}
	}()

	// Goroutine para enviar "pong"
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-pingChan
			fmt.Println("pong")
			pongChan <- true
		}
	}()

	// Aguarda a conclusão de ambas as goroutines
	wg.Wait()
	fmt.Println("\nPrograma finalizado!")
}
