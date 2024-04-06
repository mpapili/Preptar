package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"preptar/internal/config"
	"preptar/internal/llama_api"
	"preptar/internal/pdf_dejumbler"
	"preptar/internal/fileutils"
)

func SubstringAfter(s, substr string) (string, error) {
	index := strings.Index(s, substr)
	if index == -1 {
		// Substring not found
		return "", fmt.Errorf("could not find text after %s in %s", substr, s)
	}
	// Adjust the position to start after the substring
	pos := index + len(substr)
	if pos >= len(s) {
		// Substring is at the end of string
		return "", fmt.Errorf("nothing occurs after %s in %s", substr, s)
	}
	return s[pos:], nil
}

func main() {
	log.Println("Before starting, start the decoder llama cpp server on port 8080 and the Questioner server on 8081")
	infoChannel := make(chan string)
	answererChannel := make(chan string)
	reviewerChannel := make(chan string)
	cfg := config.DefaultConfig()
	defer close(infoChannel)
	defer close(answererChannel)
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

	// Questioner GoRoutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		questionLlama := llama.NewLlamaAPIHandler(cfg.Ports.QuestionerPort) // TODO - new port to add to config
		for {
			paragraph := <-infoChannel
			fmt.Printf("\nQUESTIONER: I just got work from my channel!\n%s\n", paragraph)
			ctx := context.Background()
			resp, err := questionLlama.MakeRequestAndDecode(
				ctx,
				paragraph,
				cfg.Prompts.Questioner,
				"MIKE",
				"QUESTIONER",
			)
			if err != nil {
				log.Printf("failed to generate a question and answer from paragraph : %s : %w", paragraph, err)
			}
			fmt.Printf("\nQUESTIONER: Got a question:\n%s\n", resp.Content)
			answererChannel <- fmt.Sprintf("%s\n<QUESTION>: %s\n", paragraph, resp.Content)
		}
	}()

	// Answerer GoRoutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		answerLlama := llama.NewLlamaAPIHandler(cfg.Ports.AnswererPort) // TODO - new port to add to config
		for {
			prompt := <-answererChannel
			fmt.Printf("\nANSWERER: I just got work from my channel!\n%s", prompt)
			ctx := context.Background()
			resp, err := answerLlama.MakeRequestAndDecode(
				ctx,
				prompt,
				cfg.Prompts.Answerer,
				"MIKE",
				"ANSWERER",
			)
			if err != nil {
				log.Printf("failed to generate a question and answer from prompt : %s : %v", prompt, err)
			}
			fmt.Printf("\nANSWERER: Got an answer:\n%s\n", resp.Content)
			question, err := SubstringAfter(prompt, "<QUESTION>")
			if err != nil {
				log.Printf("failed to isolate question : %w", err)
			}
			answer := resp.Content
			fmt.Printf("\nANSWERER: Your final Question and Answer Combo:\nQuesiton: %s\nAnswer: %s\n\n", question, answer)
			reviewerChannel <- fmt.Sprintf("Question: %s\nAnswer:%s", question, answer)
		}
	}()

	// Peer Reviewer GoRoutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		reviewerLlama := llama.NewLlamaAPIHandler(cfg.Ports.PeerReviewerPort) // TODO - new port to add to config
		for {
			questionAndAnswer := <-reviewerChannel
			fmt.Printf("\nPeerReviewer: I just got work from my channel!\n%s", questionAndAnswer)
			ctx := context.Background()
			resp, err := reviewerLlama.MakeRequestAndDecode(
				ctx,
				questionAndAnswer,
				cfg.Prompts.PeerReviewer,
				"MIKE",
				"REVIEWER",
			)
			if err != nil {
				log.Printf("failed to generate a question and answer from prompt : %s : %w", questionAndAnswer, err)
			}
			safeDetermination := resp.Content
			log.Println("\nPEER REVIEWER: I have determined that the question+answer combo is :", safeDetermination)
			if len(safeDetermination) > 10 {
				log.Printf("error - peer reviewer gave an invalid response : %s", safeDetermination)
			} else {
				if strings.Contains(safeDetermination, "unsafe") {
					log.Println("Unsafe! This question/answer combo will NOT be saved")
				} else if strings.Contains(safeDetermination, "safe") {
					log.Println("Safe! Saving this combo")
					err = fileutils.AppendNewParagraph("qa_combos.txt", fmt.Sprintf("%s\n", questionAndAnswer))
					if err != nil {
						log.Printf("failed to append QA combo to file : %w", err)
					}
				}
			}
		}
	}()

	wg.Wait()
	close(infoChannel)
}
