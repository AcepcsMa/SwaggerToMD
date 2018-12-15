package main

import (
	"fmt"
	"strings"
)

type HeaderLevel int
type IndentLevel int

const (
	H1 HeaderLevel = 1
	H2 HeaderLevel = 2
	H3 HeaderLevel = 3
	H4 HeaderLevel = 4
	H5 HeaderLevel = 5

	INDENT_0 IndentLevel = 0
	INDENT_1 IndentLevel = 1
	INDENT_2 IndentLevel = 2
)

type MdGenerator struct {

}

func (generator *MdGenerator) GetHeader(content string, level HeaderLevel) string {
	sharps := strings.Repeat("#", int(level))
	headerContent := fmt.Sprintf("%s %s", sharps, content)
	return headerContent
}

func (generator *MdGenerator) GetListItem(content string, level IndentLevel) string {
	indent := strings.Repeat(" ", int(level) * 4)
	listItemContent := fmt.Sprintf("%s+ %s", indent, content)
	return listItemContent
}

func NewMdGenerator() *MdGenerator {
	return &MdGenerator{}
}