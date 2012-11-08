package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Param struct {
	Param string `xml:",innerxml"`
}

type TestSuite struct {
	TestSuiteName string
	TestName      string
	InputType     string
	Params        []Param `xml:"params>param"`
}

func main() {
	mp := Param{
		"ClientType": "1",
		"Version":    "10",
	}
	v := TestSuite{
		TestSuiteName: "Test test",
		TestName:      "John",
		InputType:     "HelloInfo",
		InputData: []map[string]string{
			mp,
		},
	}
	output, err := xml.MarshalIndent(v, " ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)
}
