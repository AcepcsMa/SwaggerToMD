package main

import (
	"fmt"
	"strings"
)

const (
	H1 HeaderLevel = 1
	H2 HeaderLevel = 2
	H3 HeaderLevel = 3
	H4 HeaderLevel = 4
	H5 HeaderLevel = 5
)

type HeaderLevel int

type MdGenerator struct {

}

func (writer *MdGenerator) GetHeader(content string, level HeaderLevel) string {
	sharps := strings.Repeat("#", int(level))
	headerContent := fmt.Sprintf("%s %s", sharps, content)
	return headerContent
}

func NewMdGenerator() *MdGenerator {
	return &MdGenerator{}
}