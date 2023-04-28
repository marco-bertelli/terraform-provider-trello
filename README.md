---
page_title: "Provider: Trello Provider"
subcategory: "cloud automation"
description: |-
  a provider that create a workspace and a custom board with list using trello API.
---

# Trello Provider

this provider interact with trello API to create board (is intended for create trello board automatic for every project):
1. you need a token and a trello key (see trello docs).
2. see the bove example.

## News
1. collaborators email automatic invites for private boards, this feature help automatic invite default users for private or public boards
2. workspace members automatic invite, i can choose to add a people into the workspace but not in the board (for example external people), pass an empty array if not used


## Example Usage

Do not keep your authentication key and token in HCL for production environments, use Terraform environment variables.
if member_emails or workspace_members are not used pass and empty array or terraform will throw an error

```terraform {
  required_providers {
    trello = {
      source = "marco-bertelli/trello"
      version = "3.1.0"
    }
  }
}


resource "trello_board" "my-board-name" {
	key = "your-key"
  token = "your-token"
  workspace_name = "terraform-trello"
  board_name = "test"
  cards = ["new","todo","custom"]
  member_emails = ["marco.bertelli@testcollaborator.com"],
  workspace_members = [{"email": "marco.bertelli@testcollaborator.com", "role": "normal", "name": "marco"}]
}
```

## Help requests
open a pr if you need help or find a bug

### Optional
