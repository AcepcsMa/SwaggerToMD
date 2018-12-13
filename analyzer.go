package main

type Analyzer interface {
	Analyze(string) (string, error)
}

type SwaggerAnalyzer struct {
	content map[string]string
	generator *MdGenerator
}

func (analyzer *SwaggerAnalyzer) Analyze(jsonInput string) (string, error) {
	if analyzer.generator == nil {
		analyzer.generator = &MdGenerator{}
	}
	return "", nil
}

func (analyzer *SwaggerAnalyzer) AnalyzeOverview(jsonInput string) (string, error) {
	return "", nil
}

func (analyzer *SwaggerAnalyzer) AnalyzePaths(jsonInput string) (string, error) {
	return "", nil
}

func NewSwaggerAnalyzer() *SwaggerAnalyzer {
	analyzer := &SwaggerAnalyzer{}
	analyzer.content = make(map[string]string)
	analyzer.generator = NewMdGenerator()
	return analyzer
}