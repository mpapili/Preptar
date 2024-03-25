package config

import (
	"flag"
)

type Config struct {
	LogName       string
	LlamaMainPath string
	DecoderPath   string // path to model to decode PDF nonsense
	Prompts       *Prompts

}

type Prompts struct {
	DecodePDF string
}

func newPrompts() *Prompts {
	return &Prompts{
		DecodePDF: "This is a conversation between a Professor named Mike and a Student " +
			"named Decoder. Professor Mike will give the student a jumbled mess of text. The " +
			"Decoder student's task is to unscramble the text and response with a complete human-readable" +
			"response that retains the context and information of the original text. The student will only " +
			"reply with the answer as a properly-formatted Paragraph; nothing else. The student will ignore references" +
			"to diagrams or figures that are not present.",
	}
}

// DefaultConfig returns a config instance with default values
func DefaultConfig() *Config {
	cfg := Config{
		LogName:       "default.log",
		LlamaMainPath: "./main",
		DecoderPath:   "/home/mike/Downloads/mixtral-8x7b-instruct-v0.1.Q5_K_M.gguf",
		Prompts:       newPrompts(),
	}
	loadConfigFromCommandLine(&cfg)
	return &cfg
}

// LoadConfigFromCommandLine parses command-line arguments and updates the config accordingly
func loadConfigFromCommandLine(cfg *Config) {
	logName := flag.String("logname", cfg.LogName, "the name of the log file")
	llamaPath := flag.String("llama-main-path", cfg.LlamaMainPath, "the path to llama.cpp's 'main' executable")
	decoderPath := flag.String("decoder-path", cfg.DecoderPath, "the path to llama.cpp's 'main' executable")

	// Parse the command-line arguments
	flag.Parse()

	// Update the config values with command-line arguments if provided
	cfg.LogName = *logName
	cfg.LlamaMainPath = *llamaPath
	cfg.DecoderPath = *decoderPath
}
