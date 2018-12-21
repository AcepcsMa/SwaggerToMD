package main

type Model struct {

	OpenApi string `json:"openapi"`
	//Swagger string `json:"swagger"`
	//BasePath string `json:"basePath"`

	Info struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Contact struct {
			Email string `json:"email"`
		} `json:"contact"`
		License struct {
			Name string `json:"name"`
			Url string `json:"url"`
		} `json:"license"`
		Version string `json:"version"`
	} `json:"info"`

	Servers []struct {
		Url string `json:"url"`
		Description string `json:"description"`
	} `json:"servers"`

	//Produces []string `json:"produces"`
	//Consumes []string `json:"consumes"`
	Tags []struct {
		Name string `json:"name"`
		Description string `json:"description"`
	} `json:"tags"`

	Paths map[string]interface{} `json:"paths"`

	Components struct{
		Schemas map[string]interface{} `json:"schemas"`
	} `json:"components"`

	//Responses map[string]interface{} `json:"responses"`

}
