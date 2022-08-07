package main

// go build -gcflags '-N -l' -o go_zoo_batch zoo_batch.go

import (
	"fmt"
	"runtime"
	"time"
)

type Zoo struct {
	aarvark  int
	baboon   int
	cat      int
	donkey   int
	elephant int
	fox      int
}

func new() *Zoo {
	return &Zoo{
		aarvark:  1,
		baboon:   1,
		cat:      1,
		donkey:   1,
		elephant: 1,
		fox:      1,
	}
}

func (z *Zoo) ant() int {
	return z.aarvark
}

func (z *Zoo) banana() int {
	return z.baboon
}

func (z *Zoo) tuna() int {
	return z.cat
}

func (z *Zoo) hay() int {
	return z.donkey
}

func (z *Zoo) grass() int {
	return z.elephant
}

func (z *Zoo) mouse() int {
	return z.fox
}

func main() {
	runtime.GOMAXPROCS(1)

	zoo := new()
	sum := 0
	start := time.Now()
	batch := 0

	for time.Since(start).Seconds() < 10 {
		for i := 0; i < 10000; i++ {
			sum += (zoo.ant() +
				zoo.banana() +
				zoo.tuna() +
				zoo.hay() +
				zoo.grass() +
				zoo.mouse())
		}

		batch += 1
	}

	fmt.Printf("%d\n%d\n%f\n", sum, batch, time.Since(start).Seconds())
}
