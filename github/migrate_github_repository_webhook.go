package github

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

func resourceGithubRepositoryWebhookMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found GitHub Repository Webhook State v0; migrating to v1")
		return migrateGithubRepositoryWebhookStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateGithubRepositoryWebhookStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] GitHub Repository Webhook Attributes before migration: %#v", is.Attributes)

	prefix := "configuration."

	delete(is.Attributes, prefix+"%")

	// Read & delete old keys
	oldKeys := make(map[string]string, 0)
	for k, v := range is.Attributes {
		if strings.HasPrefix(k, prefix) {
			oldKeys[k] = v

			// Delete old keys
			delete(is.Attributes, k)
		}
	}

	// Write new keys
	for k, v := range oldKeys {
		newKey := "configuration.0." + strings.TrimPrefix(k, prefix)
		is.Attributes[newKey] = v
	}

	is.Attributes[prefix+"#"] = "1"

	log.Printf("[DEBUG] GitHub Repository Webhook Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
