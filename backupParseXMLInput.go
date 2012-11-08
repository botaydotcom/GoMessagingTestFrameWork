package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"reflect"
	//"code.google.com/p/goprotobuf/proto"
	"xmltest/btalkTest/GeneratedDataStructure"
)

type Message struct {
	MessageType string `xml:"type,attr"`
	Data        string `xml:",innerxml"`
}

type TestSuite struct {
	TestSuiteName string
	TestName      string
	InputType     string
	InputData     []Message `xml:"InputData>Message"`
	OutputData    []Message `xml:"OutputData>Message"`
}

var inputMessages []interface{}
var outputMessages []interface{}

func main() {
	// open file
	data, err := ioutil.ReadFile("test.xml")
	if err != nil {
		panic(err)
	}

	var ts TestSuite
	err = xml.Unmarshal(data, &ts)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Printf("Readed input %v", ts.InputData)

	inputMessages = make([]interface{}, 10)
	outputMessages = make([]interface{}, 10)

	for i, v := range ts.InputData {
		fmt.Printf("Readed input Number %d \n----------------------------\n%v\n---------------------------\n", i, v.Data)

		input := MagicVarFunc(v.MessageType)

		err = xml.Unmarshal([]byte(v.Data), input)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Printf("\n----------------------------\n Parsed input %v\n---------------------------\n", input)

		s := reflect.ValueOf(input).Elem()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			printValue(f)
		}
		append(inputMessages, input)
	}

	for i, v := range ts.OutputData {
		fmt.Printf("Readed output Number %d \n----------------------------\n%v\n---------------------------\n", i, v.Data)

		message := MagicVarFunc(v.MessageType)

		err = xml.Unmarshal([]byte(v.Data), message)

		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
		fmt.Printf("\n----------------------------\n Parsed output %v\n---------------------------\n", message)
		s := reflect.ValueOf(message).Elem()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			printValue(f)
		}
		append(outputMessages, message)
	}

}

func printValue(ptrValue reflect.Value) {
	if ptrValue.Kind() == reflect.Ptr {
		//fmt.Println(value.Kind(), value.Elem())
		if ptrValue.Elem().IsValid() {
			value := reflect.Indirect(reflect.ValueOf(ptrValue.Interface()))
			if _, ok := value.Interface().(interface {
				Int()
			}); ok {
				value.Int()
			}
			if _, ok := value.Interface().(interface {
				String()
			}); ok {
				value.String()
			}
			if _, ok := value.Interface().(interface {
				Float()
			}); ok {
				value.Float()
			}
		}
	}

	//fmt.Printf("Reflecting value for field number: %d, value is: %v, interface value is: %v\n", i, f, f.Interface())
	//if f == reflect.Zero(reflect.TypeOf(f)) {
	//		fmt.Println(f.Interface())
	//	}
	/*if reflect.TypeOf(f.Interface()) == reflect.TypeOf(p) && f != reflect.Zero(reflect.TypeOf(p))  {
		fmt.Println("%v\n", reflect.Indirect(reflect.ValueOf(f.Interface())).Int())
	}
	if reflect.TypeOf(f.Interface()) == reflect.TypeOf(p1) && f != reflect.Zero(reflect.TypeOf(p)) {
		fmt.Printf("%v\n", reflect.Indirect(reflect.ValueOf(f.Interface())).Uint())
	}*/
}

func MagicVarFunc(action string) interface{} {
	return reflect.New(GeneratedDataStructure.Map[action]).Interface()
}
