terraform {
  required_providers {
    vellum = {
      source = "hashicorp.com/ai/vellum"
    }
  }

}

provider "vellum" {
#   base_url = "http://localhost:8000"
}

# data "vellum_document_index" "reference" {
#   name = "reference"
# }

# resource "vellum_document_index" "managed" {
#   label = "Managed Index"
#   name  = "managed-index"
# }

# data "vellum_ml_model" "reference" {
#   name = "gpt-4o"
# }

resource "vellum_ml_model" "noas-test" {
  name = "noas-test"
  family = "CHAT_GPT"
  hosted_by = "OPENAI"
  developed_by = "OPENAI"
  visibility = "PRIVATE"
  exec_config = {
      base_url = "https://api.openai.com/v1"
      features = [
          "CHAT_MESSAGE_SYSTEM",
          "CHAT_MESSAGE_USER",
          "CHAT_MESSAGE_ASSISTANT",
          "CHAT_MESSAGE_FUNCTION_CALL",
          "CHAT_MESSAGE_IMAGE",
          "FUNCTION_DEFINITION",
          "STREAMING_SUPPORT"
      ]
      metadata = {key = "value"}
      model_identifier = "gpt-4o-mini"
  }
  parameter_config = {
#       stop = {
#         items = {
#           string = {}
#         }
#       }
      top_p = {
        format = "float"
      }
      logit_bias = {
#           pattern_properties = {
#               "^\\d+$" = {
#                 number = {
#                   format = "float"
#                   maximum = 100.0
#                   minimum = -100.0
#                 }
#               }
#           }
      }
      max_tokens = {
        maximum = 16383
        minimum = 1
      }
      temperature = {
        format = "float"
        maximum = 2.0
        minimum = 0.0
      }
      presence_penalty = {
        format = "float"
        maximum = 2.0
        minimum = -2.0
      }
      custom_parameters = {
          seed = {
              string = {
                title = "Seed"
                pattern = "^[1-9][0-9]*$"
                description = "If specified, OpenAI will make a best effort to sample deterministically, such that repeated requests with the same seed and parameters should return the same result. Determinism is not guaranteed, and you should refer to the system_fingerprint response parameter to monitor changes in the backend."
              }
          }
#           user = {
#               string = {
#                 title = "User ID",
#                 description = "A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.",
#               }
#           }
#           json_mode = {
#               boolean = {
#                 title = "JSON Mode",
#                 description = "Whether to return the response as a JSON object or a plain string.",
#               }
#           }
      }
      frequency_penalty = {
        format = "float"
        maximum = 2.0
        minimum = -2.0
      }
  }
  display_config={
      tags = ["CHAT"]
      label = "GPT-4o Mini"
      description = "OpenAI's affordable and intelligent small model for fast, lightweight tasks. GPT-4o mini is cheaper and more capable than GPT-3.5 Turbo."
      default_display_priority = 0
  }
}
