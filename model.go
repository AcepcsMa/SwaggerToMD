package main

type Model struct {

	Swagger string `json:"swagger"`
	BasePath string `json:"basePath"`

	Paths map[string]interface{} `json:"paths"`

	Info struct {
		Title string `json:"title"`
		Version string `json:"version"`
		Description string `json:"description"`
	}

	Produces []string `json:"produces"`
	Consumes []string `json:"consumes"`
	Tags []struct {
		Name string `json:"name"`
		Description string `json:"description"`
	} `json:"tags"`

	Responses map[string]interface{} `json:"responses"`

}
