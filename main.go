package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	localInput string
	webInput string
	lang string
	output string
)

func main() {

	flag.StringVar(&localInput, "local", "./", "Local path of the input json.")
	flag.StringVar(&webInput, "web", "", "Web url of the input json.")
	flag.StringVar(&lang, "lang", "en", "Language of the output markdown doc.")
	flag.StringVar(&output, "out", "./", "Output file name.")

	flag.Parse()

	fin, err := os.Open(localInput)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(fin)

	model := Model{}
	err = json.NewDecoder(reader).Decode(&model)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", model.Info.Title)

	markdownWriter := &MdGenerator{FileName: output}
	markdownWriter.WriteHeader(model.Info.Title, H2)
	markdownWriter.Close()
}
