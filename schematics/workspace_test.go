package schematics

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkspace_NewWorkspace(t *testing.T) {
	todayWorkspace := fmt.Sprintf("workspace_%s", time.Now().Format("MM_DD_YYYY"))
	type argsWorkspace struct {
		name        string
		description string
	}
	type argsVariable struct {
		name        string
		value       string
		varType     string
		description string
		secure      bool
	}
	type argsGitRepo struct {
		URL string
	}
	tests := []struct {
		name          string
		argsWorkspace argsWorkspace
		wantWorkspace *Workspace
		argsVariables []argsVariable
		wantVariables []Variable
		wantErr       []bool
		argsGitRepo   *argsGitRepo
		wantRepo      *GitRepo
	}{
		{
			"no name",
			argsWorkspace{"", ""},
			&Workspace{
				Name:        todayWorkspace,
				Description: "",
				service:     defaultService,
				Status:      "NEW",
			},
			nil, nil, nil, nil, nil,
		},
		{
			"no name in variable",
			argsWorkspace{"workspace", ""},
			&Workspace{
				Name:        "workspace",
				Description: "",
				service:     defaultService,
				Status:      "NEW",
			},
			[]argsVariable{
				{"", "", "", "", false},
			},
			nil,
			[]bool{true},
			nil, nil,
		},
		{
			"default type variable",
			argsWorkspace{"workspace", ""},
			&Workspace{
				Name:        "workspace",
				Description: "",
				service:     defaultService,
				Status:      "NEW",
			},
			[]argsVariable{
				{"a1", "b1", "string", "", false},
				{"a2", "b2", "", "", false},
				{"a3", "b3", "string", "", true},
				{"a4", "b4", "int", "faulty number, but I don't care", false},
			},
			[]Variable{
				{"a1", "b1", "string", "", false, false},
				{"a2", "b2", "string", "", false, false},
				{"a3", "b3", "string", "", true, false},
				{"a4", "b4", "int", "faulty number, but I don't care", false, false},
			},
			[]bool{false, false, false, false},
			nil, nil,
		},
		{
			"repeate variable name",
			argsWorkspace{"workspace", ""},
			&Workspace{
				Name:        "workspace",
				Description: "",
				service:     defaultService,
				Status:      "NEW",
			},
			[]argsVariable{
				{"a1", "b1", "string", "", false},
				{"a1", "b2", "", "", false},
				{"a2", "b3", "string", "", true},
				{"a2", "b4", "int", "faulty number, but I don't care", false},
			},
			[]Variable{
				{"a1", "b1", "string", "", false, false},
				{"a2", "b3", "string", "", true, false},
			},
			[]bool{false, true, false, true},
			nil, nil,
		},
		{
			"git url",
			argsWorkspace{"git_url", ""},
			&Workspace{
				Name:        "git_url",
				Description: "",
				service:     defaultService,
				Status:      "NEW",
			},
			nil, nil, nil,
			&argsGitRepo{
				URL: "https://github.com/IBM/cloud-enterprise-examples/tree/master/iac/01-getting-started",
			},
			&GitRepo{
				URL: "https://github.com/IBM/cloud-enterprise-examples/tree/master/iac/01-getting-started",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWorkspace := New(tt.argsWorkspace.name, tt.argsWorkspace.description, nil)
			assert.Equal(t, tt.wantWorkspace, gotWorkspace, "Workspaces should be equal")

			for i, argsVar := range tt.argsVariables {
				err := gotWorkspace.AddVar(argsVar.name, argsVar.value, argsVar.varType, argsVar.description, argsVar.secure)
				assert.Equalf(t, (err != nil), tt.wantErr[i], "Workspace.AddVar() error = %v, wantErr %v", err, tt.wantErr[i])
			}
			assert.Equal(t, tt.wantVariables, gotWorkspace.Variables)

			if tt.argsGitRepo != nil {
				gotWorkspace.AddRepo(tt.argsGitRepo.URL)
			}
			assert.Equal(t, tt.wantRepo, gotWorkspace.GitRepo)
		})
	}
}
