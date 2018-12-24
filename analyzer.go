package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type LanguageType int

const (
	CHINESE LanguageType = 0
	ENGLISH LanguageType = 1

	TYPE = "Type"
	DESCRIPTION = "Description"
	SCHEMA = "Schema"
	NAME = "Name"
	HTTP_CODE = "HTTP Code"
	PROPERTY_NAME = "Property Name"
	PROPERTY_TYPE = "Property Type"
	REQUIRED = "Required"
	EXAMPLE = "Example"
	TRUE = "True"
	FALSE = "False"
)

var parameterTableHeader = []string{"Type", "Name", "Description", "Schema"}
var responseTableHeader = []string{"HTTP Code", "Description", "Schema"}
var componentTableHeader = []string{"Property Name","Property Type", "Required", "Example"}

type Analyzer interface {
	Analyze(string) (string, error)
}

type SwaggerAnalyzer struct {
	content map[string]string	// content
	terms map[string]string		// terms associated with language settings
	generator *MdGenerator		// markdown format generator
}

// set language of the SwaggerAnalyzer
func (analyzer *SwaggerAnalyzer) SetLang(lang LanguageType) error {
	configFile := ""
	if lang == CHINESE {
		configFile = "zh_config.json"
	} else if lang == ENGLISH {
		configFile = "en_config.json"
	}

	fin, err := os.Open(configFile)
	defer fin.Close()
	if err != nil {
		return err
	}
	json.NewDecoder(fin).Decode(&analyzer.terms)
	return nil
}

// the main entrance of analysis
func (analyzer *SwaggerAnalyzer) Analyze(jsonInput string) (string, error) {
	if analyzer.generator == nil {
		analyzer.generator = NewMdGenerator()
	}

	model := Model{}
	json.Unmarshal([]byte(jsonInput), &model)
	title := analyzer.generator.GetHeader(model.Info.Title, H1, INDENT_0)
	overviewContent := analyzer.AnalyzeOverview(&model)
	pathsContent, pathsErr := analyzer.AnalyzePaths(model)
	if pathsErr != nil {
		return "", pathsErr
	}

	return fmt.Sprintf("%s\n%s\n%s", title, overviewContent, pathsContent), nil
}

// format info section in swagger json doc
func (analyzer *SwaggerAnalyzer) FormatInfo(swaggerModel *Model) string {
	description := analyzer.generator.GetBoldLine(swaggerModel.Info.Description)
	infoContent := fmt.Sprintf("%s\n", description)

	contactHeader := analyzer.generator.GetHeader(analyzer.terms["contact"] + "\n", H3, INDENT_0)
	infoContent += contactHeader
	for key, value := range swaggerModel.Info.Contact {
		currentLine := fmt.Sprintf("%s : %s\n", key, value)
		currentListItem := analyzer.generator.GetListItem(currentLine, INDENT_0)
		infoContent += currentListItem
	}
	infoContent += "\n"

	licenseHeader := analyzer.generator.GetHeader(analyzer.terms["license"] + "\n", H3, INDENT_0)
	infoContent += licenseHeader
	for key, value := range swaggerModel.Info.License {
		currentLine := fmt.Sprintf("%s : %s\n", key, value)
		currentListItem := analyzer.generator.GetListItem(currentLine, INDENT_0)
		infoContent += currentListItem
	}
	infoContent += "\n"

	versionHeader := analyzer.generator.GetHeader(analyzer.terms["version"] + "\n", H3, INDENT_0)
	version := fmt.Sprintf("%s\n\n", swaggerModel.Info.Version)
	infoContent += versionHeader
	infoContent += version

	return infoContent
}

// format servers section in swagger json doc
func (analyzer *SwaggerAnalyzer) FormatServers(swaggerModel *Model) string {
	serversContent := analyzer.generator.GetHeader(analyzer.terms["servers"], H3, INDENT_0)
	serversContent += "\n"

	for index, server := range swaggerModel.Servers {
		currentServerHeader := analyzer.generator.GetListItem(fmt.Sprintf("Server-%d\n", index), INDENT_0)
		currentServerUrl := analyzer.generator.GetListItem(fmt.Sprintf("url : %s\n", server.Url), INDENT_1)
		currentServerDesc := analyzer.generator.GetListItem(fmt.Sprintf("description : %s\n",
			server.Description), INDENT_1)
		serversContent += currentServerHeader
		serversContent += currentServerUrl
		serversContent += currentServerDesc
	}
	serversContent += "\n"

	return serversContent
}

// format tags section in swagger json doc
func (analyzer *SwaggerAnalyzer) FormatTags(swaggerModel *Model) string {
	tagsContent := analyzer.generator.GetHeader(analyzer.terms["tags"], H3, INDENT_0)
	tagsContent += "\n"

	for _, tag := range swaggerModel.Tags {
		boldTagName := analyzer.generator.GetBoldLine(tag.Name)
		tagName := analyzer.generator.GetItalicLine(boldTagName)
		tagDesc := tag.Description
		currentListItem := analyzer.generator.GetListItem(fmt.Sprintf("%s : %s\n", tagName, tagDesc), INDENT_0)
		tagsContent += currentListItem
	}

	return tagsContent
}

// analyze the overview part
func (analyzer *SwaggerAnalyzer) AnalyzeOverview(swaggerModel *Model) string {
	overviewContent := analyzer.generator.GetHeader(analyzer.terms["overview"], H2, INDENT_0)
	overviewContent += "\n"
	overviewContent += analyzer.FormatInfo(swaggerModel)
	overviewContent += analyzer.FormatServers(swaggerModel)
	overviewContent += analyzer.FormatTags(swaggerModel)
	return overviewContent
}

func (analyzer *SwaggerAnalyzer) AnalyzeComponents(swaggerModel *Model) string {
	components := analyzer.ExtractComponents(swaggerModel)
	componentsContent := fmt.Sprintf("%s\n",
		analyzer.generator.GetHeader(analyzer.terms["components"], H2, INDENT_0))

	componentsContent += fmt.Sprintf("%s\n", analyzer.FormatComponents(components))
	return componentsContent
}

// analyze the paths part
func (analyzer *SwaggerAnalyzer) AnalyzePaths(swaggerModel Model) (string, error) {
	pathsContent := make([]string, 0)
	pathsJson := swaggerModel.Paths

	pathsHeader := analyzer.generator.GetHeader(analyzer.terms["paths"], H2, INDENT_0)
	pathsContent = append(pathsContent, pathsHeader)

	apiIndex := 1
	for apiPath, methods := range pathsJson {
		apis := analyzer.ExtractAPIs(apiPath, methods.(map[string]interface{}))
		for _, api := range apis {
			apiInMd := analyzer.FormatAPI(apiIndex, api)
			pathsContent = append(pathsContent, apiInMd)
			apiIndex++
		}
	}

	finalPathsContent := analyzer.compact(pathsContent)
	return finalPathsContent, nil
}

func (analyzer *SwaggerAnalyzer) FormatAPI(apiIndex int, api Api) string {
	apiContent := ""
	boldDescription := analyzer.generator.GetBoldLine(api.OperationId)
	description := analyzer.generator.GetItalicLine(boldDescription)
	apiContent += fmt.Sprintf("%d. %s\n\n", apiIndex, description)

	codePath := analyzer.generator.GetMultiLineCode(fmt.Sprintf("%s %s",
		strings.ToUpper(api.Method), api.Path), INDENT_1)
	apiContent += fmt.Sprintf("%s\n", codePath)

	if len(api.Parameters) > 0 {
		parameterHeader := analyzer.generator.GetHeader(analyzer.terms["parameters"], H4, INDENT_1)
		pTableLines := make([]TableLine, 0, len(api.Parameters))
		for _, parameter := range api.Parameters {
			currentLine := TableLine{Content: make(map[string]string)}
			currentLine.Set(TYPE, parameter.In)
			currentLine.Set(NAME, parameter.Name)
			currentLine.Set(DESCRIPTION, parameter.Description)
			currentLine.Set(SCHEMA, parameter.Type)
			pTableLines = append(pTableLines, currentLine)
		}
		parameterTable := analyzer.generator.GetTable(parameterTableHeader, pTableLines, INDENT_1)
		apiContent += fmt.Sprintf("%s\n%s\n", parameterHeader, parameterTable)
	}

	if len(api.Response) > 0 {
		responseHeader := analyzer.generator.GetHeader(analyzer.terms["responses"], H4, INDENT_1)
		rTableLines := make([]TableLine, 0, len(api.Response))
		for key, value := range api.Response {
			currentLine := TableLine{Content: make(map[string]string)}
			currentLine.Set(HTTP_CODE, key)
			currentLine.Set(DESCRIPTION, value)
			rTableLines = append(rTableLines, currentLine)
		}
		responseTable := analyzer.generator.GetTable(responseTableHeader, rTableLines, INDENT_1)
		apiContent += fmt.Sprintf("%s\n%s\n", responseHeader, responseTable)
	}

	TagHeader := analyzer.generator.GetHeader(analyzer.terms["tags"], H4, INDENT_1)
	apiContent += fmt.Sprintf("%s\n", TagHeader)
	for _, tag := range api.Tags {
		currentListItem := analyzer.generator.GetListItem(tag, INDENT_1)
		apiContent += fmt.Sprintf("%s\n", currentListItem)
	}

	return apiContent
}

// format a slice of components
func (analyzer *SwaggerAnalyzer) FormatComponents(components []Component) string {
	componentsContent := ""
	for _, component := range components {
		componentsContent += fmt.Sprintf("%s\n", analyzer.FormatComponent(&component))
	}
	return componentsContent
}

// format a single component
func (analyzer *SwaggerAnalyzer) FormatComponent(component *Component) string {
	componentContent := fmt.Sprintf("%s\n", analyzer.generator.GetListItem(component.Name, INDENT_0))
	typeInCode := analyzer.generator.GetSingleLineCode(component.Type, INDENT_0)
	componentContent += fmt.Sprintf("%s\n\n", analyzer.generator.GetListItem("type : " + typeInCode, INDENT_1))
	tableLines := make([]TableLine, 0, len(component.Properties))
	for _, property := range component.Properties {
		currentLine := TableLine{Content: make(map[string]string)}
		currentLine.Set(PROPERTY_NAME, property.Name)
		currentLine.Set(PROPERTY_TYPE, property.Type)
		if property.Required {
			currentLine.Set(REQUIRED, TRUE)
		} else {
			currentLine.Set(REQUIRED, FALSE)
		}
		currentLine.Set(EXAMPLE, property.Example)
		tableLines = append(tableLines, currentLine)
	}
	componentContent += fmt.Sprintf("%s\n",
		analyzer.generator.GetTable(componentTableHeader, tableLines, INDENT_2))
	componentContent += fmt.Sprintf("%s\n",
		analyzer.generator.GetMultiLineCode(component.Code, INDENT_2))
	return componentContent
}

// extract components
func (analyzer *SwaggerAnalyzer) ExtractComponents(swaggerModel *Model) []Component {
	components := make([]Component, 0, len(swaggerModel.Components.Schemas))

	for componentName, component := range swaggerModel.Components.Schemas {
		currentComponent := Component{Name: componentName, Type: component.(map[string]interface{})["type"].(string)}
		required := make(map[string]bool)
		requiredFields := component.(map[string]interface{})["required"].([]interface{})
		for _, requiredField := range requiredFields {
			required[requiredField.(string)] = true
		}

		properties := component.(map[string]interface{})["properties"].(map[string]interface{})
		currentProperties := make([]Property, 0, len(properties))
		for propertyName, property := range properties {
			currentProperty := Property{Name: propertyName,
			Type: property.(map[string]interface{})["type"].(string)}
			if example, ok := property.(map[string]interface{})["example"]; ok {
				currentProperty.Example = fmt.Sprintf("%v", example)
			} else {
				currentProperty.Example = "/"
			}
			if isRequired, ok := required[propertyName]; ok {
				currentProperty.Required = isRequired
			}
			if currentProperty.Type == "array" {
				arrayType := property.(map[string]interface{})["items"].(map[string]interface{})["type"].(string)
				currentProperty.Type = fmt.Sprintf("array\\<%s\\>", arrayType)
			}
			currentProperty.Type = analyzer.generator.GetItalicLine(currentProperty.Type)
			currentProperties = append(currentProperties, currentProperty)
		}
		code, err := json.MarshalIndent(component, "", "    ")
		if err != nil {
			panic(err)
		}
		currentComponent.Code = string(code)
		currentComponent.Properties = currentProperties
		components = append(components, currentComponent)
	}
	return components
}

// extract APIs from a given method formatted in Json
func (analyzer *SwaggerAnalyzer) ExtractAPIs(apiPath string, methods map[string]interface{}) []Api {
	apis := make([]Api, 0, len(methods))

	for methodName, value := range methods {
		currentApi := Api{}
		currentApi.Response = make(map[string]string)
		currentApi.Path = apiPath
		currentApi.Method = methodName
		responses := value.(map[string]interface{})["responses"].(map[string]interface{})
		for statusCode, returnInfo := range responses {
			currentApi.Response[statusCode] = returnInfo.(map[string]interface{})["description"].(string)
		}
		currentApi.OperationId = value.(map[string]interface{})["operationId"].(string)

		if parameters, ok := value.(map[string]interface{})["parameters"].([]interface{}); ok {
			for _, parameter := range parameters {

				currentParameter := Parameter{
					Description: parameter.(map[string]interface{})["description"].(string),
					Name: parameter.(map[string]interface{})["name"].(string),
					Type: parameter.(map[string]interface{})["type"].(string),
					In: parameter.(map[string]interface{})["in"].(string)}
				currentApi.Parameters = append(currentApi.Parameters, currentParameter)
			}
		}

		tags := value.(map[string]interface{})["tags"].([]interface{})
		for _, tag := range tags {
			currentApi.Tags = append(currentApi.Tags, tag.(string))
		}
		apis = append(apis, currentApi)
	}
	return apis
}

// compact means removing empty & useless line contents
func (analyzer *SwaggerAnalyzer) compact(content []string) string {
	compactedContent := ""
	for _, line := range content {
		if len(line) > 0 {
			compactedContent += fmt.Sprintf("%s\n", line)
		}
	}
	return compactedContent
}

// factory for SwaggerAnalyzer
func NewSwaggerAnalyzer(lang LanguageType) *SwaggerAnalyzer {
	analyzer := &SwaggerAnalyzer{}
	analyzer.content = make(map[string]string)
	analyzer.generator = NewMdGenerator()
	err := analyzer.SetLang(lang)
	if err != nil {
		log.Fatal("language setting error, only support zh or en now")
	}
	return analyzer
}