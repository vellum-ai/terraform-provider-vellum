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
