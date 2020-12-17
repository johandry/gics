# Go IBM Cloud Schematics

**This is still a Work in Progress**, all the code documented here is in development and testing.

Go IBM Cloud Schematics (**GICS**) is a Go module and CLI to work with IBM Cloud Schematics. IBM Cloud Schematics delivers Terraform-as-a-Service using the open source Terraform provisioning engine. With IBM Cloud Schematics, you can organize your cloud resources across environments by using **workspaces**. Every workspace points to a set of Terraform configuration files and input parameters.

- [Go IBM Cloud Schematics](#go-ibm-cloud-schematics)
  - [Requirements](#requirements)
  - [How to use the GICS as a Go module](#how-to-use-the-gics-as-a-go-module)
  - [How to use the GICS CLI](#how-to-use-the-gics-cli)

## Requirements

To use GICS as a Go module or CLI you need an IBM Cloud account (any kind, including a free account) and the IBM Cloud API Key.

If you don't have an IBM Cloud account you may [create an IBM Cloud account](https://cloud.ibm.com/registration), it's free for development and do not require a credit card.

Follow these instructions to setup the IBM Cloud API Key, for more information read [Creating an API key](https://cloud.ibm.com/docs/account?topic=account-userapikey#create_user_key).

1. In the [IBM Cloud console](https://cloud.ibm.com), go to **Manage** > **Access (IAM)** > **[API keys](https://cloud.ibm.com/iam/apikeys)**.
2. Click **Create an IBM Cloud API key**.
3. Enter a *Name* and *Description* for your API key. Then, click **Create**.
4. Click on the *eye* icon to show the API key, or click **Copy**, or click **Download** to save it as a JSON file.

In a Terminal, if you have the IBM Cloud CLI (`ibmcloud`), execute:

```bash
# Login to IBM Cloud and target an account
ibmcloud login --sso
# Generate and download a new API Key
ibmcloud iam api-key-create TerraformKey -d "API Key for GICS $(date +%m-%d-%Y)" --file ~/.ibm_credentials.json
# Export the Key using the environment variable IC_API_KEY
export IC_API_KEY=$(grep '"apikey":' ~/.ibm_credentials.json | sed 's/.*: "\(.*\)".*/\1/')
# Or using jq:
export IC_API_KEY=$(jq -r .apikey ~/.ibm_credentials.json)
```

It's recommended to export the environment variable `IC_API_KEY` with the API Key in the profile file (i.e. `~/.bashrc` or `~/.zshrc`).

## How to use the GICS as a Go module

After import the Go package in your code you can use the methods `New()`, `Create()`, `Plan()`, `Apply()`, `Destroy()` and `Delete()` to execute the same actions on the Schematics Workspace in an async or non-blocking way. You may also use the method `Wait()` from the returned activity to wait for an action to be completed.

The method `Run()` can be used to create, plan and apply/execute the given Terraform code in a synchronous way, blocking the execution of the code until the entire process successfully finish or fail.

The following example creates and execute a Schematics Workspace to create a Resource Group in your account.

```go
package main

import (
  "fmt"

  "github.com/johandry/gics/schematics"
)

func run() (*schematics.Workspace, string, error) {
  w := schematics.New("GICS Demo", "")

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
    return w, "", fmt.Errorf("[ERROR] Fail the execution of the Schematics Workspace. %s", err)
  }

  output, err := w.Output()
  if err != nil {
    return w, "", fmt.Errorf("[ERROR] Fail getting the output parameters of the Schematics Workspace. %s", err)
  }

  return w, output["name"], nil
}

func main() {
  w, name, err := run()
  if err != nil {
    fmt.Println(err)
    return
  }

  fmt.Printf("Resource Group name: %s", name)

  if err := w.Delete(true); err != nil {
    fmt.Println(err)
  }
}
```

## How to use the GICS CLI

You can use `ibmcloud` with the `schematics` plugin to handle Schematics however it requires multiple calls, one per action to execute (new, plan and apply). With `gics` there is only call to the command providing all the input parameters to create and apply the code. With IBM Cloud Schematics the Terraform code is in a GitHub repo, this can be done with `gics` but also you can provide a local directory or a single file.

Read the following example to know how to use `gics`.

Create, plan and apply a workspace with a terraform code to create a Resource Group:

```bash
gics run \
  --name "GICS Demo" \
  --var "prefix=gics-demo" \
  --file - <<EOC
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
EOC
```

Or, having a directory with multiple Terraform files and a JSON file with the input parameters:

```bash
gics run \
  --name "GICS Demo with multiple files" \
  --var-file parameters.json \
  --tf-dir terraform_code
```

You can still use a code that is in a repository, in an specific folder and on some branch:

```bash
gics run \
  --name "GICS Demo with a GH repository" \
  --var "bool enable=true" \
  --var-file parameters.json \
  --url "https://github.ibm.com/pett/cloud-pak-sandboxes" \
  --branch "gics_demo" \
  --gh-token some-github_token-here
```

You can still use `gics` like you use the `ibmcloud schematics` plugin, executing one action at a time, like so:

```bash
> gics create \
  --name "GICS Demo with a GH repository" \
  --var "bool enable=true" \
  --var-file parameters.json \
  --url "https://github.com/johandry/iks" \
  --output json > gics.json

> cat gics.json
{"workspace_id": "270d8f90-fake-46"}

> gics plan --workspace-id $(jq -r .workspace_id gics.json) --show-log
 2020/11/01 18:33:06 -----  New Workspace Action  -----
 2020/11/01 18:33:06 Request: activityId=abbd9fake8368ca05519f8fake8d1, account=f06..........................887, owner=user@email.com, requestID=a5b24441-fake-fake-fake-e775c8325d53
      ...
 2020/11/01 18:36:41 Done with the Activity

> gics apply --workspace-id $(jq -r .workspace_id gics.json) --output json
{"activity_id": "abbd9f1fake68ca05519f8c260d8d1", "url": "https://cloud.ibm.com/schematics/workspaces/270d8f90-fake-46/log/abbd9f1fake68ca05519f8c260d8d1"}

> gics show-log --workspace-id $(jq -r .workspace_id gics.json) --activity-id abbd9f1fake68ca05519f8c260d8d1
 2020/11/01 18:33:06 -----  New Workspace Action  -----
 2020/11/01 18:33:06 Request: activityId=abbd9f1fake68ca05519f8c260d8d1, account=f06..........................887, owner=user@email.com, requestID=a5b24441-fake-fake-fake-e5.fake.d53
      ...
 2020/11/01 18:36:41 Done with the Activity

> gics destroy --workspace-id $(jq -r .workspace_id gics.json)
Destroy log url: https://cloud.ibm.com/schematics/workspaces/myworkspace-270d8f90-fake-46/log/abbfakea08368ca05fake60d8d1

> gics delete --workspace-id $(jq -r .workspace_id gics.json)
Schematics workspace "GICS Demo with a GH repository" (270d8f90-fake-46) has been deleted
```
