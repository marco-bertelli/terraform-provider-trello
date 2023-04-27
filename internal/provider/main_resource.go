package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ! Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TerraformResource{}

func NewTerraformResource() resource.Resource {
	return &TerraformResource{}
}

type TerraformResource struct {
	client *http.Client
}

type TerraformResourceModel struct {
	Key            types.String `tfsdk:"key"`
	Token          types.String `tfsdk:"token"`
	Workspace_name types.String `tfsdk:"workspace_name"`
	Board_name     types.String `tfsdk:"board_name"`
	Board_id       types.String `tfsdk:"board_id"`
	Workspace_id   types.String `tfsdk:"workspace_id"`
	Cards          []string     `tfsdk:"cards"`
	Member_emails  []string     `tfsdk:"member_emails"`
}

func (r *TerraformResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "trello_board"
}

func (r *TerraformResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform trello board maker",

		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				MarkdownDescription: "Trello secret key.",
				Required:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Trello secret token.",
				Required:            true,
			},
			"workspace_name": schema.StringAttribute{
				MarkdownDescription: "name of the workspace.",
				Required:            true,
			},
			"board_name": schema.StringAttribute{
				MarkdownDescription: "Name of the board to be created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"board_id": schema.StringAttribute{
				MarkdownDescription: "id of created board.",
				Computed:            true,
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "id of created workspace.",
				Computed:            true,
			},
			"cards": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "cards of which board will be filled.",
			},
			"member_emails": schema.ListAttribute{
				MarkdownDescription: "email of members to send invite.",
				ElementType:         types.StringType,
				Required:            true,
			},
		},
	}
}

type TrelloApiResponse struct {
	Id string `json:"id"`
}

func (r *TerraformResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TerraformResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *TerraformResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	key := data.Key.ValueString()
	token := data.Token.ValueString()
	workspace_name := data.Workspace_name.ValueString()
	board_name := data.Board_name.ValueString()
	cards := data.Cards
	member_emails := data.Member_emails

	workspace, err := http.Post("https://api.trello.com/1/organizations?key="+key+"&token="+token+"&displayName="+workspace_name, "application/json", nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to create the resource. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)

		return
	}

	workspaceResponse := new(TrelloApiResponse)

	workspaceResponseErr := json.NewDecoder(workspace.Body).Decode(&workspaceResponse)

	if workspaceResponseErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to parse json the resource. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)

		return
	}

	data.Workspace_id = types.StringValue(workspaceResponse.Id)

	board, err := http.Post("https://api.trello.com/1/boards?key="+key+"&token="+token+"&idOrganization="+workspaceResponse.Id+"&=&name="+board_name+"&defaultLists=false", "application/json", nil)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to create the resource. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)

		return
	}

	boardResponse := new(TrelloApiResponse)

	boardResponseErr := json.NewDecoder(board.Body).Decode(boardResponse)

	if boardResponseErr != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to parse json the resource. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)

		return
	}

	data.Board_id = types.StringValue(boardResponse.Id)

	for i := range cards {

		_, listsError := http.Post("https://api.trello.com/1/lists?key="+key+"&token="+token+"&name="+cards[i]+"&idBoard="+boardResponse.Id, "application/json", nil)

		if listsError != nil {
			resp.Diagnostics.AddError(
				"Unable to Create Cards",
				"An unexpected error occurred while attempting to create the cards. "+
					"Please retry the operation or report this issue to the provider developers.\n\n"+
					"HTTP Error: "+err.Error(),
			)
		}
	}

	for i := range member_emails {
		emailReq, _ := http.NewRequest("PUT", "https://api.trello.com/1/boards/"+boardResponse.Id+"/members?key="+key+"&token="+token+"&email="+member_emails[i], nil)

		emailReq.Header.Set("Content-Type", "application/json; charset=utf-8")

		_, emailErr := http.DefaultClient.Do(emailReq)

		if emailErr != nil {
			resp.Diagnostics.AddError(
				"Unable to Invite Members",
				"An unexpected error occurred while attempting to invite. "+
					"Please retry the operation or report this issue to the provider developers.\n\n"+
					"HTTP Error: "+err.Error(),
			)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TerraformResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (r *TerraformResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *TerraformResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	key := data.Key.ValueString()
	token := data.Token.ValueString()
	board_name := data.Board_name.ValueString()
	board_id := data.Board_id.ValueString()

	request, err := http.NewRequest("PUT", "https://api.trello.com/1/boards/"+board_id+"?key="+key+"&token="+token+"&name="+board_name, nil)

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	_, err = http.DefaultClient.Do(request)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Board",
			"An unexpected error occurred while attempting to update the board. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TerraformResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TerraformResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	key := data.Key.ValueString()
	token := data.Token.ValueString()
	board_id := data.Board_id.ValueString()
	workspace_id := data.Workspace_id.ValueString()

	workspaceRequest, _ := http.NewRequest("DELETE", "https://api.trello.com/1/organizations/"+workspace_id+"?key="+key+"&token="+token, nil)

	workspaceRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

	_, err := http.DefaultClient.Do(workspaceRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable Delete Workspace",
			"An unexpected error occurred while attempting to delete workspace. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)

		return
	}

	boardRequest, err := http.NewRequest("DELETE", "https://api.trello.com/1/boards/"+board_id+"?key="+key+"&token="+token, nil)

	boardRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

	_, err = http.DefaultClient.Do(boardRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable Delete Board",
			"An unexpected error occurred while attempting to delete board. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				"HTTP Error: "+err.Error(),
		)

		return
	}

	return
}
