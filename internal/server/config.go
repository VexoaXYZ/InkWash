package server

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/VexoaXYZ/inkwash/pkg/types"
)

const serverConfigTemplate = `## Server Identity
sv_hostname "{{.ServerName}}"
sv_licenseKey "{{.LicenseKey}}"
sv_maxclients {{.MaxPlayers}}

## Server Configuration
endpoint_add_tcp "0.0.0.0:{{.Port}}"
endpoint_add_udp "0.0.0.0:{{.Port}}"

## Resources
ensure mapmanager
ensure chat
ensure spawnmanager
ensure sessionmanager
ensure basic-gamemode
ensure hardcap
ensure rconlog

## Permissions
add_ace resource.* command allow

## Server Info
sets sv_projectName "{{.ServerName}}"
sets sv_projectDesc "FiveM Server powered by Inkwash"
sets tags "inkwash"

## Logging
set sv_logFile "logs/server.log"
`

// ConfigGenerator generates server configuration files
type ConfigGenerator struct{}

// NewConfigGenerator creates a new config generator
func NewConfigGenerator() *ConfigGenerator {
	return &ConfigGenerator{}
}

// GenerateServerConfig generates a server.cfg file
func (cg *ConfigGenerator) GenerateServerConfig(server *types.Server, licenseKey string) error {
	tmpl, err := template.New("server.cfg").Parse(serverConfigTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	configPath := filepath.Join(server.Path, "server.cfg")
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	data := struct {
		ServerName  string
		LicenseKey  string
		MaxPlayers  int
		Port        int
	}{
		ServerName: server.Name,
		LicenseKey: licenseKey,
		MaxPlayers: 32,
		Port:       server.Port,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	return nil
}

// GenerateLaunchScript generates platform-specific launch script
func (cg *ConfigGenerator) GenerateLaunchScript(server *types.Server) error {
	scriptPath, scriptContent := cg.getScriptTemplate(server)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to create launch script: %w", err)
	}

	return nil
}

// getScriptTemplate returns the script path and content for the platform
func (cg *ConfigGenerator) getScriptTemplate(server *types.Server) (string, string) {
	if isWindows() {
		scriptPath := filepath.Join(server.Path, "run.cmd")
		content := fmt.Sprintf(`@echo off
cd /d "%s"
bin\FXServer.exe +exec server.cfg
`, server.Path)
		return scriptPath, content
	}

	// Linux
	scriptPath := filepath.Join(server.Path, "run.sh")
	content := fmt.Sprintf(`#!/bin/bash
cd "%s"
bash bin/run.sh +exec server.cfg
`, server.Path)
	return scriptPath, content
}

func isWindows() bool {
	return os.PathSeparator == '\\'
}
