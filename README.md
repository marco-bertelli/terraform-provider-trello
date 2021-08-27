# terraform-provider-trello

this is a terraform provider that integrate trello. 
it create automaticaly trello organizations and boards at every project

# maintainment
this is currently mantained by Runelab S.r.l
# example

```
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
