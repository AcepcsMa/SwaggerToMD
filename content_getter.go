package main

type ContentGetter interface {
	GetLocalContent() (string, error)
	GetWebContent() (string, error)
}

type SwaggerContentGetter struct {
	ContentPath string
}

func (cg *SwaggerContentGetter) GetContent() (string, error) {
	return "", nil
}

func (cg *SwaggerContentGetter) GetLocalContent() (string, error) {
	return "", nil
}

func (cg *SwaggerContentGetter) GetWebContent() (string, error) {
	return "", nil
}

