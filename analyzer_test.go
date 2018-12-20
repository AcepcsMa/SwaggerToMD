package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		overviewContent, err := analyzer.AnalyzeOverview(model)
		if err != nil {
			log.Fatal(err)
		}
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