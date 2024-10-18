package foxglove

import (
	"testing"
	"time"
)

func TestDeviceLifecycle(t *testing.T) {
	client := NewClient("") // enter api key here

	// Step 1: Create a new device
	deviceName := "terraform_unit_test_" + time.Now().Format("20060102150405")
	createReq := CreateDeviceRequest{
		Name: deviceName,
	}

	createResp, err := client.CreateDevice(createReq)
	if err != nil {
		t.Fatalf("Failed to create device: %v", err)
		return
	}

	if createResp.Name != deviceName {
		t.Fatalf("Expected device name %s, but got %s", deviceName, createResp.Name)
	}

	createdDeviceID := createResp.ID
	t.Logf("Created device with ID: %s", createdDeviceID)

	// Step 2: List devices and verify the created device exists
	devices, err := client.ListDevices("", "", "", 100, 0)
	if err != nil {
		t.Fatalf("Failed to list devices: %v", err)
		return
	}

	var found bool
	for _, device := range devices {
		if device.ID == createdDeviceID {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Created device with ID %s not found in the list", createdDeviceID)
	}

	t.Log("Verified that the created device exists in the list")

	// Step 3: Retrieve the device by ID
	getResp, err := client.GetDevice(createdDeviceID)
	if err != nil {
		t.Fatalf("Failed to retrieve device: %v", err)
		return
	}

	if getResp.ID != createdDeviceID || getResp.Name != deviceName {
		t.Fatalf("Retrieved device details do not match. Expected ID %s and Name %s, but got ID %s and Name %s",
			createdDeviceID, deviceName, getResp.ID, getResp.Name)
	}

	t.Log("Successfully retrieved the device by ID")

	// Step 4: Change the name of the device
	deviceName = deviceName + "_updated"
	updateResp, err := client.UpdateDevice(createdDeviceID, UpdateDeviceRequest{
		Name: deviceName,
	})
	if err != nil {
		t.Fatalf("Failed to update device: %v", err)
		return
	}
	if updateResp.ID != createdDeviceID || updateResp.Name != deviceName {
		t.Fatalf("Received device details do not match after device update. Expected ID %s and Name %s, but got ID %s and Name %s",
			createdDeviceID, deviceName, getResp.ID, getResp.Name)
	}

	t.Log("Successfully updated the device name")

	// Step 5: Delete the device
	deleteResp, err := client.DeleteDevice(createdDeviceID)
	if err != nil {
		t.Fatalf("Failed to delete device: %v", err)
		return
	}

	if deleteResp.ID != createdDeviceID {
		t.Fatalf("Expected deleted device ID %s, but got %s", createdDeviceID, deleteResp.ID)
	}

	t.Log("Successfully deleted the device")

	// Optional Step 6: Verify the device is no longer listed
	time.Sleep(2 * time.Second) // Allow some time for the deletion to propagate

	devices, err = client.ListDevices("", "", "", 100, 0)
	if err != nil {
		t.Fatalf("Failed to list devices after deletion: %v", err)
	}

	for _, device := range devices {
		if device.ID == createdDeviceID {
			t.Fatalf("Deleted device with ID %s still found in the list", createdDeviceID)
		}
	}

	t.Log("Verified that the deleted device no longer exists in the list")
}
