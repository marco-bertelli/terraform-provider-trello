---
page_title: "Provider: Trello Provider"
subcategory: "cloud automation"
description: |-
  a provider that create a workspace and a custom board with list using trello API.
---

# Trello Provider

this provider interact with trello API to creeate board (is intended for create trello board automatic for every project):
1. you need a token and a trello key (see trello docs).
2. see the bove example.


## Example Usage

Do not keep your authentication key and token in HCL for production environments, use Terraform environment variables.

```terraform {
  required_providers {
    trello = {
      source = "marco-bertelli/trello"
      version = "0.2.2"
    }
  }
}


resource "trello_board" "my-board-name" {
	key = "your-key"
  token = "your-token"
  workspace_name = "terraform-trello"
  board_name = "test"
  cards = ["new","todo","custom"]
}
```

## Schema

### Optional
