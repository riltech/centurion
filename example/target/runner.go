package main

import (
	"sync"

	"github.com/riltech/centurion/example"
)

const target = "ec2-35-159-46-128.eu-central-1.compute.amazonaws.com"

func main() {
	attacker, defender := example.NewAttacker(target), example.NewDefender(target)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		attacker.Start()
		wg.Done()
	}()
	defender.Start()
}
