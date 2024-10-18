// SPDX-License-Identifier: MIT

package provider

import (
	"context"
	"fmt"
	"terraform-provider-foxglove-cloud/internal/foxglove"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &DeviceResource{}
var _ resource.ResourceWithImportState = &DeviceResource{}

func NewDeviceResource() resource.Resource {
	return &DeviceResource{}
}

// DeviceResource defines the resource implementation.
type DeviceResource struct {
	foxgloveClient *foxglove.Client
}

// DeviceResourceModel describes the resource data model.
type DeviceResourceModel struct {
	Name types.String `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

func (r *DeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *DeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Device",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the device.",
				Required:            true,
				PlanModifiers:       []planmodifier.String{},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Opaque identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	existingDevice, err := r.foxgloveClient.GetDevice(data.Name.ValueString())
	if err == nil {
		resp.State.Set(ctx, &DeviceResourceModel{
			Id:   types.StringValue(existingDevice.ID),
			Name: types.StringValue(existingDevice.Name),
		})
		return
	}

	device, err := r.foxgloveClient.CreateDevice(foxglove.CreateDeviceRequest{
		Name: data.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to create device", err.Error())
		return
	}

	data.Id = types.StringValue(device.ID)
	data.Name = types.StringValue(device.Name)

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeviceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var device *foxglove.GetDeviceResponse
	var err error
	if data.Id.IsUnknown() {
		device, err = r.foxgloveClient.GetDevice(data.Name.ValueString())
	} else {
		device, err = r.foxgloveClient.GetDevice(data.Id.ValueString())
	}

	if err != nil {
		// device not found, this is expected
		resp.State.RemoveResource(ctx)
		return
	}

	// The device exists, update the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &DeviceResourceModel{
		Id:   types.StringValue(device.ID),
		Name: types.StringValue(device.Name),
	})...)
}

func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeviceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	device, err := r.foxgloveClient.UpdateDevice(data.Id.ValueString(), foxglove.UpdateDeviceRequest{
		Name: data.Name.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("failed to update device", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &DeviceResourceModel{
		Name: types.StringValue(device.Name),
		Id:   types.StringValue(device.ID),
	})...)
}

func (r *DeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.foxgloveClient.DeleteDevice(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete device", err.Error())
		return
	}
}

func (r *DeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
