package trbac

import "github.com/BurntSushi/toml"

// Privileges defines the relationship between roles and permissions
type Privileges interface {
    // GetPermissions returns all permissions assigned to any of the given roles
    GetPermissions(roles []string) []Permission
}

// Permission represents allowed actions on resources
type Permission struct {
    Actions       []string // Allowed actions
    ResourceTypes []string // Applicable resource types
    Constraints   []string // Required constraints
}

// TOMLPrivileges loads privileges from TOML files
type TOMLPrivileges struct {
    permissions map[string][]Permission
}

// GetPermissions returns all permissions for the given roles
func (privs TOMLPrivileges) GetPermissions(roles []string) []Permission {
    result := []Permission{}
    for _, role := range roles {
        result = append(result, privs.permissions[role]...)
    }
    return result
}

// NewTOMLPrivileges creates a new TOML-based privileges provider
func NewTOMLPrivileges(tomlPath string) (Privileges, error) {
    var privs map[string][]Permission
    if _, err := toml.DecodeFile(tomlPath, &privs); err != nil {
        return nil, err
    }
    return TOMLPrivileges{privs}, nil
}
