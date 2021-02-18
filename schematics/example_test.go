package schematics_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/johandry/gics/schematics"
)

// ExampleWorkspace_Run creates and runs a workspace with a simple/short code to
// provision a Resource Group. The executing is asynchronous, it does not block
// the execution of the code.
func ExampleWorkspace_Run() {
	if icAPIKey := os.Getenv("IC_API_KEY"); len(icAPIKey) == 0 {
		fmt.Printf("[ERROR] GICS requires the IBM Cloud API Key exported in the 'IC_API_KEY' variable")
		return
	}

	w := schematics.New("GICS Demo", "", nil)
	defer w.Delete(true)

	w.AddVar("prefix", "gics-demo", "", "", false)
	w.Code = []byte(`
		variable "prefix" {}
		provider "ibm" {
			generation         = 2
			region             = "us-south"
		}
		resource "ibm_resource_group" "group" {
			name = "${var.prefix}-group"
		}
		output "name" {
			value = ibm_resource_group.group.name
		}
	`)

	if err := w.Run(); err != nil {
		fmt.Printf("[ERROR] Fail the execution of the Schematics Workspace. %s", err)
		return
	}

	output := w.GetParam("name")
	fmt.Printf("Resource Group name: %s", output["name"])
	// Output: Resource Group name: gics-demo-group
}

// ExampleWorkspace_RunSync creates and runs a workspace from a code located in
// a private GitHub repository to provision a Resource Group
func ExampleWorkspace_Run_private_repo() {
	if icAPIKey := os.Getenv("IC_API_KEY"); len(icAPIKey) == 0 {
		fmt.Printf("[ERROR] GICS requires the IBM Cloud API Key exported in the 'IC_API_KEY' variable")
		return
	}

	w := schematics.New("GICS Demo with a GH repository", "", nil)
	defer w.Delete(true)

	w.AddVar("prefix", "gics-gh-priv-demo", "", "", false)
	w.AddVar("enable", "true", "bool", "", false)
	w.GitRepo = &schematics.GitRepo{
		URL:    "https://github.com/johandry/gics-priv-test",
		Branch: "gics_demo",
		Token:  os.Getenv("GICS_GH_TOKEN"),
	}

	if err := w.Run(); err != nil {
		fmt.Printf("[ERROR] Fail the execution of the Schematics Workspace. %s", err)
		return
	}

	output := w.GetParam("name")
	fmt.Printf("Resource Group name: %s", output["name"])
	// Output: Resource Group name: gics-gh-priv-demo-group
}

// ExampleApply creates a workspace and execute all the actions to provision a
// Resource Group from a public GitHub repository
func ExampleWorkspace_Apply() {
	if icAPIKey := os.Getenv("IC_API_KEY"); len(icAPIKey) == 0 {
		fmt.Printf("[ERROR] GICS requires the IBM Cloud API Key exported in the 'IC_API_KEY' variable")
		return
	}

	w := schematics.New("GICS Demo with a GH repository", "", nil)
	defer w.Delete(true)

	w.AddVar("prefix", "gics-gh-demo", "", "", false)
	w.AddVar("enable", "true", "bool", "", false)
	w.GitRepo = &schematics.GitRepo{
		URL: "https://github.com/johandry/gics-pub-test",
	}

	act, err := w.Create()
	if err != nil {
		fmt.Printf("[ERROR] Fail to trigger the Schematics Workspace creation. %s", err)
		return
	}
	if err := act.Wait(); err != nil {
		fmt.Printf("[ERROR] Fail to create the Schematics Workspace. %s", err)
		return
	}

	act, err = w.Plan()
	if err != nil {
		fmt.Printf("[ERROR] Fail trigger to generate the plan of the Schematics Workspace. %s", err)
		return
	}
	if err := act.SetOutput(os.Stdout).Wait(); err != nil {
		fmt.Printf("[ERROR] Fail to generate the plan of the Schematics Workspace. %s", err)
		return
	}

	var buf bytes.Buffer
	act, err = w.Apply()
	if err != nil {
		fmt.Printf("[ERROR] Fail to trigger the Schematics Workspace plan execution. %s", err)
		return
	}
	go func() {
		// do something with the buffer or the activity
	}()
	if err := act.SetOutput(&buf).Wait(); err != nil {
		fmt.Printf("[ERROR] Fail to execute the plan of the Schematics Workspace. %s", err)
		return
	}
	fmt.Println(buf.String())

	output := w.GetParam("name")
	fmt.Printf("Resource Group name: %s", output["name"])
	// Output: Resource Group name: gics-gh-demo-group
}
