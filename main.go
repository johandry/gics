package main

import (
	"fmt"
	"os"

	"text/tabwriter"

	"github.com/johandry/gics/schematics"
)

const (
	version = "0.0.1"
)

// func runSchematicsWorkspace() (*schematics.Workspace, error) {
// 	w := schematics.New("GICS Demo", "")

// 	w.AddVar("prefix", "gics-demo", "", "", false)
// 	w.Code = []byte(`
//     variable "prefix" {}
//     provider "ibm" {
//       generation         = 2
//       region             = "us-south"
//     }
//     resource "ibm_resource_group" "group" {
//       name = "${var.prefix}-group"
//     }
//     output "name" {
//       value = ibm_resource_group.group.name
//     }
//   `)

// 	if err := w.Run(); err != nil {
// 		return w, fmt.Errorf("Fail the execution of the Schematics Workspace. %s", err)
// 	}

// 	return w, nil
// }

func printError(err error) {
	fmt.Printf("[ERROR] %s", err)
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
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(tw, "ID\tName\tDescription\tLocation\tOwner\tState\tCreated")
	for _, w := range wkspList.Workspaces {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", w.ID, w.Name, w.Description, w.Location, w.Owner, w.State, w.Created)
	}
	tw.Flush()
}

func main() {
	printVersions()
	// printWorkspaceList()

	// w, err := runSchematicsWorkspace()
	// if err != nil {
	// 	fmt.Printf("[ERROR] %s", err)
	// 	return
	// }

	// output, err := w.Output()
	// if err != nil {
	// 	fmt.Printf("[ERROR] %s", err)
	// 	return
	// }
	// fmt.Printf("Resource Group name: %s", output["name"])

	// if err := w.Delete(true); err != nil {
	// 	fmt.Println(err)
	// }
}
