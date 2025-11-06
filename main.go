package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type CancelOnErrorGroup struct {
	handles    []func(context.Context) error
	firstError error
}

func NewCancelOnErrorGroup(num int) CancelOnErrorGroup {
	return CancelOnErrorGroup{
		handles: make([]func(context.Context) error, 0, num),
	}
}

func (ceg *CancelOnErrorGroup) Go(handle func(ctx context.Context) error) {
	ceg.handles = append(ceg.handles, handle)
}

func (ceg *CancelOnErrorGroup) Wait(ctx context.Context) error {

	if len(ceg.handles) == 0 {
		return ceg.firstError
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, handle := range ceg.handles {
		wg.Go(func() {
			err := handle(ctx)
			if err != nil {
				ceg.firstError = err
				fmt.Println("ERROR!")
				cancel()
			}
		})

	}

	wg.Wait()
	return ceg.firstError
}

func main() {

	// create a cancel on error group - it's like a wait group
	// but when any of the go routines error, the context is cancelled
	ceg := NewCancelOnErrorGroup(2)

	// create a happy function which is doing work and checking ctx
	ceg.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				fmt.Println("Looper")
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	})

	// create a sad function that will do some work and then die
	ceg.Go(func(ctx context.Context) error {
		time.Sleep(time.Duration(5) * time.Second)
		return errors.New("sorry, I let you down")
	})

	fmt.Println("Starting")

	// now wait until our goroutines are done
	err := ceg.Wait(context.Background())
	if err != nil {
		fmt.Printf("There was an error, %v\n", err)
	}

	fmt.Println("Finished")
}
