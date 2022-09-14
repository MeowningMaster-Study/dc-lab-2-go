package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type duel struct {
	a, b int
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func arena(ctx context.Context, duels <-chan duel, winners chan<- int) {
	for {
		select {
		case <-ctx.Done():
			return
		case duel := <-duels:
			time.Sleep(500 * time.Millisecond)
			winners <- max(duel.a, duel.b)
		}
	}
}

func main() {
	arenasCount := 2 + rand.Intn(6)
	fmt.Printf("Arenas: %d\n", arenasCount)

	monksCount := 40
	monks := make([]int, monksCount)
	for i := 0; i < monksCount; i += 1 {
		monks[i] = rand.Intn(100)
	}

	ctx, cancel := context.WithCancel(context.Background())
	duels := make(chan duel, 100)
	winners := make(chan int, 100)

	for i := 0; i < arenasCount; i += 1 {
		go arena(ctx, duels, winners)
	}

	for monksCount > 1 {
		if len(monks) >= 2 {
			var duel duel
			duel.a = monks[0]
			duel.b = monks[1]
			select {
			case duels <- duel:
				monks = monks[2:]
			default:
			}
		}

		select {
		case winner := <-winners:
			monksCount -= 1
			monks = append(monks, winner)
			fmt.Printf("Duel winner: %d\n", winner)
		default:
		}
	}

	cancel()

	fmt.Printf("Competition winner: %d\n", monks[0])
}
