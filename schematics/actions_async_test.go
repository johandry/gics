package schematics

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestWorkspace_Create(t *testing.T) {
	workspaceName := "workspace"
	workspaceID := fmt.Sprintf("%s-2b1cf4ac-348f-49", workspaceName)

	httpmock.RegisterResponder("GET", `=~^https://schematics\.cloud\.ibm\.com/v1/workspaces/([\w-]+)/actions\z`,
		func(req *http.Request) (*http.Response, error) {
			// Get the fixture with the following code after getting the Token:
			// export TOKEN=$(cat .token | jq -r .access_token)
			// export WID=
			// curl -X GET "https://schematics.cloud.ibm.com/v1/workspaces/$WID/actions" -H "Authorization: Bearer $TOKEN" -H 'Accept: application/json' -H 'Content-Type: application/json'
			wID := httpmock.MustGetSubmatch(req, 1)
			assert.NotEmpty(t, wID, "Schematics Workspace should not be empty")
			assert.Equalf(t, workspaceID, wID, "The received Schematics Workspace is not %q", workspaceID)
			fixture := fmt.Sprintf(`{"workspace_name":"%s","workspace_id":"%s","actions":[{"action_id":"24481fe0f1fcd27e5642b8cc0c241eb4","name":"WORKSPACE_CREATE","status":"FAILED","message":[],"performed_by":"johandry@gmail.com","performed_at":"2020-12-17T06:21:32.573626313Z","templates":[{"template_id":"iac-f6ee24a6-1775-42","template_type":"terraform_v0.13","start_time":"2020-12-17T06:21:32.779253562Z","end_time":"2020-12-17T06:24:34.59186019Z","status":"FAILED","message":"{\"messagekey\":\"M2000_InternalError\",\"parms\":{},\"requestid\":\"de8ced31-e98f-4364-bb0d-202d6b52c864\",\"timestamp\":\"2020-12-17T06:24:34.591850876Z\"}","log_url":"https://schematics.cloud.ibm.com/v1/workspaces/%s/runtime_data/iac-f6ee24a6-1775-42/log_store/actions/24481fe0f1fcd27e5642b8cc0c241eb4","log_summary":{"activity_status":"FAILED","time_taken":181.81}}]}]}`, workspaceName, wID, wID)

			resp := httpmock.NewStringResponse(200, fixture)
			resp.Header.Add("Content-Type", "application/json; charset=utf-8")
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", "https://schematics.cloud.ibm.com/v1/workspaces",
		func(req *http.Request) (*http.Response, error) {
			// Get the fixture with the following code after getting the Token:
			// export TOKEN=$(cat .token | jq -r .access_token)
			// curl --request POST --url https://schematics.cloud.ibm.com/v1/workspaces -H "Authorization: Bearer $TOKEN" -H 'Accept: application/json' -H 'Content-Type: application/json' -d '{"name":"workspace","type": ["terraform_v0.13"],"template_repo": {"url": "https://github.com/IBM/cloud-enterprise-examples/tree/master/iac/01-getting-started"},"template_data": [{"folder": ".","type": "terraform_v0.13","variablestore": [{"name": "project_name","value": "gics"},{"name": "environment","value": "testing"},{"name": "public_key","value": "fakepublicsshkey", "secure": true}]}]}'
			fixture := fmt.Sprintf(`{"id":"%s","name":"%s","crn":"crn:v1:bluemix:public:schematics:us-south:a/f069279a778a48178fc379e99797f887:ba49546a-1075-4e1d-bfde-9d5ef3ca031e:workspace:%s","type":["terraform_v0.13"],"description":"","resource_group":"Default","location":"us-south","tags":[],"created_at":"2020-12-17T06:21:29.762423059Z","created_by":"johandry@gmail.com","status":"DRAFT","workspace_status_msg":{"status_code":"","status_msg":""},"workspace_status":{"frozen":false,"locked":false},"template_repo":{"url":"https://github.com/IBM/cloud-enterprise-examples","branch":"master","full_url":"https://github.com/IBM/cloud-enterprise-examples/tree/master/iac/01-getting-started","has_uploadedgitrepotar":false},"template_data":[{"id":"iac-f6ee24a6-1775-42","folder":"iac/01-getting-started","type":"terraform_v0.13","values_url":"https://schematics.cloud.ibm.com/v1/workspaces/%s/template_data/iac-f6ee24a6-1775-42/values","values":"","variablestore":[{"name":"project_name","secure":false,"value":"gics","type":"","description":""},{"name":"environment","secure":false,"value":"testing","type":"","description":""},{"name":"public_key","secure":true,"value":"fakepublicsshkey","type":"","description":""}],"has_githubtoken":false}],"runtime_data":[{"id":"iac-f6ee24a6-1775-42","engine_name":"terraform","engine_version":"v0.11.14","state_store_url":"https://schematics.cloud.ibm.com/v1/workspaces/%s/runtime_data/iac-f6ee24a6-1775-42/state_store","log_store_url":"https://schematics.cloud.ibm.com/v1/workspaces/%s/runtime_data/iac-f6ee24a6-1775-42/log_store"}],"shared_data":{"resource_group_id":""},"applied_shareddata_ids":null,"updated_at":"0001-01-01T00:00:00Z","last_health_check_at":"0001-01-01T00:00:00Z"}`, workspaceID, workspaceName, workspaceID, workspaceID, workspaceID, workspaceID)

			resp := httpmock.NewStringResponse(201, fixture)
			resp.Header.Add("Content-Type", "application/json; charset=utf-8")
			return resp, nil
		},
	)

	createdTime, _ := time.Parse(time.RFC3339Nano, "2020-12-17T06:21:29.762423059Z")
	expected := &Workspace{
		ID:            "workspace-2b1cf4ac-348f-49",
		Name:          "workspace",
		Description:   "",
		Location:      "us-south",
		ResourceGroup: "Default",
		Tags:          []string{},
		Folder:        "iac/01-getting-started",
		Type:          "terraform_v0.13",
		Values:        "",
		Variables: []Variable{
			{
				Name:        "project_name",
				Value:       "gics",
				Type:        "",
				Description: "",
				Secure:      false,
			},
			{
				Name:        "environment",
				Value:       "testing",
				Type:        "",
				Description: "",
				Secure:      false,
			},
			{
				Name:        "public_key",
				Value:       "fakepublicsshkey",
				Type:        "",
				Description: "",
				Secure:      true,
			},
		},
		GitRepo: &GitRepo{
			URL:    "https://github.com/IBM/cloud-enterprise-examples",
			Branch: "master",
		},
		CreatedAt: createdTime,
		CreatedBy: "johandry@gmail.com",
		Status:    "DRAFT",
		service:   defaultService,
	}

	w := New(workspaceName, "", nil)
	act, err := w.Create()

	assert.NoError(t, err, "Creating the workspace should not fail")
	if pass := assert.NotNil(t, act, "The activity should not be nil"); pass {
		assert.NotEmpty(t, act.ID, "The activity ID should not be empty")
	}
	assert.Equal(t, expected, w, "The created Workspace is not as expected")
}

func TestWorkspace_LastActivity(t *testing.T) {
	workspaceName := "workspace"
	workspaceID := fmt.Sprintf("%s-2b1cf4ac-348f-49", workspaceName)

	httpmock.RegisterResponder("GET", `=~^https://schematics\.cloud\.ibm\.com/v1/workspaces/([\w-]+)/actions\z`,
		func(req *http.Request) (*http.Response, error) {
			// Get the fixture with the following code after getting the Token:
			// export TOKEN=$(cat .token | jq -r .access_token)
			// export WID=
			// curl -X GET "https://schematics.cloud.ibm.com/v1/workspaces/$WID/actions" -H "Authorization: Bearer $TOKEN" -H 'Accept: application/json' -H 'Content-Type: application/json'
			wID := httpmock.MustGetSubmatch(req, 1)
			assert.NotEmpty(t, wID, "Schematics Workspace should not be empty")
			assert.Equalf(t, workspaceID, wID, "The received Schematics Workspace is not %q", workspaceID)
			fixture := fmt.Sprintf(`{"workspace_name":"%s","workspace_id":"%s","actions":[{"action_id":"24481fe0f1fcd27e5642b8cc0c241eb4","name":"WORKSPACE_CREATE","status":"FAILED","message":[],"performed_by":"johandry@gmail.com","performed_at":"2020-12-17T06:21:32.573626313Z","templates":[{"template_id":"iac-f6ee24a6-1775-42","template_type":"terraform_v0.13","start_time":"2020-12-17T06:21:32.779253562Z","end_time":"2020-12-17T06:24:34.59186019Z","status":"FAILED","message":"{\"messagekey\":\"M2000_InternalError\",\"parms\":{},\"requestid\":\"de8ced31-e98f-4364-bb0d-202d6b52c864\",\"timestamp\":\"2020-12-17T06:24:34.591850876Z\"}","log_url":"https://schematics.cloud.ibm.com/v1/workspaces/%s/runtime_data/iac-f6ee24a6-1775-42/log_store/actions/24481fe0f1fcd27e5642b8cc0c241eb4","log_summary":{"activity_status":"FAILED","time_taken":181.81}}]}]}`, workspaceName, wID, wID)

			resp := httpmock.NewStringResponse(200, fixture)
			resp.Header.Add("Content-Type", "application/json; charset=utf-8")
			return resp, nil
		},
	)

	performedTime, _ := time.Parse(time.RFC3339Nano, "2020-12-17T06:21:32.573626313Z")
	startTime, _ := time.Parse(time.RFC3339Nano, "2020-12-17T06:21:32.779253562Z")
	endTime, _ := time.Parse(time.RFC3339Nano, "2020-12-17T06:24:34.59186019Z")
	expectedAct := &Activity{
		ID:          "24481fe0f1fcd27e5642b8cc0c241eb4",
		Name:        "WORKSPACE_CREATE",
		Status:      "FAILED",
		StartTime:   startTime,
		EndTime:     endTime,
		PerformedAt: performedTime,
		PerformedBy: "johandry@gmail.com",
		Message:     `{"messagekey":"M2000_InternalError","parms":{},"requestid":"de8ced31-e98f-4364-bb0d-202d6b52c864","timestamp":"2020-12-17T06:24:34.591850876Z"}`,
		WorkspaceID: workspaceID,
	}
	expectedWorkspace := &Workspace{
		ID:          workspaceID,
		Name:        workspaceName,
		Description: "",
		CreatedAt:   performedTime,
		CreatedBy:   "johandry@gmail.com",
		Status:      "DRAFT",
		service:     defaultService,
	}

	w := New(workspaceName, "", nil)
	// Main fields upddated after execute w.Create(). The other fields are not
	// used by or related to LastActivity()
	w.ID = workspaceID
	w.CreatedAt = performedTime
	w.CreatedBy = "johandry@gmail.com"
	w.Status = "DRAFT"

	act, err := w.LastActivity(activityNameForCreate)

	assert.Nil(t, err, "LastActivity() should not fail")
	assert.Equal(t, expectedAct, act)
	assert.Equal(t, expectedWorkspace, w)
}
