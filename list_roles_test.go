package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListRolesOutput(t *testing.T) {
	// Save original config and stdout
	originalConfig := config
	originalStdout := os.Stdout

	defer func() {
		config = originalConfig
		os.Stdout = originalStdout
	}()

	// Set up test config
	config = Config{
		Role: "shell",
		Roles: map[string]RoleConfig{
			"default": {
				Prompt: []string{},
			},
			"shell": {
				Prompt:         []string{"you are a shell expert"},
				AllowedServers: []string{"filesystem", "shell"},
				BlockedServers: []string{"network"},
				AllowedTools:   []string{"filesystem_*", "shell_*"},
				BlockedTools:   []string{"dangerous_*"},
			},
			"analyst": {
				Prompt:         []string{"you analyze data"},
				AllowedTools:   []string{"*_read", "*_list"},
				BlockedServers: []string{"network"},
			},
			"simple": {
				Prompt: []string{"just a simple role"},
			},
		},
	}

	// Capture output
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	listRoles()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)

	output := buf.String()

	// Check that all roles are listed
	require.Contains(t, output, "default")
	require.Contains(t, output, "shell (default)")
	require.Contains(t, output, "analyst")
	require.Contains(t, output, "simple")

	// Check MCP configuration is shown
	require.Contains(t, output, "allowed servers: filesystem, shell")
	require.Contains(t, output, "blocked servers: network")
	require.Contains(t, output, "allowed tools: filesystem_*, shell_*")
	require.Contains(t, output, "blocked tools: dangerous_*")

	// Check analyst role MCP config
	require.Contains(t, output, "allowed tools: *_read, *_list")

	// Verify the output format looks reasonable
	fmt.Println("Captured output:")
	fmt.Println(output)
}
