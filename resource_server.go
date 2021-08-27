// resource_server.go
package main

import (
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        "log"
        "net/http"
        "encoding/json"
        "strconv"
)

func resourceServer() *schema.Resource {
        return &schema.Resource{
                Create: resourceServerCreate,
                Read:   resourceServerRead,
                Update: resourceServerUpdate,
                Delete: resourceServerDelete,

                Schema: map[string]*schema.Schema{
                        "key": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "token": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "workspace_name": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "board_name": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "board_id": &schema.Schema{
                                Type:     schema.TypeString,
                                Computed: true,
                        },
                        "workspace_id": &schema.Schema{
                                Type:     schema.TypeString,
                                Computed: true,
                        },
                        "cards": {
                                Type:     schema.TypeList,
                                Required: true,
                                Elem: &schema.Schema{
                                  Type: schema.TypeString,
                                },
                        },
                },
        }
}

type Body struct {
        Id string `json:"id"`
    }

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
        key := d.Get("key").(string)
        token := d.Get("token").(string)
        workspace_name := d.Get("workspace_name").(string)
        board_name := d.Get("board_name").(string)


        workspace, err := http.Post("https://api.trello.com/1/organizations?key="+key+"&token="+token+"&displayName="+workspace_name,"application/json",nil)

        if err != nil {
                log.Fatalln(err)
        }

        //lettura body.
        body := new(Body)
        
        json.NewDecoder(workspace.Body).Decode(body)

        d.Set("workspace_id",body.Id)

        board, boardError := http.Post("https://api.trello.com/1/boards?key="+key+"&token="+token+"&idOrganization="+body.Id+"&=&name="+board_name+"&defaultLists=false","application/json",nil)
        
        //lettura body.
        body1 := new(Body)
        
        json.NewDecoder(board.Body).Decode(body1)

        d.Set("board_id",body1.Id)

        if boardError != nil {
                log.Fatalln(board)
        }

        // cards for the current board read and create
        itemsRaw := d.Get("cards").([]interface{})
        items := make([]string, len(itemsRaw))
        
        for i, raw := range itemsRaw {
        items[i] = raw.(string)
        log.Println("[ERROR] "+ items[i])

        lists, listsError := http.Post("https://api.trello.com/1/lists?key="+key+"&token="+token+"&name="+items[i]+"&idBoard="+body1.Id,"application/json",nil)
        if listsError != nil {
                log.Fatalln(lists)
        }

        }
        //final close of main req
        defer board.Body.Close()
        defer workspace.Body.Close()
        
        d.SetId(board_name)
        
        return resourceServerRead(d, m)
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {

        return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {

        key := d.Get("key").(string)
        token := d.Get("token").(string)
       
        board_name := d.Get("board_name").(string)
        board_id := d.Get("board_id").(string)

        log.Println("[ERROR] "+board_id)

	request, err := http.NewRequest("PUT", "https://api.trello.com/1/boards/"+board_id+"?key="+key+"&token="+token+"&name="+board_name, nil)

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	workspaceonse, err := http.DefaultClient.Do(request)

        log.Println("[ERROR] "+strconv.Itoa(workspaceonse.StatusCode))
	if err != nil {
		log.Fatal(err)
	} else {
		defer workspaceonse.Body.Close()
		
	}


        return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
        board_id := d.Get("board_id").(string)
        workspace_id := d.Get("workspace_id").(string)

        key := d.Get("key").(string)
        token := d.Get("token").(string)

        // chiamata delete della board
        board, boardErr := http.NewRequest("DELETE", "https://api.trello.com/1/organizations/"+workspace_id+"?key="+key+"&token="+token, nil)

        board.Header.Set("Content-Type", "application/json; charset=utf-8")

	boardCall, err := http.DefaultClient.Do(board)

        if boardErr != nil {
		log.Fatal(boardErr)
	} else {
		defer boardCall.Body.Close()
		
	}

        // chiamata delete del workspace

	request, err := http.NewRequest("DELETE", "https://api.trello.com/1/boards/"+board_id+"?key="+key+"&token="+token, nil)

        request.Header.Set("Content-Type", "application/json; charset=utf-8")

	workspaceonse, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal(err)
	} else {
		defer workspaceonse.Body.Close()
		
	}

        return resourceServerRead(d, m)
}