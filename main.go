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

func main() {
	fmt.Printf("> Go IBM Cloud Schematics version: %s\n", version)

	wkspList, err := schematics.List()
	if err != nil {
		fmt.Printf("[ERROR] %s", err)
		return
	}
	fmt.Println("> Workspaces:")
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(tw, "ID\tName\tDescription\tLocation\tOwner\tState\tCreated")
	for _, w := range wkspList.Workspaces {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s", w.ID, w.Name, w.Description, w.Location, w.Owner, w.State, w.Created)
	}
	tw.Flush()
}
