// resource_server.go
package main

import (
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
        "log"
        "net/http"
        "encoding/json"
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
                                Required: false,
                        },
                },
        }
}

type Body struct {
        id string
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

        defer resp.Body.Close()

        //lettura body.
        body := new(Body)
        json.NewDecoder(resp.Body).Decode(body)

        resp1, err1 := http.Post("https://api.trello.com/1/boards?key="+key+"&token="+token+"&idOrganization="+body.id+"&=&name="+board_name,"application/json",nil)

        if err1 != nil {
                log.Fatalln(resp1)
        }

        defer resp1.Body.Close()
        
        d.SetId(board_name)
        d.Set("board_id",body.id)


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

        client := &http.Client{}
	request, err := http.NewRequest("PUT", "https://api.trello.com/1/boards?key="+key+"&token="+token+"&idOrganization="+board_id+"&=&name="+board_name, nil)
	request.SetBasicAuth("admin", "admin")
	request.ContentLength = 23
	response, err := client.Do(request)
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