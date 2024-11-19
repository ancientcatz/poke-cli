package main

import (
	"bytes"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = stdout
	_, _ = buf.ReadFrom(r)

	return buf.String()
}

func TestRunCLI(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedCode   int
	}{
		{
			name: "No Arguments",
			args: []string{},
			expectedOutput: "╭──────────────────────────────────────────────────────╮\n" +
				"│Welcome! This tool displays data related to Pokémon!  │\n" +
				"│                                                      │\n" +
				"│ USAGE:                                               │\n" +
				"│    poke-cli [flag]                                   │\n" +
				"│    poke-cli <command> [flag]                         │\n" +
				"│    poke-cli <command> <subcommand> [flag]            │\n" +
				"│                                                      │\n" +
				"│ FLAGS:                                               │\n" +
				"│    -h, --help      Shows the help menu               │\n" +
				"│    -l, --latest    Prints the latest available       │\n" +
				"│                    version of the program            │\n" +
				"│                                                      │\n" +
				"│ AVAILABLE COMMANDS:                                  │\n" +
				"│    pokemon         Get details of a specific Pokémon │\n" +
				"│    types           Get details of a specific typing  │\n" +
				"╰──────────────────────────────────────────────────────╯\n",
			expectedCode: 0,
		},
		{
			name: "Help Flag Short",
			args: []string{"-h"},
			expectedOutput: "╭──────────────────────────────────────────────────────╮\n" +
				"│Welcome! This tool displays data related to Pokémon!  │\n" +
				"│                                                      │\n" +
				"│ USAGE:                                               │\n" +
				"│    poke-cli [flag]                                   │\n" +
				"│    poke-cli <command> [flag]                         │\n" +
				"│    poke-cli <command> <subcommand> [flag]            │\n" +
				"│                                                      │\n" +
				"│ FLAGS:                                               │\n" +
				"│    -h, --help      Shows the help menu               │\n" +
				"│    -l, --latest    Prints the latest available       │\n" +
				"│                    version of the program            │\n" +
				"│                                                      │\n" +
				"│ AVAILABLE COMMANDS:                                  │\n" +
				"│    pokemon         Get details of a specific Pokémon │\n" +
				"│    types           Get details of a specific typing  │\n" +
				"╰──────────────────────────────────────────────────────╯\n",
			expectedCode: 0,
		},
		{
			name: "Help Flag Long",
			args: []string{"--help"},
			expectedOutput: "╭──────────────────────────────────────────────────────╮\n" +
				"│Welcome! This tool displays data related to Pokémon!  │\n" +
				"│                                                      │\n" +
				"│ USAGE:                                               │\n" +
				"│    poke-cli [flag]                                   │\n" +
				"│    poke-cli <command> [flag]                         │\n" +
				"│    poke-cli <command> <subcommand> [flag]            │\n" +
				"│                                                      │\n" +
				"│ FLAGS:                                               │\n" +
				"│    -h, --help      Shows the help menu               │\n" +
				"│    -l, --latest    Prints the latest available       │\n" +
				"│                    version of the program            │\n" +
				"│                                                      │\n" +
				"│ AVAILABLE COMMANDS:                                  │\n" +
				"│    pokemon         Get details of a specific Pokémon │\n" +
				"│    types           Get details of a specific typing  │\n" +
				"╰──────────────────────────────────────────────────────╯\n",
			expectedCode: 0,
		},
		{
			name:           "Invalid Command",
			args:           []string{"invalid"},
			expectedOutput: "Error!",
			expectedCode:   1,
		},
		{
			name:           "Latest Flag",
			args:           []string{"-l"},
			expectedOutput: "Latest Docker image version: v0.7.1\nLatest release tag: v0.7.1\n",
			expectedCode:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exit = func(code int) {}
			output := captureOutput(func() {
				code := runCLI(tt.args)
				if code != tt.expectedCode {
					t.Errorf("Expected exit code %d, got %d", tt.expectedCode, code)
				}
			})

			if !bytes.Contains([]byte(output), []byte(tt.expectedOutput)) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}
