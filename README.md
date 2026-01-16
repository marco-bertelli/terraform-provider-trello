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
3. if you need help write to marcobert37@gmail.com

## News
1. collaborators email automatic invites for private boards, this feature help automatic invite default users for private or public boards
2. workspace members automatic invite, i can choose to add a people into the workspace but not in the board (for example external people), pass an empty array if not used
3. Custom cards for each board! Now you can specify different cards (lists) for each board using the `boards` array of objects. Each board has its own `name` and `cards` array.
4. **NEW**: Custom labels for each board! Now you can specify labels with name and color for each board. Valid colors: yellow, purple, blue, red, green, orange, black, sky, pink, lime.


## Example Usage

Do not keep your authentication key and token in HCL for production environments, use Terraform environment variables.
if member_emails or workspace_members are not used pass and empty array or terraform will throw an error

```terraform {
  required_providers {
    trello = {
      source = "marco-bertelli/trello"
      version = "3.4.0"
    }
  }
}


resource "trello_board" "my-board-name" {
	key = "your-key"
  token = "your-token"
  workspace_name = "terraform-trello"
  boards = [
    {
      name = "Development Board"
      cards = ["Backlog", "In Progress", "Done"]
      labels = [
        {
          name = "Bug"
          color = "red"
        },
        {
          name = "Feature"
          color = "green"
        },
        {
          name = "Enhancement"
          color = "blue"
        }
      ]
    },
    {
      name = "Marketing Board"
      cards = ["Ideas", "Planning", "Execution", "Complete"]
      labels = [
        {
          name = "Campaign"
          color = "purple"
        },
        {
          name = "Content"
          color = "orange"
        }
      ]
    },
    {
      name = "HR Board"
      cards = ["Recruiting", "Onboarding", "Training"]
      labels = [
        {
          name = "Urgent"
          color = "red"
        },
        {
          name = "Interview"
          color = "yellow"
        }
      ]
    }
  ]
  member_emails = ["marco.bertelli@testcollaborator.com"]
  workspace_members = [{email = "marco.bertelli@testcollaborator.com", role = "normal", name = "marco"}]
}
```

## Help requests
open a pr if you need help or find a bug
