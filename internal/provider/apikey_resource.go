// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-foxglove-cloud/internal/foxglove"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ApikeyResource{}
var _ resource.ResourceWithImportState = &ApikeyResource{}

func NewApikeyResource() resource.Resource {
	return &ApikeyResource{}
}

// ApikeyResource defines the resource implementation.
type ApikeyResource struct {
	foxgloveClient *foxglove.Client
}

// ApikeyResourceModel describes the resource data model.
type ApikeyResourceModel struct {
	Label        types.String `tfsdk:"label"`
	Capabilities types.List   `tfsdk:"capabilities"`
	Id           types.String `tfsdk:"id"`
	Secret       types.String `tfsdk:"secret"`
}

func (akr *ApikeyResourceModel) CapabilitiesValue() []string {
	capabilities := []string{}
	for _, capability := range akr.Capabilities.Elements() {
		capabilities = append(capabilities, capability.(types.String).ValueString())
	}
	return capabilities
}

func (r *ApikeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikey"
}

func (r *ApikeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Device",
		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				MarkdownDescription: "The human-readable label for this key.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{},
			},
			"capabilities": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "Capabilities of this key",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Opaque identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The secret token",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ApikeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	foxgloveClient, ok := req.ProviderData.(*foxglove.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *foxglove.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.foxgloveClient = foxgloveClient
}

func (r *ApikeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApikeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	newDevice, err := r.foxgloveClient.CreateAPIKey(foxglove.CreateAPIKeyRequest{
		Label:        data.Label.ValueString(),
		Capabilities: data.CapabilitiesValue(),
	})

	if err != nil {
		resp.Diagnostics.AddError("failed to create device "+strings.Join(data.CapabilitiesValue(), ";"), err.Error())
		return
	}
	trueCapabilities, _ := types.ListValueFrom(ctx, types.StringType, newDevice.Capabilities)

	resp.Diagnostics.Append(resp.State.Set(ctx, &ApikeyResourceModel{
		Id:           types.StringValue(newDevice.ID),
		Secret:       types.StringValue(newDevice.SecretToken),
		Label:        types.StringValue(newDevice.Label),
		Capabilities: trueCapabilities,
	})...)
}

func (r *ApikeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApikeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApikeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApikeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey, err := r.foxgloveClient.UpdateAPIKey(data.Id.ValueString(), foxglove.UpdateAPIKeyRequest{
		Label:        data.Label.ValueString(),
		Capabilities: data.CapabilitiesValue(),
	})

	if err != nil {
		resp.Diagnostics.AddError("failed to update apiKey", err.Error())
		return
	}

	trueCapabilities, _ := types.ListValueFrom(ctx, types.StringType, apiKey.Capabilities)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &ApikeyResourceModel{
		Label:        types.StringValue(apiKey.Label),
		Id:           types.StringValue(apiKey.ID),
		Secret:       data.Secret,
		Capabilities: trueCapabilities,
	})...)
}

func (r *ApikeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApikeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.foxgloveClient.DeleteAPIKey(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete apiKey", err.Error())
		return
	}
}

func (r *ApikeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
