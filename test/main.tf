terraform {
  required_providers {
    trello = {
      source  = "terraform.local/edu/trello"
      version = "2.0"
    }
  }
}

resource "trello_board" "my-board-name" {
  key            = "e49bb022fb242733d8be9cffdcad6a96"
  token          = "858bef8d6ed9e8bbc65e215b324a77be37e8c781410b1d76a4cc2333a6417005"
  workspace_name = "a1"
  board_name     = "b3"
  cards          = ["Done(front)", "Done(Back)", "Bugs", "InProgress", "QA", "ToDo(front)", "ToDo(back)", "UX/UI", "Backlog", "Links"]
  member_emails = ["marcobert37@gmail.com"]
}
