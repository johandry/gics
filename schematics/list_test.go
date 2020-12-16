package schematics

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Get the fixture with the following code after getting the API Key:
	// export IC_API_KEY='xxxxxxxxxxxxxxxxxxxxxx'
	// curl -s -k -X POST --header "Content-Type: application/x-www-form-urlencoded"  --data-urlencode "grant_type=urn:ibm:params:oauth:grant-type:apikey"  --data-urlencode "apikey=$IC_API_KEY" "https://iam.cloud.ibm.com/identity/token"
	fixture := `{"access_token":"eyJraWQiOiIyMDIwMTEyMTE4MzQiLCJhbGciOiJSUzI1NiJ9.eyJpYW1faWQiOiJJQk1pZC01NTAwMDVUQ1VEIiwiaWQiOiJJQk1pZC01NTAwMDVUQ1VEIiwicmVhbG1pZCI6IklCTWlkIiwianRpIjoiNmVjOGRjMjQtN2Q1ZS00NTc4LTk5ZDMtNTNkOGUxY2Y0ZGUzIiwiaWRlbnRpZmllciI6IjU1MDAwNVRDVUQiLCJnaXZlbl9uYW1lIjoiSm9oYW5kcnkiLCJmYW1pbHlfbmFtZSI6IkFtYWRvciIsIm5hbWUiOiJKb2hhbmRyeSBBbWFkb3IiLCJlbWFpbCI6ImpvaGFuZHJ5QGdtYWlsLmNvbSIsInN1YiI6ImpvaGFuZHJ5QGdtYWlsLmNvbSIsImFjY291bnQiOnsidmFsaWQiOnRydWUsImJzcyI6ImYwNjkyNzlhNzc4YTQ4MTc4ZmMzNzllOTk3OTdmODg3IiwiZnJvemVuIjp0cnVlfSwiaWF0IjoxNjA3OTI1Njk1LCJleHAiOjE2MDc5MjkyOTUsImlzcyI6Imh0dHBzOi8vaWFtLmNsb3VkLmlibS5jb20vaWRlbnRpdHkiLCJncmFudF90eXBlIjoidXJuOmlibTpwYXJhbXM6b2F1dGg6Z3JhbnQtdHlwZTphcGlrZXkiLCJzY29wZSI6ImlibSBvcGVuaWQiLCJjbGllbnRfaWQiOiJkZWZhdWx0IiwiYWNyIjoxLCJhbXIiOlsicHdkIl19.I9c-LWVK5j1RkE9wmKTkV6M-d27nd7o6fYHhxePlTo-RqW3_D5g1DI_9O4rqCGe9O0M_CymDyRuEK2t7x7qCRvsp5G9GCPWbroUjIVrkZNVYxIgEO7aQ_nYXhmyr2lr5WDS6l5kZOEwFXozPmYKVICn3jcMPlgGF0D6OavFYabQ6YrX9wFwuAgwW8KCZrIdXAbGU2HBVcKT-uByElTOxvjsz0ctgmVSKTdphfm-3bKJvBxno1eXjHyShznV89MsnbTsld_lqzJ-452bHnoLgWilS-sL3Wm_hdlyB37GTSHWKI_9aTTR-LNWWE6SRsDgstOxEHlZSkMhjqJnF94MPfA","refresh_token":"OKDKZxG56yTYw1az7q7EExwlbHeDbrfwdvAzQIDtz7oDxv2TNH3QEH-1L8e1601eE4fIEPL77A0GMDeAc5F1lyslAW0U4PWiZSMyq0tpBoS_v3iTvkCElj6TStivwScsF5znXKYfHRU_sZ7-P6hX2_Z3kMvdrpU410OR9aAe6a53YCAK4StQlSdqtc5JpVd0o3ohDtt7IAzqmGB9xMFDfVCbNbF08IY22PWOZzWVv6G8yqYrZw-dC0enNiJ35fQsljNd_vp3cGIYjNDUgh1THi2HFnlttngPW2kcE7lBE7fq3Z7azOpbFD2HT7jEl_YuudW2DYzgrrGDiujQ8_rYL_eLuDNRR8twkqnFYKV3oOIWmE0s10hgnWBZbFpTw4CQzPqt3wKJ7nKvp8-HuNMclKW0ltX-BO0VUKvr5j83zXDxdeE75OE1eTasg7RLKgkEgV1BvuOCD4CtLIzWycnPYv0haOJbK1n--YsDhFOv3GuP3XN0Mi4lJ1esPpjQxHrep_E2Z3wDSd6qJfUK33iatWBfSVnnq5qxbqShGG6IVg1Ki9L71995aEbc5fYmVZuMpZxn6Qttljp97gmAae8_pGc84rzSObB6cbQC3SJv-XCBCmj0n0RQES9aG_eHd5q143o6_FB0_eEtT-4vLwrmE_EgLjA-S_3QzHgxmQC9KZasau9F44bfuI8ZS1oVSerT7ElYrVfg4kBRUT-p1noX6nrjiU7jlAmzPy6ETrRHPqAQ3aEs5Mjd7aCuzZJJYqy5e-kTquvlQJWFTFydSoqciuWOkLtRtTDmX_YBVqqWWS6MUAUGbymvFOx7ItMz6cAvLwkZM6f5S5qKCcNlQdq1-G7vFTDGKwL_ncnkecCpthffxJMIWp-qQ_arLcMICd0wVHn_ZerdmjdEuHCPIYitKDECdx5EE_0pqQgpb9eJYTt5Xu61Byr9Zge_10OAZiKxs6esjorrktqAADg0MZlyzQUYyi95zYiQeYXQ0Y4669TyxyjanMfI1kOaMOtN4x_J-tpPhCU_0P3cwZ2LOBRrFaFfwVleuTkMNcnz43PI6Nms9dB_MtY8hNAMIfMOdyIMahfBjZWGFrHsgA6UIJvNatDnjhfiFKmgjWkiho6Qa0Lpvg","token_type":"Bearer","expires_in":3600,"expiration":1607929295,"refresh_token_expiration":1610517695,"scope":"ibm openid"}`
	httpmock.RegisterResponder("POST", "https://iam.cloud.ibm.com/identity/token",
		httpmock.NewStringResponder(200, fixture))

	os.Exit(m.Run())
}

func TestList(t *testing.T) {
	httpmock.RegisterResponder("GET", "https://schematics.cloud.ibm.com/v1/workspaces",
		func(req *http.Request) (*http.Response, error) {
			// Get the fixture with the following code after getting the Token:
			// export TOKEN=$(cat .token | jq -r .access_token)
			// curl -s -X GET "https://schematics.cloud.ibm.com/v1/workspaces" -H "Authorization: Bearer $TOKEN" -H 'Accept: application/json' -H 'Content-Type: application/json'
			fixture := `{"offset":0,"limit":100,"count":1,"workspaces":[{"id":"myworkspace-270d8f90-b67a-46","name":"myworkspace","crn":"crn:v1:bluemix:public:schematics:us-south:a/f069279a778a48178fc379e99797f887:ba49546a-1075-4e1d-bfde-9d5ef3ca031e:workspace:myworkspace-270d8f90-b67a-46","type":["terraform_v0.12"],"description":"","resource_group":"Default","location":"us-east","tags":[],"created_at":"2020-11-01T18:30:15.364093434Z","created_by":"johandry@gmail.com","status":"FAILED","workspace_status_msg":{"status_code":"500","status_msg":"rpc error: code = ResourceExhausted desc = trying to send message larger than max (129563174 vs. 83886080)"},"workspace_status":{"frozen":false,"locked":false},"template_repo":{"url":"https://github.com/IBM/cloud-enterprise-examples","branch":"master","full_url":"https://github.com/IBM/cloud-enterprise-examples/tree/master/iac/02-schematics","has_uploadedgitrepotar":false},"template_data":[{"id":"e86f76d0-b77f-49","folder":"iac/02-schematics","type":"terraform_v0.12","values_url":"https://schematics.cloud.ibm.com/v1/workspaces/myworkspace-270d8f90-b67a-46/template_data/e86f76d0-b77f-49/values","values":"","values_metadata":[{"description":"","name":"project_name","type":"string"},{"description":"","name":"environment","type":"string"},{"description":"","name":"public_key","type":"string"},{"default":"8080","description":"","name":"port","type":"string"}],"has_githubtoken":false}],"runtime_data":[{"id":"e86f76d0-b77f-49","engine_name":"terraform","engine_version":"v0.12.20","state_store_url":"https://schematics.cloud.ibm.com/v1/workspaces/myworkspace-270d8f90-b67a-46/runtime_data/e86f76d0-b77f-49/state_store","log_store_url":"https://schematics.cloud.ibm.com/v1/workspaces/myworkspace-270d8f90-b67a-46/runtime_data/e86f76d0-b77f-49/log_store"}],"shared_data":{"resource_group_id":""},"applied_shareddata_ids":null,"updated_by":"johandry@gmail.com","updated_at":"2020-11-01T18:33:05.807531893Z","last_health_check_at":"0001-01-01T00:00:00Z"}]}`

			resp := httpmock.NewStringResponse(200, fixture)
			resp.Header.Add("Content-Type", "application/json; charset=utf-8")
			return resp, nil
		},
	)

	createdTime, _ := time.Parse(time.RFC3339Nano, "2020-11-01T18:30:15.364093434Z")
	expected := &WorkspaceList{
		Workspaces: []WorkspaceSummary{
			WorkspaceSummary{
				ID:          "myworkspace-270d8f90-b67a-46",
				Name:        "myworkspace",
				Description: "",
				Location:    "us-east",
				Owner:       "johandry@gmail.com",
				State:       "FAILED",
				Created:     createdTime,
			},
		},
	}

	resp, err := List()

	assert.Nil(t, err, "List of workspaces should not return an error")
	assert.Equal(t, expected, resp)
}
