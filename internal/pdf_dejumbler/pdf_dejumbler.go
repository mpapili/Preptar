package dejumbler

import (
	"context"
	"fmt"
	"log"

	"preptar/internal/config"
	"preptar/internal/fileutils"
	"preptar/internal/llama_api"
)

type Dejumbler struct {
	llama     *llama.LlamaAPIHandler
	prePrompt string

	infoChannel chan string
}

func NewDejumbler(cfg *config.Config, infoChannel chan string) *Dejumbler {
	return &Dejumbler{
		infoChannel: infoChannel,
		// TODO - set port in config
		llama:     llama.NewLlamaAPIHandler("8080"),
		prePrompt: cfg.Prompts.DecodePDF,
	}
}

func (dj *Dejumbler) dejumbleParagraph(ctx context.Context, paragraph string) error {
	log.Println("starting to dejumble another paragraph")
	response, err := dj.llama.MakeRequestAndDecode(
		ctx,
		paragraph,
		dj.prePrompt,
		"MIKE",
		"DECODER",
	)
	if err != nil {
		return fmt.Errorf("failed to decode scrambled PDF text : paragraph %s : error %w", paragraph, err)
	}
	log.Println("dejumbler finished processing another paragraph")
	// append to file and send to info channel
	err = fileutils.AppendNewParagraph("dejumbled-pdf.txt", response.Content)
	if err != nil {
		return fmt.Errorf("failed appending dejumbled pdf paragraph to file : %w", err)
	}
	dj.infoChannel <- response.Content
	return nil
}

func (dj *Dejumbler) DejumblePDF(pdfpath string) error {
	paragraphs, err := fileutils.GatherParagraphs(pdfpath)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	for _, paragraph := range paragraphs {
		err = dj.dejumbleParagraph(ctx, paragraph)
		if err != nil {
			log.Printf("error decoding paragraph %s : %v", paragraph, err)
		}
	}
	return nil
}
