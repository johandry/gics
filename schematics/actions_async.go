package schematics

import (
	"context"
	"fmt"
	"time"

	apiv1 "github.com/johandry/gics/schematics/api/v1"
)

const (
	createWorkspaceTimeout = 50
)

// Create creates a Schematics Workspace and returns the activity in charge of
// this task
func (w *Workspace) Create() (*Activity, error) {
	w.Output = nil

	// Create Timeout
	ctx, cancelFunc := context.WithTimeout(context.Background(), createWorkspaceTimeout*time.Second)
	defer cancelFunc()

	variables := []apiv1.WorkspaceVariableRequest{}
	for _, v := range w.Variables {
		variable := apiv1.WorkspaceVariableRequest{
			Description: &v.Description,
			Name:        &v.Name,
			Secure:      &v.Secure,
			Type:        &v.Type,
			UseDefault:  &v.UseDefault,
			Value:       &v.Value,
		}
		variables = append(variables, variable)
	}
	variableStore := apiv1.VariablesRequest(variables)

	templateData := &apiv1.TemplateData{
		apiv1.TemplateSourceDataRequest{
			Folder:              &w.Folder,
			InitStateFile:       &w.InitStateFile,
			Type:                &w.Type,
			UninstallScriptName: &w.UninstallScriptName,
			Values:              &w.Values,
			Variablestore:       &variableStore,
		},
	}
	var templateRepo *apiv1.TemplateRepoRequest
	if w.GitRepo != nil {
		templateRepo = &apiv1.TemplateRepoRequest{
			Branch:       &w.GitRepo.Branch,
			Release:      &w.GitRepo.Release,
			RepoShaValue: &w.GitRepo.RepoShaValue,
			RepoUrl:      &w.GitRepo.RepoURL,
			Url:          &w.GitRepo.URL,
		}
	}
	workspaceCreateRequest := apiv1.WorkspaceCreateRequest{
		Description:   &w.Description,
		Location:      &w.Location,
		Name:          &w.Name,
		ResourceGroup: &w.ResourceGroup,
		Tags:          &w.Tags,
		TemplateData:  templateData,
		TemplateRepo:  templateRepo,
		Type:          &[]string{w.Type},
	}

	params := &apiv1.CreateWorkspaceParams{}
	body := apiv1.CreateWorkspaceJSONRequestBody(apiv1.CreateWorkspaceJSONBody(workspaceCreateRequest))
	resp, err := w.service.clientWithResponses.CreateWorkspaceWithResponse(ctx, params, body)
	if err != nil {
		return nil, err
	}
	if code := resp.StatusCode(); code != 201 {
		return nil, fmt.Errorf(`{"status_code": %d, "status": %q}`, code, resp.Status())
	}
	response := resp.JSON201 // WorkspaceResponse

	if response.CreatedAt != nil {
		w.CreatedAt = *response.CreatedAt
	}

	w.CreatedBy = stringValue(response.CreatedBy)
	w.Description = stringValue(response.Description)
	w.ID = stringValue(response.Id)
	w.Location = stringValue(response.Location)
	w.Name = stringValue(response.Name)
	w.ResourceGroup = stringValue(response.ResourceGroup)

	t := *response.Type
	if len(t) > 0 {
		w.Type = t[0]
	}

	if response.Status != nil {
		w.Status = WorkspaceStatus(*response.Status)
	}
	if response.Tags != nil {
		w.Tags = *response.Tags
	}

	if response.TemplateRepo != nil {
		w.GitRepo = &GitRepo{
			URL:          stringValue(response.TemplateRepo.Url),
			RepoURL:      stringValue(response.TemplateRepo.RepoUrl),
			Branch:       stringValue(response.TemplateRepo.Branch),
			Release:      stringValue(response.TemplateRepo.Release),
			RepoShaValue: stringValue(response.TemplateRepo.RepoShaValue),
		}
	}

	// EnvValues *EnvVariableRequest `json:"env_values,omitempty"`
	if len(*response.TemplateData) > 0 {
		templateData := *response.TemplateData

		w.Folder = stringValue(templateData[0].Folder)
		w.UninstallScriptName = stringValue(templateData[0].UninstallScriptName)
		w.Values = stringValue(templateData[0].Values)

		wv := []apiv1.WorkspaceVariableResponse(*templateData[0].Variablestore)
		var variables []Variable
		if len(wv) > 0 {
			variables = []Variable{}
		}
		for _, v := range wv {
			variable := Variable{
				Name:        *v.Name,
				Value:       *v.Value,
				Type:        *v.Type,
				Description: *v.Description,
				Secure:      *v.Secure,
			}
			variables = append(variables, variable)
		}
		w.Variables = variables
	}

	act, err := w.LastActivity(activityNameForCreate)
	if err != nil {
		return nil, err
	}

	// TODO: Verify this is the correct status. It may not be.
	// w.Status = WorkspaceStatus(act.Status)

	return act, err
}

// Plan executes the planning of the Schematics Workspace
func (w *Workspace) Plan() (*Activity, error) {
	w.Output = nil
	// TODO: Complete the Plan method of Workspace
	return &Activity{}, nil
}

// Apply executes the applying of the Schematics Workspace. It 'executes' the
// Terraform code in the workspace
func (w *Workspace) Apply() (*Activity, error) {
	w.Output = nil
	// TODO: Complete the Apply method of Workspace
	return &Activity{}, nil
}

// Destroy destroyes the resources created by the Terraform code in the Schematics
// Workspace, it does not delete the workspace
func (w *Workspace) Destroy() (*Activity, error) {
	// TODO: Complete the Destroy method of Workspace
	return &Activity{}, nil
}

// LastActivity returns the last executed activity. It should be call after every
// action (i.e. Plan, Apply) to return the action activity.
func (w *Workspace) LastActivity(name string) (*Activity, error) {
	// Get all the activities of the workspace
	activities, err := getActivities(w.service, w.ID)
	if err != nil {
		return nil, err
	}
	if activities == nil || len(activities) == 0 {
		return nil, nil
	}

	// Filter the activities by Name and PerformedBy
	activitiesWithName := []Activity{}
	for _, act := range activities {
		if (act.Name == name) && (act.PerformedBy == w.CreatedBy) {
			activitiesWithName = append(activitiesWithName, act)
		}
	}

	// Trying to safe some time
	l := len(activitiesWithName)
	if l == 0 {
		return nil, nil
	}
	if l == 1 {
		return &activitiesWithName[0], nil
	}

	var activity Activity

	// Get the latest activity if there is more than one
	now := time.Now()
	t0, _ := time.Parse(time.RFC822, "31 Oct 20 15:30 PST") // <- GICS Birthday
	diff := now.Sub(t0)
	for _, act := range activitiesWithName {
		d := now.Sub(act.PerformedAt)
		if d < diff {
			diff = d
			activity = act
		}
	}

	return &activity, nil
}
