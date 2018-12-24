package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// test AnalyzeOverview in SwaggerAnalyzer
func TestSwaggerAnalyzer_AnalyzeOverview(t *testing.T) {
	t.Log("Test swagger analyzer - AnalyzeOverview")
	{
		model := Model{}
		t.Log("Generate a swagger model")
		{
			fin, err := os.Open("test.json")
			defer fin.Close()
			if err != nil {
				t.Fatal(err)
			}
			json.NewDecoder(fin).Decode(&model)
		}
		analyzer := NewSwaggerAnalyzer(ENGLISH)
		overviewContent := analyzer.AnalyzeOverview(&model)
		fmt.Println(overviewContent)
	}
}

// test ExtractAPIs in SwaggerAnalyzer
func TestSwaggerAnalyzer_ExtractAPIs(t *testing.T) {
	t.Log("Test swagger analyzer - ExtractAPIs")
	{
		model := Model{}
		t.Log("Generate a swagger model")
		{
			fin, err := os.Open("test.json")
			defer fin.Close()
			if err != nil {
				t.Fatal(err)
			}
			json.NewDecoder(fin).Decode(&model)
		}
		analyzer := NewSwaggerAnalyzer(ENGLISH)
		paths := model.Paths
		for apiPath, methods := range paths {
			apis := analyzer.ExtractAPIs(apiPath, methods.(map[string]interface{}))
			for _, api := range apis {
				fmt.Println(api)
			}
		}
	}
}

// test FormatAPI in SwaggerAnalyzer
func TestSwaggerAnalyzer_FormatAPI(t *testing.T) {
	t.Log("Test swagger analyzer - FormatAPI")
	{
		model := Model{}
		t.Log("Generate a swagger model")
		{
			fin, err := os.Open("test.json")
			defer fin.Close()
			if err != nil {
				t.Fatal(err)
			}
			json.NewDecoder(fin).Decode(&model)
		}
		analyzer := NewSwaggerAnalyzer(ENGLISH)
		paths := model.Paths
		for apiPath, methods := range paths {
			apis := analyzer.ExtractAPIs(apiPath, methods.(map[string]interface{}))
			for index, api := range apis {
				t.Log(fmt.Sprintf("API #%d\n", index))
				{
					fmt.Println(analyzer.FormatAPI(index, api))
				}
			}
		}
	}
}

// test Analyze in SwaggerAnalyzer
func TestSwaggerAnalyzer_Analyze(t *testing.T) {
	t.Log("Test swagger analyzer - Analyze")
	{
		t.Log("Read json from swagger test file")
		{
			jsonString, readErr := ioutil.ReadFile("test.json")
			if readErr != nil {
				t.Fatal(readErr)
			}
			t.Log("Analyze")
			{
				analyzer := NewSwaggerAnalyzer(ENGLISH)
				result, analyzeErr := analyzer.Analyze(string(jsonString))
				if analyzeErr != nil {
					t.Fatal(analyzeErr)
				}
				//fmt.Println(result)
				resultFile, resultFileErr := os.OpenFile("test.md", os.O_RDWR | os.O_CREATE, 0755)
				defer resultFile.Close()
				if resultFileErr != nil {
					t.Fatal(resultFileErr)
				}
				resultWriter := bufio.NewWriter(resultFile)
				resultWriter.WriteString(result)
				resultWriter.Flush()
			}
		}
	}
}

// test FormatInfo in SwaggerAnalyzer
func TestSwaggerAnalyzer_FormatInfo(t *testing.T) {
	t.Log("Test SwaggerAnalyzer for OAS3.0 - FormatInfo")
	{
		model := Model{}
		t.Log("Generate a swagger model")
		{
			fin, err := os.Open("test.json")
			defer fin.Close()
			if err != nil {
				t.Fatal(err)
			}
			json.NewDecoder(fin).Decode(&model)
		}
		analyzer := NewSwaggerAnalyzer(ENGLISH)
		result := analyzer.FormatInfo(&model)
		fmt.Println(result)
	}
}

// test FormatServers in SwaggerAnalyzer
func TestSwaggerAnalyzer_FormatServers(t *testing.T) {
	t.Log("Test SwaggerAnalyzer for OAS3.0 - FormatServers")
	{
		model := Model{}
		t.Log("Generate a swagger model")
		{
			fin, err := os.Open("test.json")
			defer fin.Close()
			if err != nil {
				t.Fatal(err)
			}
			json.NewDecoder(fin).Decode(&model)
		}
		analyzer := NewSwaggerAnalyzer(ENGLISH)
		result := analyzer.FormatServers(&model)
		fmt.Println(result)
	}
}

// test FormatTags in SwaggerAnalyzer
func TestSwaggerAnalyzer_FormatTags(t *testing.T) {
	t.Log("Test SwaggerAnalyzer for OAS3.0 - FormatTags")
	{
		model := Model{}
		t.Log("Generate a swagger model")
		{
			fin, err := os.Open("test.json")
			defer fin.Close()
			if err != nil {
				t.Fatal(err)
			}
			json.NewDecoder(fin).Decode(&model)
		}
		analyzer := NewSwaggerAnalyzer(ENGLISH)
		result := analyzer.FormatTags(&model)
		fmt.Println(result)
	}
}

// test ExtractComponents in SwaggerAnalyzer
func TestSwaggerAnalyzer_ExtractComponents(t *testing.T) {
	t.Log("Test SwaggerAnalyzer for OAS3.0 - ExtractComponents")
	{
		model := Model{}
		t.Log("Generate a swagger model")
		{
			fin, err := os.Open("test.json")
			defer fin.Close()
			if err != nil {
				t.Fatal(err)
			}
			json.NewDecoder(fin).Decode(&model)
		}
		analyzer := NewSwaggerAnalyzer(ENGLISH)
		components := analyzer.ExtractComponents(&model)
		for _, component := range components {
			fmt.Println(component)
		}
	}
}