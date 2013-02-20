package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	PATH_SEP          = "\\/"
	DEFAULT_EXTENSION = ".func.go"
	SRC_DIR           = "src"
	templateCode      = `
package {{.GeneratedPackage}}

import (
	
	"reflect"
	{{range $path, $value:= .ImportPath}}
	"{{$path}}"
	{{end}}
)
// exported FuncMap
var FuncMap map[string] reflect.Value

func init() {
	FuncMap = make(map[string] reflect.Value)
	{{range $index, $dataType:= .Types}}
	FuncMap["{{$dataType.PackageName}}_{{$dataType.TypeName}}"] = reflect.ValueOf({{$dataType.PackageName}}.{{$dataType.TypeName}})
	{{end}}
}
`
)

type DataType struct {
	TypeName    string
	PackageName string
}

type FileDataType struct {
	ImportPath  string
	PackageName string
	Types       map[string]DataType
}

type OutputData struct {
	GeneratedPackage string
	ImportPath       map[string]bool
	Types            map[string]DataType
}

var PackageNameMap map[string]bool

func isPathSep(r rune) bool {
	for _, pS := range PATH_SEP {
		if r == pS {
			return true
		}
	}
	return false
}

// Adding the package name to import path if it does not contain
// normalize the import path according to go convention
func cleanImportPath(importPath string) string {
	importPath = strings.Replace(importPath, "\\", "/", -1)
	return strings.TrimFunc(importPath, isPathSep)
}

// parse Protobuf file and extract all data type from that file
// inputFile: input file to parse
// import path: import path to this file.
// output: list of all data type from input file, parsed into a FileDataType type
func parseProtobufFile(inputFilePath string, importPath string) FileDataType {
	fset := token.NewFileSet() // positions are relative to fset	
	// Parse this file
	parsedFile, err := parser.ParseFile(fset, inputFilePath, nil, 0)
	if err != nil {
		fmt.Println(err)
		return FileDataType{}
	}

	fmt.Println("----Parsing file: ", inputFilePath, " import path: ", path.Dir(inputFilePath))
	var data FileDataType
	fmt.Println(parsedFile.Name)
	data.PackageName = parsedFile.Name.Name
	data.ImportPath = cleanImportPath(importPath)

	fmt.Println("---Import Path: ", data.ImportPath)

	data.Types = make(map[string]DataType)

	for typeName, objType := range parsedFile.Scope.Objects {
		if objType.Kind == ast.Fun && (strings.ToUpper(typeName[:1]) == typeName[:1]){
			key := data.PackageName + "_" + typeName
			fmt.Println(key)
			newType := DataType{TypeName: typeName, PackageName: data.PackageName}
			data.Types[key] = newType
		}
	}
	return data
}

// visit a file path
// if the file is a .func.go (function helper file), call the parse function and add result to output
func visit(currentPath string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", currentPath)
	if !f.IsDir() {
		var fileName = f.Name()
		var strLen = len(fileName)
		if strLen >= len(DEFAULT_EXTENSION) && fileName[strLen-len(DEFAULT_EXTENSION):strLen] == DEFAULT_EXTENSION {
			absPath, err1 := filepath.Abs(currentPath)
			if err1 == nil {
				importPath := filepath.Dir(absPath)
				if importPath[0:len(goPath)] == goPath {
					importPath = importPath[len(goPath):len(importPath)]
				}
				srcIndex := strings.Index(importPath, SRC_DIR)
				importPath = importPath[srcIndex+len(SRC_DIR) : len(importPath)]
				fmt.Printf("file path from go path: %d %s\n", srcIndex, importPath)

				result := parseProtobufFile(absPath, importPath)

				if !output.ImportPath[result.ImportPath] {
					// add import path to output
					output.ImportPath[result.ImportPath] = true
				}
				for key, dataType := range result.Types {
					if _, ok := output.Types[key]; !ok {
						output.Types[key] = dataType
					}
				}
			}
		}
	}
	return err
}

var goPath string
var output OutputData

func main() {
	var generatedPackage string
	var inputDir string
	var outputFile string

	// get go path for later use
	goPath = os.Getenv("GOPATH")

	flag.StringVar(&inputDir, "inputDir", ".", "The inputDir to be parsed")
	flag.StringVar(&outputFile, "output", "GeneratedDataStructure/generatedDataFuncStructure.go", "The output file")
	flag.Parse()

	generatedPackage = path.Base(path.Dir(outputFile))
	fmt.Println("Input directory: ", inputDir)

	output.GeneratedPackage = generatedPackage
	output.ImportPath = make(map[string]bool)
	output.Types = make(map[string]DataType)

	fmt.Println("Start parsing now ...")
	fmt.Println("Parsing directory: ", inputDir, "")
	root := inputDir
	err := filepath.Walk(root, visit)
	fmt.Println("End of parsing")

	fmt.Println("Output file: ", outputFile)
	fmt.Println("Generated Package: ", generatedPackage)
	fmt.Println("Writing output file ...")
	t := template.Must(template.New("DataStructure").Parse(templateCode))

	err = t.Execute(os.Stdout, output)
	if err != nil {
		log.Println("Try Executing template, get error:", err)
	}

	fo, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	w := bufio.NewWriter(fo)

	err = t.Execute(w, output)
	if err != nil {
		log.Println("Try executing template and write to file, get error:", err)
	}
	if err = w.Flush(); err != nil {
		panic(err)
	}
}
