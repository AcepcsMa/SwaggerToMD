package main

import "errors"

type Transformer struct {
	Input       string
	Output      string
	ContentFrom ContentSource
	LangType 	LanguageType
	JsonContent string
	MdContent   string

	contentGetter ContentGetter
	analyzer      Analyzer
}

func (t *Transformer) GetContent() error {
	if t.contentGetter == nil {
		t.contentGetter = NewSwaggerContentGetter(t.Input, t.ContentFrom)
	}

	content, err := t.contentGetter.GetContent()
	if err != nil {
		return err
	}
	t.JsonContent = content
	return nil
}

func (t *Transformer) Analyze() error {
	if t.analyzer == nil {
		t.analyzer = NewSwaggerAnalyzer(ENGLISH)
	}

	if len(t.JsonContent) == 0 {
		return errors.New("empty json content")
	}

	result, err := t.analyzer.Analyze(t.JsonContent)
	if err != nil {
		return err
	}
	t.MdContent = result
	return nil
}

func (t *Transformer) WriteToOutput() {

}

func NewTransformer(input string, output string, contentSource ContentSource, langType LanguageType) *Transformer {
	transformer := &Transformer{Input:input, Output:output, LangType:langType}
	transformer.contentGetter = NewSwaggerContentGetter(input, contentSource)
	transformer.analyzer = NewSwaggerAnalyzer(langType)
	return transformer
}