package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestRoleConfigBackwardsCompatibility(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected RoleConfig
	}{
		{
			name:  "old format - simple array",
			input: `["you are a shell expert", "output only commands"]`,
			expected: RoleConfig{
				Prompt:         []string{"you are a shell expert", "output only commands"},
				AllowedTools:   nil,
				BlockedTools:   nil,
				AllowedServers: nil,
				BlockedServers: nil,
			},
		},
		{
			name: "new format - full config",
			input: `
prompt:
  - "you are a shell expert"
  - "output only commands"
allowed_servers: ["filesystem", "shell"]
blocked_servers: ["network"]
allowed_tools: ["filesystem_*", "shell_*"]
blocked_tools: ["dangerous_*"]`,
			expected: RoleConfig{
				Prompt:         []string{"you are a shell expert", "output only commands"},
				AllowedTools:   []string{"filesystem_*", "shell_*"},
				BlockedTools:   []string{"dangerous_*"},
				AllowedServers: []string{"filesystem", "shell"},
				BlockedServers: []string{"network"},
			},
		},
		{
			name: "new format - only prompt",
			input: `
prompt:
  - "you are helpful"`,
			expected: RoleConfig{
				Prompt:         []string{"you are helpful"},
				AllowedTools:   nil,
				BlockedTools:   nil,
				AllowedServers: nil,
				BlockedServers: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config RoleConfig
			err := yaml.Unmarshal([]byte(tt.input), &config)
			require.NoError(t, err)
			require.Equal(t, tt.expected, config)
		})
	}
}

func TestRoleFiltering(t *testing.T) {
	tests := []struct {
		name       string
		toolName   string
		serverName string
		roleConfig RoleConfig
		expected   bool
	}{
		{
			name:       "no restrictions - allow all",
			toolName:   "read_file",
			serverName: "filesystem",
			roleConfig: RoleConfig{},
			expected:   true,
		},
		{
			name:       "blocked tool - deny",
			toolName:   "dangerous_delete",
			serverName: "filesystem",
			roleConfig: RoleConfig{
				BlockedTools: []string{"*_delete"},
			},
			expected: false,
		},
		{
			name:       "allowed tool pattern - allow",
			toolName:   "read_file",
			serverName: "filesystem",
			roleConfig: RoleConfig{
				AllowedTools: []string{"filesystem_read_*"},
			},
			expected: true,
		},
		{
			name:       "allowed tool pattern - deny non-matching",
			toolName:   "write_file",
			serverName: "filesystem",
			roleConfig: RoleConfig{
				AllowedTools: []string{"filesystem_read_*"},
			},
			expected: false,
		},
		{
			name:       "blocked takes precedence over allowed",
			toolName:   "read_file",
			serverName: "filesystem",
			roleConfig: RoleConfig{
				AllowedTools: []string{"*_file"},
				BlockedTools: []string{"*_read_*"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isToolAllowedForRole(tt.toolName, tt.serverName, tt.roleConfig)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestServerFiltering(t *testing.T) {
	tests := []struct {
		name       string
		serverName string
		roleConfig RoleConfig
		expected   bool
	}{
		{
			name:       "no restrictions - allow all",
			serverName: "filesystem",
			roleConfig: RoleConfig{},
			expected:   true,
		},
		{
			name:       "blocked server - deny",
			serverName: "network",
			roleConfig: RoleConfig{
				BlockedServers: []string{"network"},
			},
			expected: false,
		},
		{
			name:       "allowed server - allow",
			serverName: "filesystem",
			roleConfig: RoleConfig{
				AllowedServers: []string{"filesystem", "shell"},
			},
			expected: true,
		},
		{
			name:       "allowed server - deny non-matching",
			serverName: "network",
			roleConfig: RoleConfig{
				AllowedServers: []string{"filesystem", "shell"},
			},
			expected: false,
		},
		{
			name:       "blocked takes precedence over allowed",
			serverName: "filesystem",
			roleConfig: RoleConfig{
				AllowedServers: []string{"*"},
				BlockedServers: []string{"filesystem"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isServerAllowedForRole(tt.serverName, tt.roleConfig)
			require.Equal(t, tt.expected, result)
		})
	}
}
