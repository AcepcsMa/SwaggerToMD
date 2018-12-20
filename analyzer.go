package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type LanguageType int

const (
	CHINESE LanguageType = 0
	ENGLISH LanguageType = 1
)

var parameterTableHeader = []string{"Type", "Name", "Description", "Schema"}
var responseTableHeader = []string{"HTTP Code", "Description", "Schema"}

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
	overviewContent, overviewErr := analyzer.AnalyzeOverview(model)
	if overviewErr != nil {
		return "", overviewErr
	}
	pathsContent, pathsErr := analyzer.AnalyzePaths(model)
	if pathsErr != nil {
		return "", pathsErr
	}

	return overviewContent + pathsContent, nil
}

// analyze the overview part
func (analyzer *SwaggerAnalyzer) AnalyzeOverview(swaggerModel Model) (string, error) {

	overviewContent := make([]string, 0)

	overviewHeader := analyzer.generator.GetHeader(analyzer.terms["overview"], H2)
	overview := swaggerModel.Info.Description
	overviewContent = append(overviewContent, overviewHeader, overview)

	versionHeader := analyzer.generator.GetHeader(analyzer.terms["version_info"], H4)
	version := fmt.Sprintf("Version: %s", swaggerModel.Info.Version)
	overviewContent = append(overviewContent, versionHeader, version)

	uriHeader := analyzer.generator.GetHeader(analyzer.terms["uri_scheme"], H4)
	basePath := fmt.Sprintf("BasePath: %s", swaggerModel.BasePath)
	overviewContent = append(overviewContent, uriHeader, basePath)

	tags := make([]string, len(swaggerModel.Tags))
	tagsHeader := analyzer.generator.GetHeader(analyzer.terms["tags"], H4)
	for _, tag := range swaggerModel.Tags {
		listItemContent := fmt.Sprintf("%s : %s", tag.Name, tag.Description)
		tags = append(tags, analyzer.generator.GetListItem(listItemContent, INDENT_0))
	}
	overviewContent = append(overviewContent, tagsHeader)
	overviewContent = append(overviewContent, tags...)

	consumesHeader := analyzer.generator.GetHeader(analyzer.terms["consumes"], H4)
	consumes := make([]string, len(swaggerModel.Consumes))
	for _, consume := range swaggerModel.Consumes {
		codeConsume := analyzer.generator.GetSingleLineCode(consume)
		consumes = append(consumes, analyzer.generator.GetListItem(codeConsume, INDENT_0))
	}
	overviewContent = append(overviewContent, consumesHeader)
	overviewContent = append(overviewContent, consumes...)

	producesHeader := analyzer.generator.GetHeader(analyzer.terms["produces"], H4)
	produces := make([]string, len(swaggerModel.Produces))
	for _, produce := range swaggerModel.Produces {
		codeProduce := analyzer.generator.GetSingleLineCode(produce)
		produces = append(produces, analyzer.generator.GetListItem(codeProduce, INDENT_0))
	}
	overviewContent = append(overviewContent, producesHeader)
	overviewContent = append(overviewContent, produces...)

	finalOverviewContent := analyzer.compact(overviewContent)

	analyzer.content["overview"] = finalOverviewContent
	return analyzer.content["overview"], nil
}

// analyze the paths part
func (analyzer *SwaggerAnalyzer) AnalyzePaths(swaggerModel Model) (string, error) {
	pathsContent := make([]string, 0)
	pathsJson := swaggerModel.Paths

	pathsHeader := analyzer.generator.GetHeader("Paths", H2)
	pathsContent = append(pathsContent, pathsHeader)

	for apiPath, methods := range pathsJson {
		apis := analyzer.ExtractAPIs(apiPath, methods.(map[string]interface{}))
		for _, api := range apis {
			apiInMd := analyzer.FormatAPI(api)
			pathsContent = append(pathsContent, apiInMd)
		}
	}

	finalPathsContent := analyzer.compact(pathsContent)
	return finalPathsContent, nil
}

func (analyzer *SwaggerAnalyzer) FormatAPI(api Api) string {
	apiContent := ""
	boldDescription := analyzer.generator.GetBoldLine(api.OperationId)
	description := analyzer.generator.GetItalicLine(boldDescription)
	apiContent += fmt.Sprintf("%s\n", description)

	codePath := analyzer.generator.GetMultiLineCode(api.Path)
	apiContent += fmt.Sprintf("%s\n", codePath)

	parameterHeader := analyzer.generator.GetHeader("Parameters", H4)
	pTableLines := make([]TableLine, len(api.Parameters))
	for _, parameter := range api.Parameters {
		currentLine := TableLine{Content: make(map[string]string)}
		currentLine.Set("Type", parameter.In)
		currentLine.Set("Name", parameter.Name)
		currentLine.Set("Description", parameter.Description)
		currentLine.Set("Schema", parameter.Type)
		pTableLines = append(pTableLines, currentLine)
	}
	parameterTable := analyzer.generator.GetTable(parameterTableHeader, pTableLines)
	apiContent += fmt.Sprintf("%s\n%s\n", parameterHeader, parameterTable)

	responseHeader := analyzer.generator.GetHeader("Responses", H4)
	rTableLines := make([]TableLine, len(api.Response))
	for key, value := range api.Response {
		currentLine := TableLine{Content: make(map[string]string)}
		currentLine.Set("HTTP Code", key)
		currentLine.Set("Description", value)
		rTableLines = append(rTableLines, currentLine)
	}
	responseTable := analyzer.generator.GetTable(responseTableHeader, rTableLines)
	apiContent += fmt.Sprintf("%s\n%s\n", responseHeader, responseTable)

	TagHeader := analyzer.generator.GetHeader("Tags", H4)
	apiContent += fmt.Sprintf("%s\n", TagHeader)
	for _, tag := range api.Tags {
		currentListItem := analyzer.generator.GetListItem(tag, INDENT_0)
		apiContent += fmt.Sprintf("%s\n", currentListItem)
	}

	return apiContent
}

// extract APIs from a given method formatted in Json
func (analyzer *SwaggerAnalyzer) ExtractAPIs(apiPath string, methods map[string]interface{}) []Api {
	apis := make([]Api, len(methods))

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