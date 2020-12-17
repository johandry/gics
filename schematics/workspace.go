package schematics

import (
	"fmt"
	"time"
)

const (
	schematicsWorkspaceBaseURL = "https://cloud.ibm.com/schematics/workspaces"
)

// // Status is the status of a Schematics workspace
// type Status int

// const (
// 	StatusNew Status = iota
// 	StatusInactive
// )

// func (s Status) String() string {
// 	return [...]string{"New", "Inactive"}[s]
// }

// Variable encapsulate the parameters for a Schematics Workspace or Terraform variable
type Variable struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Secure      bool   `json:"secure,omitempty" yaml:"secure,omitempty"`
	UseDefault  bool   `json:"use_default,omitempty" yaml:"use_default,omitempty"`
}

// EnvVariable is an environment variable to set in the Schematics Workspace
type EnvVariable map[string]string

// GitRepo encapsulate the parameters for a Git repository (GitHub or GitLab)
// storing a Terraform code to be executed by Schematics
type GitRepo struct {
	Branch       string `json:"branch,omitempty" yaml:"branch,omitempty"`
	Release      string `json:"release,omitempty" yaml:"release,omitempty"`
	RepoShaValue string `json:"repo_sha_value,omitempty" yaml:"repo_sha_value,omitempty"`
	RepoURL      string `json:"repo_url,omitempty" yaml:"repo_url,omitempty"`
	URL          string `json:"url,omitempty" yaml:"url,omitempty"`
	Token        string `json:"token,omitempty" yaml:"token,omitempty"`
}

// WorkspaceStatus is the status of a Schematics workspace
type WorkspaceStatus string

// Workspace is a Schematics workspace
type Workspace struct {
	ID                  string                 `json:"id,omitempty" yaml:"id,omitempty"`
	Name                string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Description         string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Location            string                 `json:"location,omitempty" yaml:"location,omitempty"`
	ResourceGroup       string                 `json:"resource_group,omitempty" yaml:"resource_group,omitempty"`
	Tags                []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	EnvValues           []EnvVariable          `json:"env_values,omitempty" yaml:"env_values,omitempty"`
	Folder              string                 `json:"folder,omitempty" yaml:"folder,omitempty"`
	InitStateFile       string                 `json:"init_state_file,omitempty" yaml:"init_state_file,omitempty"`
	Type                string                 `json:"type,omitempty" yaml:"type,omitempty"`
	UninstallScriptName string                 `json:"uninstall_script_name,omitempty" yaml:"uninstall_script_name,omitempty"`
	Values              string                 `json:"values,omitempty" yaml:"values,omitempty"`
	Variables           []Variable             `json:"variables,omitempty" yaml:"variables,omitempty"`
	GitRepo             *GitRepo               `json:"git_repo,omitempty" yaml:"git_repo,omitempty"`
	CreatedAt           time.Time              `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	CreatedBy           string                 `json:"created_by,omitempty" yaml:"created_by,omitempty"`
	Status              WorkspaceStatus        `json:"status,omitempty" yaml:"status,omitempty"`
	Output              map[string]interface{} `json:"output,omitempty" yaml:"output,omitempty"`
	Code                []byte
	TarCode             []byte

	// Other possible parameters used by the API:
	// TemplateData         TemplateData                   `json:"template_data,omitempty" yaml:"template_data,omitempty"`
	// Type                 []string                       `json:"type,omitempty" yaml:"type,omitempty"`
	// WorkspaceStatus      *WorkspaceStatusResponse       `json:"workspace_status,omitempty"`
	// WorkspaceStatusMsg   *WorkspaceStatusMessage        `json:"workspace_status_msg,omitempty"`
	// Crn                  *string                        `json:"crn,omitempty"`
	// LastHealthCheckAt    *time.Time                     `json:"last_health_check_at,omitempty"`
	// RuntimeData          *[]TemplateRunTimeDataResponse `json:"runtime_data,omitempty"`
	// UpdatedAt            *time.Time                     `json:"updated_at,omitempty"`
	// UpdatedBy            *string                        `json:"updated_by,omitempty"`
	// AppliedShareddataIds []string                       `json:"applied_shareddata_ids,omitempty"`
	// CatalogRef           CatalogRef                     `json:"catalog_ref,omitempty"`
	// SharedData           SharedTargetData               `json:"shared_data,omitempty"`
	// TemplateRef          string                         `json:"template_ref,omitempty"`
	// WorkspaceStatus      WorkspaceStatusRequest         `json:"workspace_status,omitempty"`
}

// type TemplateSourceData struct {
// 	EnvValues           EnvVariable              `json:"env_values,omitempty"`
// 	Folder              string                   `json:"folder,omitempty"`
// 	InitStateFile       string                   `json:"init_state_file,omitempty"`
// 	Type                string                   `json:"type,omitempty"`
// 	UninstallScriptName string                   `json:"uninstall_script_name,omitempty"`
// 	Values              string                   `json:"values,omitempty"`
// 	ValuesMetadata      []map[string]interface{} `json:"values_metadata,omitempty"`
// 	Variablestore       Variables                `json:"variablestore,omitempty"`
// }
// type TemplateData []TemplateSourceData

// New creates an empty Workspace structure which is linked to a Schematics
// workspace and used to execute actions on it
func New(name, description string) *Workspace {
	if len(name) == 0 {
		name = fmt.Sprintf("workspace_%s", time.Now().Format("MM_DD_YYYY"))
	}
	return &Workspace{
		Name:        name,
		Description: description,
	}
}

// AddVar appends a new variable to the workspace settings
func (w *Workspace) AddVar(name, value, varType, description string, secure bool) error {
	if len(name) == 0 {
		return fmt.Errorf("invalid variable name, it cannot be an empty string")
	}
	for _, v := range w.Variables {
		if v.Name == name {
			return fmt.Errorf("Variable %q already exist", name)
		}
	}
	if len(varType) == 0 {
		varType = "string"
	}
	v := Variable{
		Name:        name,
		Value:       value,
		Type:        varType,
		Description: description,
		Secure:      secure,
		UseDefault:  false,
	}

	if len(w.Variables) == 0 {
		w.Variables = []Variable{}
	}

	w.Variables = append(w.Variables, v)

	return nil
}

// Run is used to create, generate and apply the plan of the Schematics
// workspace in a synchronous way, blocking the execution of the code until the
// process is completed or fail
func (w *Workspace) Run() error {
	return nil
}

// Output collect the outputs of the execution of the Schematics
// workspace plan to set the output parameter of the Workspace
// func (w *Workspace) Output() (map[string]interface{}, error) {
// 	w.output = map[string]interface{}{}
// 	return w.output, nil
// }

// Delete ...
func (w *Workspace) Delete(destroy bool) error {
	// if _, err := w.DestroySync(nil); err != nil {
	// 	fmt.Printf("Fail to destroy the resources provisioned by the Schematics Workspace. %s", err)
	// 	fmt.Printf("Please, delete the resources provisioned by the Schematics Workspace manually. The Schematics Workspace URL is: %s", w.URL)
	// 	return
	// }

	// fmt.Printf("The resources provisioned by the Schematics workspace %q (%s) has been destroyed", w.Name, w.ID)

	// if err := w.Delete(); err != nil {
	// 	fmt.Printf("Fail to delete the Schematics Workspace. %s", err)
	// 	fmt.Printf("Please, delete the Schematics Workspace manually. The Schematics Workspace URL is: %s", w.URL)
	// 	return
	// }

	// fmt.Printf("Schematics workspace %q (%s) has been deleted", w.Name, w.ID)

	return nil
}
