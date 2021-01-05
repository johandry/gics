package v1

// To re-generate or update the client.gen.go file:
// 1. Download the JSON OpenAPI
// 	1.1 Go to https://cloud.ibm.com/apidocs/schematics
// 	1.2 Select the 3 vertical dots at the top of the left menu
// 	1.3 Click on "Download OpenAPI definition"
// 2. Replace the schematics.json file
// 3. Execute: go generate

// Run oapi-codegen to regenerate the schematics boilerplate version 1.0
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=v1 --generate types,client -o ./client.gen.go ./schematics.json

// APIVersion is the API version, it is located in the JSON OpenAPI in info.version
const APIVersion = "1.0"
