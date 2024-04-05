// This file was auto-generated by Fern from our API Definition.

package vellum

import (
	json "encoding/json"
	core "github.com/vellum-ai/terraform-provider-vellum/internal/vellum/core"
)

type BadRequestError struct {
	*core.APIError
	Body interface{}
}

func (b *BadRequestError) UnmarshalJSON(data []byte) error {
	var body interface{}
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	b.StatusCode = 400
	b.Body = body
	return nil
}

func (b *BadRequestError) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Body)
}

func (b *BadRequestError) Unwrap() error {
	return b.APIError
}

type ConflictError struct {
	*core.APIError
	Body *RegisterPromptErrorResponse
}

func (c *ConflictError) UnmarshalJSON(data []byte) error {
	var body *RegisterPromptErrorResponse
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	c.StatusCode = 409
	c.Body = body
	return nil
}

func (c *ConflictError) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Body)
}

func (c *ConflictError) Unwrap() error {
	return c.APIError
}

type ForbiddenError struct {
	*core.APIError
	Body interface{}
}

func (f *ForbiddenError) UnmarshalJSON(data []byte) error {
	var body interface{}
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	f.StatusCode = 403
	f.Body = body
	return nil
}

func (f *ForbiddenError) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Body)
}

func (f *ForbiddenError) Unwrap() error {
	return f.APIError
}

type InternalServerError struct {
	*core.APIError
	Body interface{}
}

func (i *InternalServerError) UnmarshalJSON(data []byte) error {
	var body interface{}
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	i.StatusCode = 500
	i.Body = body
	return nil
}

func (i *InternalServerError) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Body)
}

func (i *InternalServerError) Unwrap() error {
	return i.APIError
}

type NotFoundError struct {
	*core.APIError
	Body interface{}
}

func (n *NotFoundError) UnmarshalJSON(data []byte) error {
	var body interface{}
	if err := json.Unmarshal(data, &body); err != nil {
		return err
	}
	n.StatusCode = 404
	n.Body = body
	return nil
}

func (n *NotFoundError) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Body)
}

func (n *NotFoundError) Unwrap() error {
	return n.APIError
}
