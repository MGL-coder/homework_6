package homework_6

import (
	"errors"
	"sync"
)

func Execute(tasks []func() error, N int, E int) error {
	if N < 1 || E < 1 || len(tasks) == 0 {
		return errors.New("incorrect arguments: N, E, len(tasks) must be positive")
	}

	var (
		ch = make(chan struct{}, N)
		done = make(chan struct{})
		mu sync.Mutex
	)

	for index, task := range tasks {
		ch <- struct{}{}
		mu.Lock()
		if E < 1 {
			break
		}
		mu.Unlock()

		i, t := index, task
		go func() {
			//fmt.Println("Running task", i)
			err := t()
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				E--
				//fmt.Println("Error detected")
			}
			<-ch
			if E == 0 || i == len(tasks) - 1 {
				close(done)
			}
		}()
	}
	<-done
	if E < 1 {
		return errors.New("exceeded the permissible number of errors")
	}
	return nil
}
