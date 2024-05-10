package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	Key               types.String            `tfsdk:"key"`
	Token             types.String            `tfsdk:"token"`
	Workspace_name    types.String            `tfsdk:"workspace_name"`
	Board_names       []string                `tfsdk:"board_names"`
	Board_id_1        types.String            `tfsdk:"board_id_1"`
	Board_id_2        types.String            `tfsdk:"board_id_2"`
	Board_id_3        types.String            `tfsdk:"board_id_3"`
	Workspace_id      types.String            `tfsdk:"workspace_id"`
	Cards             []string                `tfsdk:"cards"`
	Member_emails     []string                `tfsdk:"member_emails"`
	Workspace_members []*WorkspaceMemberModel `tfsdk:"workspace_members"`
}

type WorkspaceMemberModel struct {
	Email types.String `tfsdk:"email"`
	Name  types.String `tfsdk:"name"`
	Role  types.String `tfsdk:"role"`
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
			"board_names": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Names of the boards to be created.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "id of created workspace.",
				Computed:            true,
			},
			"board_id_1": schema.StringAttribute{
				MarkdownDescription: "id of created workspace.",
				Computed:            true,
			},
			"board_id_2": schema.StringAttribute{
				MarkdownDescription: "id of created workspace.",
				Computed:            true,
			},
			"board_id_3": schema.StringAttribute{
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
			"workspace_members": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"email": schema.StringAttribute{
						Required: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"role": schema.StringAttribute{
						Required: true,
					},
				}},
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
	board_names := data.Board_names
	cards := data.Cards
	member_emails := data.Member_emails
	workspace_members := data.Workspace_members

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
				"HTTP Error: "+workspaceResponseErr.Error(),
		)

		return
	}

	data.Workspace_id = types.StringValue(workspaceResponse.Id)

	for i := range board_names {
		board, err := http.Post("https://api.trello.com/1/boards?key="+key+"&token="+token+"&idOrganization="+workspaceResponse.Id+"&=&name="+board_names[i]+"&defaultLists=false", "application/json", nil)

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
					"HTTP Error: "+boardResponseErr.Error(),
			)

			return
		}

		if i == 0 {
			data.Board_id_1 = types.StringValue(boardResponse.Id)
		} else if i == 1 {
			data.Board_id_2 = types.StringValue(boardResponse.Id)
		} else if i == 2 {
			data.Board_id_3 = types.StringValue(boardResponse.Id)
		}

		for i := range cards {
			tflog.Debug(ctx, "https://api.trello.com/1/lists?key="+key+"&token="+token+"&name="+cards[i]+"&idBoard="+boardResponse.Id)
			_, listsError := http.Post("https://api.trello.com/1/lists?key="+key+"&token="+token+"&name="+cards[i]+"&idBoard="+boardResponse.Id, "application/json", nil)

			if listsError != nil {
				resp.Diagnostics.AddError(
					"Unable to Create Cards",
					"An unexpected error occurred while attempting to create the cards. "+
						"Please retry the operation or report this issue to the provider developers.\n\n"+
						"HTTP Error: "+listsError.Error(),
				)
			}
		}

		for i := range member_emails {
			tflog.Debug(ctx, "https://api.trello.com/1/boards/"+boardResponse.Id+"/members?key="+key+"&token="+token+"&email="+member_emails[i])
			emailReq, _ := http.NewRequest("PUT", "https://api.trello.com/1/boards/"+boardResponse.Id+"/members?key="+key+"&token="+token+"&email="+member_emails[i], nil)

			emailReq.Header.Set("Content-Type", "application/json; charset=utf-8")

			_, emailErr := http.DefaultClient.Do(emailReq)

			if emailErr != nil {
				resp.Diagnostics.AddError(
					"Unable to Invite Members",
					"An unexpected error occurred while attempting to invite. "+
						"Please retry the operation or report this issue to the provider developers.\n\n"+
						"HTTP Error: "+emailErr.Error(),
				)
			}
		}
	}

	for i := range workspace_members {
		member_email := workspace_members[i].Email.ValueString()
		member_name := workspace_members[i].Name.ValueString()
		member_role := workspace_members[i].Role.ValueString()

		tflog.Debug(ctx, "https://api.trello.com/1/organizations/"+workspaceResponse.Id+"/members?email="+member_email+"&fullName="+member_name+"&type="+member_role+"&key="+key+"&token="+token)
		workspaceMemberRequest, _ := http.NewRequest("PUT", "https://api.trello.com/1/organizations/"+workspaceResponse.Id+"/members?email="+member_email+"&fullName="+member_name+"&type="+member_role+"&key="+key+"&token="+token, nil)

		workspaceMemberRequest.Header.Set("Content-Type", "application/json; charset=utf-8")

		_, workspaceMemberErr := http.DefaultClient.Do(workspaceMemberRequest)

		if workspaceMemberErr != nil {
			resp.Diagnostics.AddError(
				"Unable to Invite Workspace Members",
				"An unexpected error occurred while attempting to invite. "+
					"Please retry the operation or report this issue to the provider developers.\n\n"+
					"HTTP Error: "+workspaceMemberErr.Error(),
			)
		}
	}

	// ! we have to fill nulls empty board ids
	if data.Board_id_1.ValueString() == "" {
		data.Board_id_1 = types.StringValue("null")
	}

	if data.Board_id_2.ValueString() == "" {
		data.Board_id_2 = types.StringValue("null")
	}

	if data.Board_id_3.ValueString() == "" {
		data.Board_id_3 = types.StringValue("null")
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
	board_names := data.Board_names
	board_ids := []types.String{data.Board_id_1, data.Board_id_2, data.Board_id_3}

	for i := range board_ids {
		// ! we have to skip null board ids
		if board_ids[i].ValueString() == "null" {
			continue
		}

		board_name := board_names[i]

		request, err := http.NewRequest("PUT", "https://api.trello.com/1/boards/"+board_ids[i].ValueString()+"?key="+key+"&token="+token+"&name="+board_name, nil)

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
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TerraformResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *TerraformResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	key := data.Key.ValueString()
	token := data.Token.ValueString()
	workspace_id := data.Workspace_id.ValueString()
	board_ids := []types.String{data.Board_id_1, data.Board_id_2, data.Board_id_3}

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

	for i := range board_ids {
		if board_ids[i].ValueString() == "null" {
			continue
		}

		boardRequest, err := http.NewRequest("DELETE", "https://api.trello.com/1/boards/"+board_ids[i].ValueString()+"?key="+key+"&token="+token, nil)

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
	}

	return
}
