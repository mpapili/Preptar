package main

import (
	"fmt"
	"log"

	"preptar/internal/config"
	"preptar/internal/pdf_dejumbler"
)

func main() {
	cfg := config.DefaultConfig()
	dejumbler := dejumbler.NewDejumbler(cfg)
	// NOTE - you must run PDF files through training/materials/massage-pdf.sh prior to using this. TXT files only
	err := dejumbler.DejumblePDF("training-materials/arc-first-aid-handbook.txt")
	if err != nil {
		panic(fmt.Errorf("error dejumbling PDF : %v", err))
	}
}
