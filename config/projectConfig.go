package config

// ProjectConfig contains project settings
type ProjectConfig struct {
	Path          string `json:"path"`
	Name          string `json:"name"`
	MainContainer string `json:"main_container"`
}

// GetPath gets path to project
func (p *ProjectConfig) GetPath() string {
	return p.Path
}

// GetName gets project name
func (p *ProjectConfig) GetName() string {
	return p.Name
}

// GetMainContainer gets php container name
func (p *ProjectConfig) GetMainContainer() string {
	return p.MainContainer
}
