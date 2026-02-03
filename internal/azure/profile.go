package azure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Subscription represents an Azure subscription
type Subscription struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	State            string `json:"state"`
	TenantID         string `json:"tenantId"`
	IsDefault        bool   `json:"isDefault"`
	HomeTenantID     string `json:"homeTenantId,omitempty"`
	ManagedByTenants []any  `json:"managedByTenants,omitempty"`
}

// Profile represents the Azure CLI profile
type Profile struct {
	InstallationID string         `json:"installationId"`
	Subscriptions  []Subscription `json:"subscriptions"`
}

// GetProfilePath returns the path to the Azure profile file
// Supports AZURE_CONFIG_DIR environment variable (same as Azure CLI)
func GetProfilePath() (string, error) {
	// Check for custom Azure config directory
	if configDir := os.Getenv("AZURE_CONFIG_DIR"); configDir != "" {
		return filepath.Join(configDir, "azureProfile.json"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".azure", "azureProfile.json"), nil
}

// LoadProfile loads the Azure profile from disk
func LoadProfile() (*Profile, error) {
	path, err := GetProfilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read Azure profile: %w (have you run 'az login'?)", err)
	}

	// Handle BOM (Azure CLI sometimes adds it)
	data = []byte(strings.TrimPrefix(string(data), "\ufeff"))

	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("failed to parse Azure profile: %w", err)
	}

	return &profile, nil
}

// SaveProfile saves the Azure profile to disk
func SaveProfile(profile *Profile) error {
	path, err := GetProfilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// GetCurrentSubscription returns the current default subscription
func (p *Profile) GetCurrentSubscription() *Subscription {
	for i := range p.Subscriptions {
		if p.Subscriptions[i].IsDefault {
			return &p.Subscriptions[i]
		}
	}
	return nil
}

// SetSubscription sets the specified subscription as default
func (p *Profile) SetSubscription(nameOrID string) error {
	found := false
	nameOrIDLower := strings.ToLower(nameOrID)

	for i := range p.Subscriptions {
		sub := &p.Subscriptions[i]
		if strings.ToLower(sub.Name) == nameOrIDLower || strings.ToLower(sub.ID) == nameOrIDLower {
			sub.IsDefault = true
			found = true
		} else {
			sub.IsDefault = false
		}
	}

	if !found {
		return fmt.Errorf("subscription '%s' not found", nameOrID)
	}

	return nil
}

// FindSubscription finds a subscription by name or ID
func (p *Profile) FindSubscription(nameOrID string) *Subscription {
	nameOrIDLower := strings.ToLower(nameOrID)
	for i := range p.Subscriptions {
		sub := &p.Subscriptions[i]
		if strings.ToLower(sub.Name) == nameOrIDLower || strings.ToLower(sub.ID) == nameOrIDLower {
			return sub
		}
	}
	return nil
}
