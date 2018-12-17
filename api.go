package main

import "fmt"

// Parameter struct, demonstrating the essential info of a parameter for an API
type Parameter struct {
	Description string
	Name string
	Type string
	In string
}

// API struct, demonstrating the essential info of an API
type Api struct {
	Path string
	Method string
	Response map[string]string
	OperationId string
	Parameters []Parameter
	Tags []string
}

func (p Parameter) String() string {
	return fmt.Sprintf("{\n\tDescription: %s\n\tName: %s\n\tType: %s\n\tIn: %s\n}",
		p.Description, p.Name, p.Type, p.In)
}

func (api Api) String() string {
	return fmt.Sprintf("{\n\tPath: %s\n\tMethod: %s\n\tResponses: %v\n\t" +
		"OperationId: %s\n\tParameters: %v\n\tTags: %v\n}",
		api.Path, api.Method, api.Response, api.OperationId, api.Parameters, api.Tags)
}
