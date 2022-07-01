package utils

import "fmt"

func ProgressDisplay(total int) (statusCh chan int, doneCh chan struct{}) {
	statusCh = make(chan int, 10)
	doneCh = make(chan struct{})
	go func() {
		for ch := range statusCh {
			fmt.Printf("\rProcessing: %d/%d", ch, total)
		}
		fmt.Printf("\rProcessing Complete!\n")
		close(doneCh)
	}()
	return
}
