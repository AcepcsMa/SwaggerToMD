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
	title := analyzer.generator.GetHeader(model.Info.Title, H1, INDENT_0)
	overviewContent, overviewErr := analyzer.AnalyzeOverview(model)
	if overviewErr != nil {
		return "", overviewErr
	}
	pathsContent, pathsErr := analyzer.AnalyzePaths(model)
	if pathsErr != nil {
		return "", pathsErr
	}

	return fmt.Sprintf("%s\n%s\n%s", title, overviewContent, pathsContent), nil
}

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
	version := fmt.Sprintf("%s\n", swaggerModel.Info.Version)
	infoContent += versionHeader
	infoContent += version

	return infoContent
}

func (analyzer *SwaggerAnalyzer) FormatServers(swaggerModel *Model) string {
	return ""
}

func (analyzer *SwaggerAnalyzer) FormatTags(swaggerModel *Model) string {
	return ""
}

// analyze the overview part
func (analyzer *SwaggerAnalyzer) AnalyzeOverview(swaggerModel Model) (string, error) {

	overviewContent := make([]string, 0)

	overviewHeader := analyzer.generator.GetHeader(analyzer.terms["overview"], H2, INDENT_0)
	overview := swaggerModel.Info.Description + "\n"
	overviewContent = append(overviewContent, overviewHeader, overview)

	versionHeader := analyzer.generator.GetHeader(analyzer.terms["version_info"], H3,INDENT_0)
	version := fmt.Sprintf("Version: %s\n", swaggerModel.Info.Version)
	overviewContent = append(overviewContent, versionHeader, version)

	//uriHeader := analyzer.generator.GetHeader(analyzer.terms["uri_scheme"], H3, INDENT_0)
	//basePath := fmt.Sprintf("BasePath: %s\n", swaggerModel.BasePath)
	//overviewContent = append(overviewContent, uriHeader, basePath)

	tags := make([]string, len(swaggerModel.Tags))
	tagsHeader := analyzer.generator.GetHeader(analyzer.terms["tags"], H3, INDENT_0)
	for index, tag := range swaggerModel.Tags {
		format := "%s : %s"
		if index == len(swaggerModel.Tags) - 1 {
			format += "\n"
		}
		listItemContent := fmt.Sprintf(format, tag.Name, tag.Description)
		tags = append(tags, analyzer.generator.GetListItem(listItemContent, INDENT_0))
	}
	overviewContent = append(overviewContent, tagsHeader)
	overviewContent = append(overviewContent, tags...)

	//consumesHeader := analyzer.generator.GetHeader(analyzer.terms["consumes"], H3, INDENT_0)
	//consumes := make([]string, len(swaggerModel.Consumes))
	//for index, consume := range swaggerModel.Consumes {
	//	codeConsume := analyzer.generator.GetSingleLineCode(consume, INDENT_0)
	//	if index == len(swaggerModel.Consumes) - 1 {
	//		codeConsume += "\n"
	//	}
	//	consumes = append(consumes, analyzer.generator.GetListItem(codeConsume, INDENT_0))
	//}
	//overviewContent = append(overviewContent, consumesHeader)
	//overviewContent = append(overviewContent, consumes...)

	//producesHeader := analyzer.generator.GetHeader(analyzer.terms["produces"], H3, INDENT_0)
	//produces := make([]string, len(swaggerModel.Produces))
	//for _, produce := range swaggerModel.Produces {
	//	codeProduce := analyzer.generator.GetSingleLineCode(produce, INDENT_0)
	//	produces = append(produces, analyzer.generator.GetListItem(codeProduce, INDENT_0))
	//}
	//overviewContent = append(overviewContent, producesHeader)
	//overviewContent = append(overviewContent, produces...)

	finalOverviewContent := analyzer.compact(overviewContent)

	analyzer.content["overview"] = finalOverviewContent
	return analyzer.content["overview"], nil
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