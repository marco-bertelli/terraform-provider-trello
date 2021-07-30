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


        resp, err := http.Post("https://api.trello.com/1/organizations?key="+key+"&token="+token+"&displayName="+workspace_name,"application/json",nil)

        if err != nil {
                log.Fatalln(err)
        }

        //lettura body.
        body := new(Body)
        
        json.NewDecoder(resp.Body).Decode(body)


        resp1, err1 := http.Post("https://api.trello.com/1/boards?key="+key+"&token="+token+"&idOrganization="+body.Id+"&=&name="+board_name,"application/json",nil)
        
        //lettura body.
        body1 := new(Body)
        
        json.NewDecoder(resp1.Body).Decode(body1)

        d.Set("board_id",body1.Id)

        if err1 != nil {
                log.Fatalln(resp1)
        }

        defer resp1.Body.Close()
        defer resp.Body.Close()
        
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

	response, err := http.DefaultClient.Do(request)

        log.Println("[ERROR] "+strconv.Itoa(response.StatusCode))
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		
	}


        return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}