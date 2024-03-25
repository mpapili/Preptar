package llama

type LlamaCppRequest struct {
	Stream           bool     `json:"stream"`
	NPredict         int      `json:"n_predict"`
	Temperature      float64  `json:"temperature"`
	Stop             []string `json:"stop"`
	RepeatLastN      int      `json:"repeat_last_n"`
	RepeatPenalty    float64  `json:"repeat_penalty"`
	TopK             int      `json:"top_k"`
	TopP             float64  `json:"top_p"`
	MinP             float64  `json:"min_p"`
	TfsZ             int      `json:"tfs_z"`
	TypicalP         int      `json:"typical_p"`
	PresencePenalty  int      `json:"presence_penalty"`
	FrequencyPenalty int      `json:"frequency_penalty"`
	Mirostat         int      `json:"mirostat"`
	MirostatTau      int      `json:"mirostat_tau"`
	MirostatEta      float64  `json:"mirostat_eta"`
	Grammar          string   `json:"grammar"`
	NProbs           int      `json:"n_probs"`
	MinKeep          int      `json:"min_keep"`
	ImageData        []any    `json:"image_data"`
	CachePrompt      bool     `json:"cache_prompt"`
	APIKey           string   `json:"api_key"`
	SlotID           int      `json:"slot_id"`
	Prompt           string   `json:"prompt"`
}

type LlamaApiResponse struct {
	Content            string `json:"content"`
	GenerationSettings struct {
		DynatempExponent       float64  `json:"dynatemp_exponent"`
		DynatempRange          float64  `json:"dynatemp_range"`
		FrequencyPenalty       float64  `json:"frequency_penalty"`
		Grammar                string   `json:"grammar"`
		IgnoreEos              bool     `json:"ignore_eos"`
		LogitBias              []any    `json:"logit_bias"`
		MinKeep                int      `json:"min_keep"`
		MinP                   float64  `json:"min_p"`
		Mirostat               int      `json:"mirostat"`
		MirostatEta            float64  `json:"mirostat_eta"`
		MirostatTau            float64  `json:"mirostat_tau"`
		Model                  string   `json:"model"`
		NCtx                   int      `json:"n_ctx"`
		NKeep                  int      `json:"n_keep"`
		NPredict               int      `json:"n_predict"`
		NProbs                 int      `json:"n_probs"`
		PenalizeNl             bool     `json:"penalize_nl"`
		PenaltyPromptTokens    []any    `json:"penalty_prompt_tokens"`
		PresencePenalty        float64  `json:"presence_penalty"`
		RepeatLastN            int      `json:"repeat_last_n"`
		RepeatPenalty          float64  `json:"repeat_penalty"`
		Samplers               []string `json:"samplers"`
		Seed                   int64    `json:"seed"`
		Stop                   []string `json:"stop"`
		Stream                 bool     `json:"stream"`
		Temperature            float64  `json:"temperature"`
		TfsZ                   float64  `json:"tfs_z"`
		TopK                   int      `json:"top_k"`
		TopP                   float64  `json:"top_p"`
		TypicalP               float64  `json:"typical_p"`
		UsePenaltyPromptTokens bool     `json:"use_penalty_prompt_tokens"`
	} `json:"generation_settings"`
	Model        string `json:"model"`
	Prompt       string `json:"prompt"`
	SlotID       int    `json:"slot_id"`
	Stop         bool   `json:"stop"`
	StoppedEos   bool   `json:"stopped_eos"`
	StoppedLimit bool   `json:"stopped_limit"`
	StoppedWord  bool   `json:"stopped_word"`
	StoppingWord string `json:"stopping_word"`
	Timings      struct {
		PredictedMs         float64 `json:"predicted_ms"`
		PredictedN          int     `json:"predicted_n"`
		PredictedPerSecond  float64 `json:"predicted_per_second"`
		PredictedPerTokenMs float64 `json:"predicted_per_token_ms"`
		PromptMs            float64 `json:"prompt_ms"`
		PromptN             int     `json:"prompt_n"`
		PromptPerSecond     float64 `json:"prompt_per_second"`
		PromptPerTokenMs    float64 `json:"prompt_per_token_ms"`
	} `json:"timings"`
	TokensCached    int  `json:"tokens_cached"`
	TokensEvaluated int  `json:"tokens_evaluated"`
	TokensPredicted int  `json:"tokens_predicted"`
	Truncated       bool `json:"truncated"`
}
