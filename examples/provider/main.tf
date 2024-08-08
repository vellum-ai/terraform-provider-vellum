terraform {
  required_providers {
    vellum = {
      source = "hashicorp.com/ai/vellum"
    }
  }

}

provider "vellum" {}

data "vellum_document_index" "reference" {
  name = "reference"
}

resource "vellum_document_index" "managed" {
  label = "Managed Index"
  name  = "managed-index"
}

data "vellum_ml_model" "reference" {
  name = "gpt-4o"
}

resource "vellum_ml_model" "managed" {
  name = "my-test-model"
  family = "GPT3"
  hosted_by = "OPENAI"
  developed_by = "OPENAI"
  visibility = "PRIVATE"
}
