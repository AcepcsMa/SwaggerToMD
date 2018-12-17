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

// generate a header in markdown
func (generator *MdGenerator) GetHeader(content string, level HeaderLevel) string {
	sharps := strings.Repeat("#", int(level))
	headerContent := fmt.Sprintf("%s %s", sharps, content)
	return headerContent
}

// generate a list item in markdown
func (generator *MdGenerator) GetListItem(content string, level IndentLevel) string {
	indent := strings.Repeat(" ", int(level) * 4)
	listItemContent := fmt.Sprintf("%s+ %s", indent, content)
	return listItemContent
}

// generate a single line of code in markdown
func (generator *MdGenerator) GetSingleLineCode(content string) string {
	return fmt.Sprintf("`%s`", content)
}

// generate multiple lines of code in markdown
func (generator *MdGenerator) GetMultiLineCode(content string) string {
	return fmt.Sprintf("```\n%s\n```", content)
}

// factory for MdGenerator
func NewMdGenerator() *MdGenerator {
	return &MdGenerator{}
}