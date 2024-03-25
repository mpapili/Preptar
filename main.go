package main

import (
	"fmt"
	"sync"

	"preptar/internal/config"
	"preptar/internal/pdf_dejumbler"
)

func main() {
	infoChannel := make(chan string)
	cfg := config.DefaultConfig()
	defer close(infoChannel)
	var wg sync.WaitGroup

	// DeJumbler GoRoutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		dejumbler := dejumbler.NewDejumbler(cfg, infoChannel)
		// NOTE - you must run PDF files through training/materials/massage-pdf.sh prior to using this. TXT files only
		err := dejumbler.DejumblePDF("training-materials/arc-first-aid-handbook.txt")
		if err != nil {
			panic(fmt.Errorf("error dejumbling PDF : %v", err))
		}
	}()

	// First-Pass GoRoutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		paragraph := <-infoChannel
		fmt.Printf("HEY! I just got the channel!\n%s", paragraph)
	}()

	wg.Wait()
	close(infoChannel)
}
