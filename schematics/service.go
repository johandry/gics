package schematics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/IBM/go-sdk-core/core"
	apiv1 "github.com/johandry/gics/schematics/api/v1"
)

// Run oapi-codegen to regenerate the schematics boilerplate version 1.0
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=v1 --generate types,client -o ./api/v1/client.gen.go ./api/v1/schematics.json

const (
	defaultAPIEndpoint = "https://schematics.cloud.ibm.com"
	defaultAPIVersion  = "v1"
	userAgent          = "GICS"
	timeout            = 20 // Seconds, longest timeout. Set shorter timeouts with context
)

// Service is the Schematics service
type Service struct {
	client              *apiv1.Client
	clientWithResponses *apiv1.ClientWithResponses
	apiVersion          string
	authenticator       *core.IamAuthenticator
}

// ServiceOptions are the parameters to pass to create a new Schematics Service
type ServiceOptions struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
	APIVersion string
}

var defaultService = NewService(nil)

// NewService creates a new Schematics service to communicate with the Schematics API endpoint
func NewService(opt *ServiceOptions) *Service {
	// Set default values to missing input parameters
	if opt == nil {
		opt = &ServiceOptions{}
	}
	if len(opt.BaseURL) == 0 {
		opt.BaseURL = defaultAPIEndpoint
	}
	if opt.HTTPClient == nil {
		schHTTPpClient := &http.Client{
			Timeout: timeout * time.Second,
		}
		opt.HTTPClient = schHTTPpClient
	}

	// Get API Key
	if len(opt.APIKey) == 0 {
		opt.APIKey = os.Getenv("IC_API_KEY")
	}

	if len(opt.APIVersion) == 0 {
		opt.APIVersion = defaultAPIVersion
	}

	authenticator := &core.IamAuthenticator{
		ApiKey: opt.APIKey,
	}
	icc := &ICClient{
		UserAgent:     userAgent,
		http:          opt.HTTPClient,
		authenticator: authenticator,
	}

	c, _ := apiv1.NewClient(opt.BaseURL, apiv1.WithHTTPClient(icc))
	cwr, _ := apiv1.NewClientWithResponses(opt.BaseURL, apiv1.WithHTTPClient(icc))

	return &Service{
		client:              c,
		clientWithResponses: cwr,
		apiVersion:          "/" + opt.APIVersion,
		authenticator:       authenticator,
	}
}

// ICClient is an HTTP Client wrapped by the Schematics client to communicate
// with the IBM Cloud endpoint API and provide the provide the authentication
type ICClient struct {
	UserAgent     string
	http          *http.Client
	authenticator *core.IamAuthenticator
}

// Do implements the Do method so ICClient is a HttpRequestDoer interface
func (c *ICClient) Do(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	if err := c.authenticator.Authenticate(req); err != nil {
		return nil, err
	}

	return c.http.Do(req)
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// APIError encapsulate a returned error from the API
type APIError struct {
	Message string                  `json:"message,omitempty"`
	Error   apiv1.APIMessagePayload `json:"error,omitempty"`
}

func getAPIError(msg string, body []byte) error {
	var msgPayload apiv1.APIMessagePayload
	if err := json.Unmarshal(body, &msgPayload); err != nil {
		return err
	}

	apiErr := APIError{
		Message: msg,
		Error:   msgPayload,
	}

	jsonAPIError, err := json.Marshal(apiErr)
	if err != nil {
		return err
	}

	return fmt.Errorf("%s", jsonAPIError)
}

func (w *Workspace) getToken() (string, error) {
	req, _ := http.NewRequest("GET", "", nil)
	if err := w.service.authenticator.Authenticate(req); err != nil {
		return "", err
	}
	token := req.Header.Get("Authorization")

	return token, nil
}
