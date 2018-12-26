package main

import "fmt"

// Parameter struct, demonstrating the essential info of a parameter for an API
type Parameter struct {
	Description string
	Name string
	Type string
	In string
	Example string
}

func (p Parameter) String() string {
	return fmt.Sprintf("{\n\tDescription: %s\n\tName: %s\n\tType: %s\n\tIn: %s\n\tExample: %s\n}",
		p.Description, p.Name, p.Type, p.In, p.Example)
}

// Responses struct
type Response struct {
	StatusCode string
	Description string
	Schema string
}

func (r Response) String() string {
	return fmt.Sprintf("StatusCode: %s, Description: %s, SchemaType: %s",
		r.StatusCode, r.Description, r.Schema)
}

// API struct, demonstrating the essential info of an API
type Api struct {
	Path        string
	Method      string
	Responses   []Response
	ResponseInJson string
	OperationId string
	Parameters  []Parameter
	Tags        []string
}

func (api Api) String() string {
	return fmt.Sprintf("{\n\tPath: %s\n\tMethod: %s\n\tResponses: %v\n\t" +
		"OperationId: %s\n\tParameters: %v\n\tTags: %v\n\tResponseInJson: %s\n}",
		api.Path, api.Method, api.Responses, api.OperationId, api.Parameters, api.Tags, api.ResponseInJson)
}
