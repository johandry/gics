package schematics

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	httpmock.RegisterResponder("GET", "https://schematics.cloud.ibm.com/v1/version",
		func(req *http.Request) (*http.Response, error) {
			// Get the fixture with the following code after getting the Token:
			// curl -X GET https://schematics.cloud.ibm.com/v1/version
			fixture := `{"commitsha":"40db47f4aae3b84f23b173e5ab9c3efe71092d37","builddate":"2020-12-11T10:33:56Z","buildno":"6525","terraform_version":"v0.11.14","terraform_provider_version":"v0.31.0","helm_version":"v2.14.2","helm_provider_version":"v0.10.0","supported_template_types":[{"ansible":"v2.9.7","ansible_provisioner":"v2.3.3","helm":"v2.14.2","helm_provider":"v0.10.4a","ibm_cloud_provider":"v0.31.0","kubernetes_provider":"v1.10.0a","oc_client":"v3.11.0","provider_restapi":"v1.10.0","template_name":"terraform_v0.11","terraform":"v0.11.14"},{"ansible":"v2.9.7","ansible_provisioner":"v2.3.3","helm":"v3.1.1","helm_provider":"v0.10.4a","ibm_cloud_provider":"v1.17.0","kubernetes_provider":"v1.10.0a","oc_client":"v3.11.0","provider_restapi":"v1.10.0","template_name":"terraform_v0.12","terraform":"v0.12.20"},{"ansible":"v2.9.7","ansible_provisioner":"v2.3.3","helm":"v3.1.1","helm_provider":"v0.10.4a","ibm_cloud_provider":"v1.13.1","kubernetes_provider":"v1.10.0a","oc_client":"v3.11.0","provider_restapi":"v1.10.0","template_name":"terraform_v0.13","terraform":"v0.13.5"}]}`

			resp := httpmock.NewStringResponse(200, fixture)
			resp.Header.Add("Content-Type", "application/json; charset=utf-8")
			return resp, nil
		},
	)

	expected := &SupportedVersion{
		Commitsha:                  "40db47f4aae3b84f23b173e5ab9c3efe71092d37",
		Builddate:                  "2020-12-11T10:33:56Z",
		Buildno:                    "6525",
		APIVersion:                 "1.0",
		TerraformVersions:          []string{"v0.11.14", "v0.12.20", "v0.13.5"},
		IBMCloudProviderVersions:   []string{"v0.31.0", "v1.17.0", "v1.13.1"},
		HelmVersions:               []string{"v2.14.2", "v3.1.1", "v3.1.1"},
		HelmProviderVersions:       []string{"v0.10.4a", "v0.10.4a", "v0.10.4a"},
		AnsibleVersions:            []string{"v2.9.7", "v2.9.7", "v2.9.7"},
		AnsibleProvisionerVersions: []string{"v2.3.3", "v2.3.3", "v2.3.3"},
		KubernetesProviderVersions: []string{"v1.10.0a", "v1.10.0a", "v1.10.0a"},
		OCClientVersions:           []string{"v3.11.0", "v3.11.0", "v3.11.0"},
		RestAPIProviderVersions:    []string{"v1.10.0", "v1.10.0", "v1.10.0"},
		TemplateNames:              []string{"terraform_v0.11", "terraform_v0.12", "terraform_v0.13"},
	}

	resp, err := Version()

	assert.Nil(t, err, "Version should not return an error")
	assert.Equal(t, expected, resp)
}
