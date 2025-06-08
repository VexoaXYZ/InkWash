package services

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/vexoa/inkwash/internal/domain"
)

// templateServiceImpl implements TemplateService
type templateServiceImpl struct {
	fileService   FileService
	templatesDir  string
	defaultTemplates map[string]*domain.Template
}

// NewTemplateService creates a new template service
func NewTemplateService(fileService FileService, templatesDir string) TemplateService {
	service := &templateServiceImpl{
		fileService:      fileService,
		templatesDir:     templatesDir,
		defaultTemplates: domain.GetDefaultTemplates(),
	}

	return service
}

// GetTemplate gets a template by name
func (s *templateServiceImpl) GetTemplate(ctx context.Context, name string) (*domain.Template, error) {
	// Check default templates first
	if template, exists := s.defaultTemplates[name]; exists {
		return template, nil
	}

	// Check custom templates
	templatePath := filepath.Join(s.templatesDir, name+".json")
	if !s.fileService.FileExists(templatePath) {
		return nil, domain.NewError(domain.ErrorTypeNotFound, "template not found").
			WithDetail("template_name", name)
	}

	data, err := s.fileService.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	var template domain.Template
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, domain.NewError(domain.ErrorTypeInternal, "failed to parse template").WithCause(err)
	}

	return &template, nil
}

// ListTemplates lists all available templates
func (s *templateServiceImpl) ListTemplates(ctx context.Context) ([]*domain.Template, error) {
	var templates []*domain.Template

	// Add default templates
	for _, template := range s.defaultTemplates {
		templates = append(templates, template)
	}

	// Add custom templates
	if s.fileService.FileExists(s.templatesDir) {
		entries, err := s.fileService.ListDirectory(s.templatesDir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if filepath.Ext(entry) == ".json" {
				templateName := entry[:len(entry)-5] // Remove .json extension
				template, err := s.GetTemplate(ctx, templateName)
				if err != nil {
					continue // Skip invalid templates
				}
				templates = append(templates, template)
			}
		}
	}

	return templates, nil
}

// ApplyTemplate applies a template to a server
func (s *templateServiceImpl) ApplyTemplate(ctx context.Context, serverID string, templateName string) error {
	template, err := s.GetTemplate(ctx, templateName)
	if err != nil {
		return err
	}

	// Template application logic would go here
	// For now, this is a placeholder
	_ = template
	return nil
}

// CreateTemplate creates a custom template
func (s *templateServiceImpl) CreateTemplate(ctx context.Context, template *domain.Template) error {
	// Create templates directory if it doesn't exist
	if err := s.fileService.CreateDirectory(s.templatesDir, 0755); err != nil {
		return err
	}

	templatePath := filepath.Join(s.templatesDir, template.Name+".json")
	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return domain.NewError(domain.ErrorTypeInternal, "failed to marshal template").WithCause(err)
	}

	return s.fileService.WriteFile(templatePath, data, 0644)
}

// UpdateTemplate updates a template
func (s *templateServiceImpl) UpdateTemplate(ctx context.Context, template *domain.Template) error {
	// Don't allow updating default templates
	if _, exists := s.defaultTemplates[template.Name]; exists {
		return domain.NewError(domain.ErrorTypeValidation, "cannot update default template")
	}

	return s.CreateTemplate(ctx, template)
}

// DeleteTemplate deletes a custom template
func (s *templateServiceImpl) DeleteTemplate(ctx context.Context, templateName string) error {
	// Don't allow deleting default templates
	if _, exists := s.defaultTemplates[templateName]; exists {
		return domain.NewError(domain.ErrorTypeValidation, "cannot delete default template")
	}

	templatePath := filepath.Join(s.templatesDir, templateName+".json")
	return s.fileService.DeleteFile(templatePath)
}

// ExportTemplate exports a server configuration as a template
func (s *templateServiceImpl) ExportTemplate(ctx context.Context, serverID string, templateName string) (*domain.Template, error) {
	// This would read server configuration and create a template
	// For now, this is a placeholder
	template := domain.NewTemplate(templateName, domain.TemplateTypeCustom)
	template.Description = "Exported from server " + serverID
	
	return template, nil
}