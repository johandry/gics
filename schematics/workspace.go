package schematics

import (
	"context"
	"fmt"
	"io"
	"time"

	apiv1 "github.com/johandry/gics/schematics/api/v1"
)

const (
	schematicsWorkspaceBaseURL = "https://cloud.ibm.com/schematics/workspaces"
	templateIDDefault          = "terraform_v0.13"
)

const (
	deleteWorkspaceTimeout = 50
)

// Variable encapsulate the parameters for a Schematics Workspace or Terraform variable
type Variable struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Secure      bool   `json:"secure,omitempty" yaml:"secure,omitempty"`
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

	Code    []byte
	TarCode []byte

	tfCodeFiles map[string]string
	tfBuf       io.Reader
	service     *Service
	logOutput   io.Writer

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
func New(name, description string, service *Service) *Workspace {
	if len(name) == 0 {
		name = fmt.Sprintf("workspace_%s", time.Now().Format("MM_DD_YYYY"))
	}
	if service == nil {
		service = defaultService
	}
	return &Workspace{
		Name:        name,
		Description: description,
		Status:      WorkspaceStatus("NEW"),
		service:     service,
		Type:        templateIDDefault,
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
	}

	if len(w.Variables) == 0 {
		w.Variables = []Variable{}
	}

	w.Variables = append(w.Variables, v)

	return nil
}

// AddRepo assign a Git URL from GitHub, GitLab, BitBucket or any other supported
// by Schematics, to the Workspace
func (w *Workspace) AddRepo(url string) {
	if w.GitRepo != nil {
		w.GitRepo.URL = url
		return
	}

	w.GitRepo = &GitRepo{
		URL: url,
	}

	return
}

// SetOutput sets the output used for the logger. It won't log by default
func (w *Workspace) SetOutput(out io.Writer) {
	w.logOutput = out
}

// LoadCode tar and loads the given code to the workspace
func (w *Workspace) LoadCode(code string) error {
	w.tfCodeFiles = map[string]string{
		"main.tf": code,
	}

	r, err := w.tarMemFiles()
	if err != nil {
		return fmt.Errorf("failed to tar the files in memory. %s", err)
	}

	w.tfBuf = r

	return nil
}

// Run is used to create, generate and apply the plan of the Schematics
// workspace in a synchronous way, blocking the execution of the code until the
// process is completed or fail
func (w *Workspace) Run() error {
	// Create the Schematics workspace
	actCreate, err := w.Create()
	if err != nil {
		return err
	}
	// the activity should be a NilActivity, anyway we wait in case the API change
	// in the future
	if err := actCreate.Wait(); err != nil {
		return err
	}

	if err := w.UploadTar(w.tfBuf); err != nil {
		return err
	}

	// Generate the workspace plan
	actPlan, err := w.Plan()
	if err != nil {
		return err
	}
	if err := actPlan.Wait(); err != nil {
		return err
	}

	// Apply the workspace plan
	actApply, err := w.Apply()
	if err != nil {
		return err
	}
	if err := actApply.Wait(); err != nil {
		return err
	}

	return nil
}

// GetParam collect and returns the requested output parameters of the execution
// of the Schematics workspace
func (w *Workspace) GetParam(params ...string) map[string]interface{} {
	output := map[string]interface{}{}
	for _, key := range params {
		if value, ok := w.Output[key]; ok {
			output[key] = value
		}
	}
	return output
}

// Delete deletes an existing Schematics workspace, it may destroy the resources
// and wait for them to be destroyed
func (w *Workspace) Delete(destroy bool) error {
	if !destroy {
		// TODO: Should we return an error if there are resources created?
		// 			 Maybe not, it may be possible to request to delete the workspace
		// 			 and not the resources
		return w.delete(false)
	}

	actDestroy, err := w.Destroy()
	if err != nil {
		return err
	}

	return actDestroy.Wait()
}

// delete deletes an existing Schematics workspace, it also destroy the resources
// if the `destroy` parameter is set to true.
func (w *Workspace) delete(destroy bool) error {
	// Delete Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), deleteWorkspaceTimeout*time.Second)
	defer cancelFunc()

	token, err := w.getToken()
	if err != nil {
		return err
	}

	params := &apiv1.DeleteWorkspaceParams{
		DestroyResources: &destroy,
		RefreshToken:     &token,
	}
	resp, err := w.service.clientWithResponses.DeleteWorkspaceWithResponse(ctx, w.ID, params)
	if err != nil {
		return err
	}
	if code := resp.StatusCode(); code != 200 {
		return getAPIError("failed to delete the workspace", resp.Body)
	}

	// response := resp.JSON200 // *WorkspaceDeleteResponse => *String
	// fmt.Printf("[DEBUG] Delete response: %+v\n", *response)

	w.Status = WorkspaceStatusDeleted

	return nil
}

// ListResources returns the list of resources created by the workspace
func (w *Workspace) ListResources() ([]string, error) {
	return []string{}, nil
}
