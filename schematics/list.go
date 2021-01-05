package schematics

import (
	"context"
	"time"

	apiv1 "github.com/johandry/gics/schematics/api/v1"
)

const listTimeout = 50

// WorkspaceSummary is a data structure for the Schematic workspace information
// to return a list of workspaces
type WorkspaceSummary struct {
	ID          string    `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string    `json:"name,omitempty" yaml:"name,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Location    string    `json:"location,omitempty" yaml:"location,omitempty"`
	Owner       string    `json:"owner,omitempty" yaml:"owner,omitempty"`
	State       string    `json:"state,omitempty" yaml:"state,omitempty"`
	Created     time.Time `json:"created,omitempty" yaml:"created,omitempty"`
}

// WorkspaceList is a data structure for a list of Schematics workspaces used
// by List()
type WorkspaceList struct {
	Workspaces []WorkspaceSummary `json:"workspaces,omitempty" yaml:"workspaces,omitempty"`
}

// map[count:1 limit:100 offset:0
// 	workspaces:[map[applied_shareddata_ids:<nil> created_at:2020-11-01T18:30:15.364093434Z created_by:johandry@gmail.com crn:crn:v1:bluemix:public:schematics:us-south:a/f069279a778a48178fc379e99797f887:ba49546a-1075-4e1d-bfde-9d5ef3ca031e:workspace:myworkspace-270d8f90-b67a-46 description: id:myworkspace-270d8f90-b67a-46 last_health_check_at:0001-01-01T00:00:00Z location:us-east name:myworkspace resource_group:Default runtime_data:[map[engine_name:terraform engine_version:v0.12.20 id:e86f76d0-b77f-49 log_store_url:https://schematics.cloud.ibm.com/v1/workspaces/myworkspace-270d8f90-b67a-46/runtime_data/e86f76d0-b77f-49/log_store state_store_url:https://schematics.cloud.ibm.com/v1/workspaces/myworkspace-270d8f90-b67a-46/runtime_data/e86f76d0-b77f-49/state_store]] shared_data:map[resource_group_id:] status:FAILED tags:[] template_data:[map[folder:iac/02-schematics has_githubtoken:false id:e86f76d0-b77f-49 type:terraform_v0.12 values: values_metadata:[map[description: name:project_name type:string] map[description: name:environment type:string] map[description: name:public_key type:string] map[default:8080 description: name:port type:string]] values_url:https://schematics.cloud.ibm.com/v1/workspaces/myworkspace-270d8f90-b67a-46/template_data/e86f76d0-b77f-49/values]] template_repo:map[branch:master full_url:https://github.com/IBM/cloud-enterprise-examples/tree/master/iac/02-schematics has_uploadedgitrepotar:false url:https://github.com/IBM/cloud-enterprise-examples] type:[terraform_v0.12] updated_at:2020-11-01T18:33:05.807531893Z updated_by:johandry@gmail.com workspace_status:map[frozen:false locked:false] workspace_status_msg:map[status_code:500 status_msg:rpc error: code = ResourceExhausted desc = trying to send message larger than max (129563174 vs. 83886080)]]]
// ]

// List return a list of existing Schematics workspaces in your IBM Cloud account
// using the default Schematics service
func List() (*WorkspaceList, error) {
	ctx := context.Background()
	return defaultService.List(ctx)
}

// List return a list of existing Schematics workspaces in your IBM Cloud account
func (s *Service) List(ctx context.Context) (*WorkspaceList, error) {
	// List Timeout
	ctx, cancelFunc := context.WithTimeout(ctx, listTimeout*time.Second)
	defer cancelFunc()

	resp, err := s.clientWithResponses.ListWorkspacesWithResponse(ctx, &apiv1.ListWorkspacesParams{})
	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode(); code != 200 {
		return nil, getAPIError("failed to list the workspaces", resp.Body)
	}
	response := resp.JSON200

	if response.Workspaces == nil {
		return &WorkspaceList{
			Workspaces: []WorkspaceSummary{},
		}, nil
	}

	workspaces := []WorkspaceSummary{}
	for _, wksp := range *response.Workspaces {
		workspace := WorkspaceSummary{
			ID:          *wksp.Id,
			Name:        *wksp.Name,
			Description: *wksp.Description,
			Location:    *wksp.Location,
			Owner:       *wksp.CreatedBy,
			State:       (string)(*wksp.Status),
			Created:     *wksp.CreatedAt,
		}
		workspaces = append(workspaces, workspace)
	}

	return &WorkspaceList{
		Workspaces: workspaces,
	}, nil
}
