package foxglove

import (
	"testing"
	"time"
)

func TestAPIKeyLifecycle(t *testing.T) {
	client := NewClient("") // enter api key here

	// Step 1: Create a new API key
	apiKeyName := "terraform_unit_test_" + time.Now().Format("20060102150405")
	createReq := CreateAPIKeyRequest{
		Label: apiKeyName,
		Capabilities: []string{
			"devices.create",
			"devices.delete",
		},
	}

	createResp, err := client.CreateAPIKey(createReq)
	if err != nil {
		t.Fatalf("Failed to create API key: %v", err)
		return
	}

	if createResp.Label != apiKeyName {
		t.Fatalf("Expected API key name %s, but got %s", apiKeyName, createResp.Label)
	}

	createdAPIKeyID := createResp.ID
	t.Logf("Created API key with ID: %s", createdAPIKeyID)

	// Step 2: List API keys and verify the created API key exists
	apiKeys, err := client.ListAPIKeys()
	if err != nil {
		t.Fatalf("Failed to list API keys: %v", err)
		return
	}

	var found bool
	for _, apiKey := range apiKeys {
		if apiKey.ID == createdAPIKeyID {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Created API key with ID %s not found in the list", createdAPIKeyID)
		return
	}

	t.Log("Verified that the created API key exists in the list")

	// Step 3: Change the name of the API key
	apiKeyName = apiKeyName + "_updated"
	updateResp, err := client.UpdateAPIKey(createdAPIKeyID, UpdateAPIKeyRequest{
		Label:        apiKeyName,
		Capabilities: []string{"recordings.list", "data.topics.list"},
	})
	if err != nil {
		t.Fatalf("Failed to update API key: %v", err)
		return
	}
	if updateResp.ID != createdAPIKeyID || updateResp.Label != apiKeyName {
		t.Fatalf("Received API key details do not match after API key update. Expected ID %s and Name %s, but got ID %s and Name %s",
			createdAPIKeyID, apiKeyName, updateResp.ID, updateResp.Label)
	}

	t.Log("Successfully updated the API key name")

	// Step 4: Delete the API key
	err = client.DeleteAPIKey(createdAPIKeyID)
	if err != nil {
		t.Fatalf("Failed to delete API key: %v", err)
		return
	}

	t.Log("Successfully deleted the API key")

	// Optional Step 6: Verify the API key is no longer listed
	time.Sleep(2 * time.Second) // Allow some time for the deletion to propagate

	apiKeys, err = client.ListAPIKeys()
	if err != nil {
		t.Fatalf("Failed to list API keys after deletion: %v", err)
		return
	}

	for _, apiKey := range apiKeys {
		if apiKey.ID == createdAPIKeyID {
			t.Fatalf("Deleted API key with ID %s still found in the list", createdAPIKeyID)
		}
	}

	t.Log("Verified that the deleted API key no longer exists in the list")
}
