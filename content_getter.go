package main

import "errors"

type ContentSource int

const (
	LOCAL_SOURCE ContentSource = 0
	WEB_SOURCE     ContentSource = 1
)

// Invalid content source, it's unable to retrieve content from this source
var InvalidContentSource = errors.New("invalid content source")

// An interface retrieving content from a local file OR a web url
type ContentGetter interface {
	GetContent() (string, error)
	GetLocalContent() (string, error)
	GetWebContent() (string, error)
}

type SwaggerContentGetter struct {
	contentPath   string
	contentSource ContentSource
}

func (scg *SwaggerContentGetter) GetContent() (string, error) {
	switch scg.contentSource {
	case LOCAL_SOURCE:
		return scg.GetLocalContent()
	case WEB_SOURCE:
		return scg.GetWebContent()
	default:
		return "", InvalidContentSource
	}
}

func (scg *SwaggerContentGetter) GetLocalContent() (string, error) {
	return "", nil
}

func (scg *SwaggerContentGetter) GetWebContent() (string, error) {
	return "", nil
}

func NewSwaggerContentGetter(contentPath string, origin ContentSource) *SwaggerContentGetter {
	getter := &SwaggerContentGetter{contentPath: contentPath, contentSource:origin}
	return getter
}