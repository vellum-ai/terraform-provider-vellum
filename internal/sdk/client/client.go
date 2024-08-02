package client

import (
	http "net/http"
	"terraform-provider-vellum/internal/sdk/mlmodels"

	core "terraform-provider-vellum/internal/sdk/core"
	documentindexes "terraform-provider-vellum/internal/sdk/documentindexes"
)

type Client struct {
	baseURL string
	caller  *core.Caller
	header  http.Header

	DocumentIndexes *documentindexes.Client
	MLModels        *mlmodels.Client
}

func NewClient(opts ...core.ClientOption) *Client {
	options := core.NewClientOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &Client{
		baseURL:         options.BaseURL,
		caller:          core.NewCaller(options.HTTPClient),
		header:          options.ToHeader(),
		DocumentIndexes: documentindexes.NewClient(opts...),
		MLModels:        mlmodels.NewClient(opts...),
	}
}
