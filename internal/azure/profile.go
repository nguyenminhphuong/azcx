package azure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Subscription represents an Azure subscription with all fields preserved
type Subscription map[string]interface{}

// Profile represents the Azure CLI profile with all fields preserved
type Profile struct {
	raw           map[string]interface{}
	Subscriptions []Subscription
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

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse Azure profile: %w", err)
	}

	profile := &Profile{raw: raw}

	// Extract subscriptions
	if subs, ok := raw["subscriptions"].([]interface{}); ok {
		profile.Subscriptions = make([]Subscription, len(subs))
		for i, sub := range subs {
			if subMap, ok := sub.(map[string]interface{}); ok {
				profile.Subscriptions[i] = Subscription(subMap)
			}
		}
	}

	return profile, nil
}

// SaveProfile saves the Azure profile to disk
func SaveProfile(profile *Profile) error {
	path, err := GetProfilePath()
	if err != nil {
		return err
	}

	// Update subscriptions in raw data
	subs := make([]interface{}, len(profile.Subscriptions))
	for i, sub := range profile.Subscriptions {
		subs[i] = map[string]interface{}(sub)
	}
	profile.raw["subscriptions"] = subs

	data, err := json.MarshalIndent(profile.raw, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// GetCurrentSubscription returns the current default subscription
func (p *Profile) GetCurrentSubscription() *Subscription {
	for i := range p.Subscriptions {
		if p.Subscriptions[i].IsDefault() {
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
		if strings.ToLower(sub.GetName()) == nameOrIDLower || strings.ToLower(sub.GetID()) == nameOrIDLower {
			sub.SetDefault(true)
			found = true
		} else {
			sub.SetDefault(false)
		}
	}

	if !found {
		return fmt.Errorf("subscription '%s' not found", nameOrID)
	}

	return nil
}

// Helper methods for Subscription

func (s Subscription) GetID() string {
	if id, ok := s["id"].(string); ok {
		return id
	}
	return ""
}

func (s Subscription) GetName() string {
	if name, ok := s["name"].(string); ok {
		return name
	}
	return ""
}

func (s Subscription) IsDefault() bool {
	if isDefault, ok := s["isDefault"].(bool); ok {
		return isDefault
	}
	return false
}

func (s Subscription) SetDefault(value bool) {
	s["isDefault"] = value
}
