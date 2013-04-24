
package GeneratedDataStructure

import (
	
	"reflect"
	
	"xmltest/btalkTest/Helper"
	
)
// exported FuncMap
var FuncMap map[string] reflect.Value

func init() {
	FuncMap = make(map[string] reflect.Value)
	
	FuncMap["Helper_FileContent"] = reflect.ValueOf(Helper.FileContent)
	
	FuncMap["Helper_FileContentMultiPart"] = reflect.ValueOf(Helper.FileContentMultiPart)
	
	FuncMap["Helper_FileSize"] = reflect.ValueOf(Helper.FileSize)
	
	FuncMap["Helper_GetCurrentUserName"] = reflect.ValueOf(Helper.GetCurrentUserName)
	
	FuncMap["Helper_GetCurrentUserNameWithParam"] = reflect.ValueOf(Helper.GetCurrentUserNameWithParam)
	
	FuncMap["Helper_GetNextUniqueNumber"] = reflect.ValueOf(Helper.GetNextUniqueNumber)
	
	FuncMap["Helper_GetNextUserEmail"] = reflect.ValueOf(Helper.GetNextUserEmail)
	
	FuncMap["Helper_GetNextUserEmailWithParam"] = reflect.ValueOf(Helper.GetNextUserEmailWithParam)
	
	FuncMap["Helper_GetOAuthRaw"] = reflect.ValueOf(Helper.GetOAuthRaw)
	
	FuncMap["Helper_IncTimestamp"] = reflect.ValueOf(Helper.IncTimestamp)
	
	FuncMap["Helper_Md5For"] = reflect.ValueOf(Helper.Md5For)
	
	FuncMap["Helper_NumFilePart"] = reflect.ValueOf(Helper.NumFilePart)
	
	FuncMap["Helper_RequestId"] = reflect.ValueOf(Helper.RequestId)
	
	FuncMap["Helper_Timestamp"] = reflect.ValueOf(Helper.Timestamp)
	
}
