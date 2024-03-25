package dejumbler

import (
	"log"

	"preptar/internal/fileutils"
	"preptar/internal/config"
	"preptar/internal/llama_api"
)

type Dejumbler struct {
	llama *llama.LlamaAPIHandler
	prePrompt string
}

func NewDejumbler(cfg *config.Config) *Dejumbler {
	return &Dejumbler{
		// TODO - set port in config
		llama: llama.NewLlamaAPIHandler("8080"),
		prePrompt: cfg.Prompts.DecodePDF,
	}
}

func (dj *Dejumbler) DejumblePDF(pdfpath string) error {
	paragraphs, err := fileutils.GatherParagraphs(pdfpath)
	if err != nil {
		panic(err)
	}
	for _, paragraph := range paragraphs {
		response, err := dj.llama.MakeRequestAndDecode(
			paragraph,
			dj.prePrompt,
			"MIKE",
			"DECODER",
		)
		if err != nil {
			log.Printf("failed to decode scrambled PDF text : paragraph %s : error %w :", paragraph, err)
			continue
		}
		log.Printf("got a new good output! \nChanged:\n%s\nTo:\n:%s", paragraph, response.Content)
		err = fileutils.AppendNewParagraph("dejumbled-pdf.txt", response.Content)
		if err != nil {
			log.Printf("failed appending dejumbled pdf paragraph to file : %w", err)
		}
	}
	return nil
}