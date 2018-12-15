package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

// test AnalyzeOverview in SwaggerAnalyzer
func TestSwaggerAnalyzer_AnalyzeOverview(t *testing.T) {
	t.Log("Test swagger analyzer")
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
