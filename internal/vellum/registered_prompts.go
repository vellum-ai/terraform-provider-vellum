// This file was auto-generated by Fern from our API Definition.

package vellum

type RegisterPromptRequestRequest struct {
	// A human-friendly label for corresponding entities created in Vellum.
	Label string `json:"label" url:"label"`
	// A uniquely-identifying name for corresponding entities created in Vellum.
	Name string `json:"name" url:"name"`
	// Information about how to execute the prompt template.
	Prompt *RegisterPromptPromptInfoRequest `json:"prompt,omitempty" url:"prompt,omitempty"`
	// The initial LLM provider to use for this prompt
	//
	// * `ANTHROPIC` - Anthropic
	// * `AWS_BEDROCK` - AWS Bedrock
	// * `AZURE_OPENAI` - Azure OpenAI
	// * `COHERE` - Cohere
	// * `GOOGLE` - Google
	// * `HOSTED` - Hosted
	// * `MOSAICML` - MosaicML
	// * `OPENAI` - OpenAI
	// * `FIREWORKS_AI` - Fireworks AI
	// * `HUGGINGFACE` - HuggingFace
	// * `MYSTIC` - Mystic
	// * `PYQ` - Pyq
	// * `REPLICATE` - Replicate
	Provider *ProviderEnum `json:"provider,omitempty" url:"provider,omitempty"`
	// The initial model to use for this prompt
	Model string `json:"model" url:"model"`
	// The initial model parameters to use for  this prompt
	Parameters *RegisterPromptModelParametersRequest `json:"parameters,omitempty" url:"parameters,omitempty"`
	// Optionally include additional metadata to store along with the prompt.
	Meta map[string]interface{} `json:"meta,omitempty" url:"meta,omitempty"`
}