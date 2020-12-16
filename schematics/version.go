package schematics

import (
	"context"
	"fmt"
	"time"

	apiv1 "github.com/johandry/gics/schematics/api/v1"
)

// SupportedVersion contains the version of the Schamtics, used components and supported
// providers
type SupportedVersion struct {
	Builddate                  string   `json:"builddate,omitempty" yaml:"builddate,omitempty"`
	Buildno                    string   `json:"buildno,omitempty" yaml:"buildno,omitempty"`
	Commitsha                  string   `json:"commitsha,omitempty" yaml:"commitsha,omitempty"`
	APIVersion                 string   `json:"api_version,omitempty" yaml:"api_version,omitempty"`
	TerraformVersions          []string `json:"terraform_versions,omitempty" yaml:"terraform_versions,omitempty"`
	IBMCloudProviderVersions   []string `json:"ibm_cloud_provider_versions,omitempty" yaml:"ibm_cloud_provider_versions,omitempty"`
	HelmVersions               []string `json:"helm_versions,omitempty" yaml:"helm_version,omitempty"`
	HelmProviderVersions       []string `json:"helm_provider_versions,omitempty" yaml:"helm_provider_version,omitempty"`
	AnsibleVersions            []string `json:"ansible_versions,omitempty" yaml:"ansible_versions,omitempty"`
	AnsibleProvisionerVersions []string `json:"ansible_provisioner_versions,omitempty" yaml:"ansible_provisioner_versions,omitempty"`
	KubernetesProviderVersions []string `json:"kubernetes_provider_versions,omitempty" yaml:"kubernetes_provider_versions,omitempty"`
	OCClientVersions           []string `json:"oc_client_versions,omitempty" yaml:"oc_client_versions,omitempty"`
	RestAPIProviderVersions    []string `json:"rest_api_provider_versions,omitempty" yaml:"rest_api_provider_versions,omitempty"`
	TemplateNames              []string `json:"template_names,omitempty" yaml:"template_names,omitempty"`
}

// Version returns the Schematics versions and supported components using the
// default Schematics service
func Version() (*SupportedVersion, error) {
	ctx := context.Background()
	return defaultService.Version(ctx)
}

// Version returns the Schematics versions and supported components
func (s *Service) Version(ctx context.Context) (*SupportedVersion, error) {
	// Version Timeout
	ctx, cancelFunc := context.WithTimeout(ctx, listTimeout*time.Second)
	defer cancelFunc()

	resp, err := s.clientWithResponses.GetSchematicsVersionWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode(); code != 200 {
		return nil, fmt.Errorf(`{"status_code": %d, "status": %q}`, code, resp.Status())
	}
	response := resp.JSON200

	if response == nil {
		return &SupportedVersion{}, nil
	}

	supportedVersions := map[string][]string{
		"ansible":             []string{},
		"ansible_provisioner": []string{},
		"helm":                []string{},
		"helm_provider":       []string{},
		"ibm_cloud_provider":  []string{},
		"kubernetes_provider": []string{},
		"oc_client":           []string{},
		"provider_restapi":    []string{},
		"template_name":       []string{},
		"terraform":           []string{},
	}
	for _, template := range *response.SupportedTemplateTypes {
		for key, iface := range template {
			version := fmt.Sprintf("%v", iface)
			supportedVersions[key] = append(supportedVersions[key], version)
		}
	}

	v := &SupportedVersion{
		Builddate:                  *response.Builddate,
		Buildno:                    *response.Buildno,
		Commitsha:                  *response.Commitsha,
		APIVersion:                 apiv1.APIVersion,
		TerraformVersions:          supportedVersions["terraform"],
		IBMCloudProviderVersions:   supportedVersions["ibm_cloud_provider"],
		HelmVersions:               supportedVersions["helm"],
		HelmProviderVersions:       supportedVersions["helm_provider"],
		AnsibleVersions:            supportedVersions["ansible"],
		AnsibleProvisionerVersions: supportedVersions["ansible_provisioner"],
		KubernetesProviderVersions: supportedVersions["kubernetes_provider"],
		OCClientVersions:           supportedVersions["oc_client"],
		RestAPIProviderVersions:    supportedVersions["provider_restapi"],
		TemplateNames:              supportedVersions["template_name"],
	}

	return v, nil
}
