package domain

// TemplateType represents the type of server template
type TemplateType string

const (
	TemplateTypeBasic      TemplateType = "basic"
	TemplateTypeRoleplay   TemplateType = "roleplay"
	TemplateTypeDrifting   TemplateType = "drifting"
	TemplateTypeRacing     TemplateType = "racing"
	TemplateTypeDeathmatch TemplateType = "deathmatch"
	TemplateTypeCustom     TemplateType = "custom"
)

// Template represents a server configuration template
type Template struct {
	Name         string            `json:"name"`
	Type         TemplateType      `json:"type"`
	Description  string            `json:"description"`
	Version      string            `json:"version"`
	Author       string            `json:"author"`
	Resources    []string          `json:"resources"`    // List of required resources
	Config       map[string]string `json:"config"`       // Server.cfg values
	ConVars      map[string]string `json:"convars"`      // Server convars
	Permissions  map[string]string `json:"permissions"`  // ACL permissions
	Requirements Requirements      `json:"requirements"` // System requirements
}

// Requirements represents the system requirements for a template
type Requirements struct {
	MinRAM      int      `json:"min_ram_mb"`      // Minimum RAM in MB
	MinCPU      int      `json:"min_cpu_cores"`   // Minimum CPU cores
	MinStorage  int      `json:"min_storage_mb"`  // Minimum storage in MB
	Ports       []int    `json:"ports"`           // Required ports
	Database    bool     `json:"database"`        // Requires database
	Dependencies []string `json:"dependencies"`    // External dependencies
}

// NewTemplate creates a new template instance
func NewTemplate(name string, templateType TemplateType) *Template {
	return &Template{
		Name:        name,
		Type:        templateType,
		Resources:   []string{},
		Config:      make(map[string]string),
		ConVars:     make(map[string]string),
		Permissions: make(map[string]string),
		Requirements: Requirements{
			MinRAM:     2048,  // 2GB default
			MinCPU:     2,     // 2 cores default
			MinStorage: 5120,  // 5GB default
			Ports:      []int{30120, 30110}, // Default FiveM ports
		},
	}
}

// AddResource adds a required resource to the template
func (t *Template) AddResource(resourceName string) {
	for _, r := range t.Resources {
		if r == resourceName {
			return // Already exists
		}
	}
	t.Resources = append(t.Resources, resourceName)
}

// SetConfig sets a configuration value
func (t *Template) SetConfig(key, value string) {
	t.Config[key] = value
}

// SetConVar sets a server convar
func (t *Template) SetConVar(key, value string) {
	t.ConVars[key] = value
}

// SetPermission sets an ACL permission
func (t *Template) SetPermission(key, value string) {
	t.Permissions[key] = value
}

// GetDefaultTemplates returns a list of default templates
func GetDefaultTemplates() map[string]*Template {
	templates := make(map[string]*Template)

	// Basic template
	basic := NewTemplate("Basic Server", TemplateTypeBasic)
	basic.Description = "A minimal FiveM server setup"
	basic.AddResource("chat")
	basic.AddResource("spawnmanager")
	basic.AddResource("sessionmanager")
	basic.AddResource("hardcap")
	basic.SetConfig("sv_hostname", "My FiveM Server")
	basic.SetConfig("sv_maxclients", "32")
	templates["basic"] = basic

	// Roleplay template
	roleplay := NewTemplate("Roleplay Server", TemplateTypeRoleplay)
	roleplay.Description = "A roleplay-focused server with essential RP resources"
	roleplay.AddResource("chat")
	roleplay.AddResource("spawnmanager")
	roleplay.AddResource("sessionmanager")
	roleplay.AddResource("hardcap")
	roleplay.SetConfig("sv_hostname", "My Roleplay Server")
	roleplay.SetConfig("sv_maxclients", "64")
	roleplay.Requirements.MinRAM = 4096
	roleplay.Requirements.MinCPU = 4
	templates["roleplay"] = roleplay

	// Drifting template
	drifting := NewTemplate("Drifting Server", TemplateTypeDrifting)
	drifting.Description = "A server optimized for drifting and car meets"
	drifting.AddResource("chat")
	drifting.AddResource("spawnmanager")
	drifting.AddResource("sessionmanager")
	drifting.AddResource("hardcap")
	drifting.SetConfig("sv_hostname", "My Drift Server")
	drifting.SetConfig("sv_maxclients", "48")
	templates["drifting"] = drifting

	return templates
}