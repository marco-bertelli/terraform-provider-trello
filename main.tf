terraform {
  required_providers {
    trello = {
      source = "hashicorp.com/edu/customprovider"
      version = "0.2"
    }
  }
}


resource "trello_board" "my-server-name" {
	key = "e49bb022fb242733d8be9cffdcad6a96"
  token = "858bef8d6ed9e8bbc65e215b324a77be37e8c781410b1d76a4cc2333a6417005"
  workspace_name = "custom"
  board_name = "prova"
  cards = ["new","todo","kek"]
}
