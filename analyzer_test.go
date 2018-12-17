package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// test AnalyzeOverview in SwaggerAnalyzer
//func TestSwaggerAnalyzer_AnalyzeOverview(t *testing.T) {
//	t.Log("Test swagger analyzer - AnalyzeOverview")
//	{
//		model := Model{}
//		t.Log("Generate a swagger model")
//		{
//			fin, err := os.Open("test.json")
//			defer fin.Close()
//			if err != nil {
//				t.Fatal(err)
//			}
//			json.NewDecoder(fin).Decode(&model)
//		}
//		analyzer := NewSwaggerAnalyzer(ENGLISH)
//		overviewContent, err := analyzer.AnalyzeOverview(model)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println(overviewContent)
//	}
//}

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
