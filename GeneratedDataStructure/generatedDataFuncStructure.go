
package GeneratedDataStructure

import (
	
	"reflect"
	
	"xmltest/btalkTest/Helper"
	
)
// exported FuncMap
var FuncMap map[string] reflect.Value

func init() {
	FuncMap = make(map[string] reflect.Value)
	
	FuncMap["Helper_IncTimestamp"] = reflect.ValueOf(Helper.IncTimestamp)
	
	FuncMap["Helper_Md5For"] = reflect.ValueOf(Helper.Md5For)
	
	FuncMap["Helper_RequestId"] = reflect.ValueOf(Helper.RequestId)
	
	FuncMap["Helper_Timestamp"] = reflect.ValueOf(Helper.Timestamp)
	
}
