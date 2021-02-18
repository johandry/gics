package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"text/tabwriter"

	"github.com/johandry/gics/schematics"
)

const (
	version = "0.0.1"
)

func runSchematicsWorkspaceWithCode() *schematics.Workspace {
	w := schematics.New("GICS-Demo-with-Code", "", nil)

	w.SetOutput(os.Stderr)
	w.AddVar("prefix", "gics-demo", "", "", false)
	code := `
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
	`
	if err := w.LoadCode(code); err != nil {
		printError(err)
	}

	if err := w.Run(); err != nil {
		printError(err)
	}

	output := w.GetParam("name")
	fmt.Printf("Resource Group name: %s\n", output["name"])

	return w
}

func runSchematicsWorkspaceWithRepo() *schematics.Workspace {
	pubSSHKey, err := ioutil.ReadFile("~/.ssh/id_rsa.pub")
	if err != nil {
		printError(err)
	}

	w := schematics.New("GICS-Demo-with-Repo", "", nil)

	w.AddVar("project_name", "gics", "", "", false)
	w.AddVar("environment", "testing", "", "", false)
	w.AddVar("public_key", string(pubSSHKey), "", "", false)

	w.AddRepo("https://github.com/IBM/cloud-enterprise-examples/tree/master/iac/01-getting-started")

	if err := w.Run(); err != nil {
		printError(err)
	}

	return w
}

func printError(err error) {
	fmt.Printf("[ERROR] %s\n", err)
	os.Exit(1)
}

func printVersions() {
	fmt.Printf("> Go IBM Cloud Schematics version: %s\n", version)
	v, err := schematics.Version()
	if err != nil {
		printError(err)
	}

	fmt.Printf("> Build %s %s (SHA: %s)\n", v.Buildno, v.Builddate, v.Commitsha)
	fmt.Printf("> API Version: %s\n", v.APIVersion)
	fmt.Println("> Supported Template versions:")
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	// Header
	fmt.Fprint(tw, "Template Name")
	for _, name := range v.TemplateNames {
		fmt.Fprintf(tw, "\t%s", name)
	}
	fmt.Fprint(tw, "\n")
	// Body
	row := func(title string, r []string) {
		fmt.Fprint(tw, title)
		for _, ver := range r {
			fmt.Fprintf(tw, "\t%s", ver)
		}
		fmt.Fprint(tw, "\n")
	}
	row("Terraform", v.TerraformVersions)
	row("IBM Cloud Provider", v.IBMCloudProviderVersions)
	row("Helm", v.HelmVersions)
	row("Helm Provider", v.HelmProviderVersions)
	row("Ansible", v.AnsibleVersions)
	row("Ansible Provisioner", v.AnsibleProvisionerVersions)
	row("Kubernetes Provider", v.KubernetesProviderVersions)
	row("OC Client", v.OCClientVersions)
	row("Rest API Provider", v.RestAPIProviderVersions)

	// This code print the table in landscape: wider and shorter
	// fmt.Fprintln(tw, "Template Name\tTerraform\tIBM Cloud Provider\tHelm\tHelm Provider\tAnsible\tAnsible Provisioner\tKubernetes Provider\tOC Client\tRest API Provider")
	// for i := range v.TemplateNames {
	// 	fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
	// 		v.TemplateNames[i],
	// 		v.TerraformVersions[i],
	// 		v.IBMCloudProviderVersions[i],
	// 		v.HelmVersions[i],
	// 		v.HelmProviderVersions[i],
	// 		v.AnsibleVersions[i],
	// 		v.AnsibleProvisionerVersions[i],
	// 		v.KubernetesProviderVersions[i],
	// 		v.OCClientVersions[i],
	// 		v.RestAPIProviderVersions[i],
	// 	)
	// }
	tw.Flush()
}

func printWorkspaceList() {
	wkspList, err := schematics.List()
	if err != nil {
		printError(err)
	}
	fmt.Println("> Workspaces:")
	if len(wkspList.Workspaces) == 0 {
		fmt.Println("NONE")
		return
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(tw, "ID\tName\tDescription\tLocation\tOwner\tState\tCreated")
	for _, w := range wkspList.Workspaces {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", w.ID, w.Name, w.Description, w.Location, w.Owner, w.State, w.Created)
	}
	tw.Flush()
}

func main() {
	if icAPIKey := os.Getenv("IC_API_KEY"); len(icAPIKey) == 0 {
		printError(fmt.Errorf("[ERROR] GICS requires the IBM Cloud API Key exported in the 'IC_API_KEY' variable"))
	}

	// printVersions()
	// printWorkspaceList()

	runSchematicsWorkspaceWithCode()
	// if err := w.Delete(true); err != nil {
	// 	printError(err)
	// }

	// w = runSchematicsWorkspaceWithRepo()
	// if err := w.Delete(true); err != nil {
	// 	printError(err)
	// }
}
