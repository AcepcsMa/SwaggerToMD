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

	NO_CONTENT = "No Content"
)

type MdGenerator struct {

}

type TableLine struct {
	Content map[string]string
}

func (line *TableLine) Get(header string) string {
	if content, ok := line.Content[header]; ok {
		return content
	} else {
		return NO_CONTENT
	}
}

func (line *TableLine) Set(key string, value string) {
	line.Content[key] = value
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

func (generator *MdGenerator) GetBoldLine(content string) string {
	return fmt.Sprintf("**%s**", content)
}

func (generator *MdGenerator) GetItalicLine(content string) string {
	return fmt.Sprintf("*%s*", content)
}

func (generator *MdGenerator) getTableLine(header []string, terms []string) TableLine {
	if len(header) != len(terms) {
		return TableLine{}
	}

	line := TableLine{}
	termCount := len(terms)
	for i := 0;i < termCount;i++ {
		line.Set(header[i], terms[i])
	}
	return line
}

func (generator *MdGenerator) GetTable(header []string, lines []TableLine) string {
	headerLine := ""
	headerSepLine := ""
	for _, colHeader := range header {
		headerLine += fmt.Sprintf("|%s", colHeader)
		headerSepLine += "|---"
	}
	headerSepLine += "|"
	headerLine += fmt.Sprintf("|\n%s", headerSepLine)

	lineContents := ""
	for _, line := range lines {
		currentLine := ""
		for _, colHeader := range header {
			currentLine += fmt.Sprintf("|%s", line.Get(colHeader))
		}
		currentLine += "|\n"
		lineContents += currentLine
	}

	return headerLine + lineContents
}

// factory for MdGenerator
func NewMdGenerator() *MdGenerator {
	return &MdGenerator{}
}