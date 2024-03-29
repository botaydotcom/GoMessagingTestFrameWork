package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
	//"code.google.com/p/goprotobuf/proto"
	"xmltest/btalkTest/GeneratedDataStructure"
)

type LoginUserInfo_LastLoginInfo struct {
	IP      *int32  `protobuf:"varint,1,opt" json:"IP,omitempty"`
	Time    *int32  `protobuf:"varint,2,opt" json:"Time,omitempty"`
	Country *string `protobuf:"bytes,3,opt" json:"Country,omitempty"`
}

type LoginUserInfo struct {
	LastLogin  []LoginUserInfo_LastLoginInfo `protobuf:"bytes,2,opt" json:"LastLogin,omitempty"`
	ServerTime *int32                        `protobuf:"varint,3,opt" json:"ServerTime,omitempty"`
	TestByte   []byte                        `protobuf:"bytes,4,req" json:"InviteId,omitempty"`
}

func visitGetNameFieldPtr(w reflect.Value, fieldName string) {
	domainName := fieldName
	fmt.Printf("field name: %s\n", domainName)
	fv := reflect.Indirect(w)

	if fv.Kind() == reflect.Struct {
		visitGetNameFieldStruct(fv, domainName)
	}
	if fv.Kind() == reflect.Slice {
		visitGetNameFieldSlice(fv, domainName)
	}
	if fv.Kind() == reflect.Ptr {
		visitGetNameFieldPtr(fv, domainName)
	}
}

func visitGetNameFieldSlice(w reflect.Value, fieldName string) {
	for i := 0; i < w.Len(); i++ {
		fv := w.Index(i)
		domainName := fmt.Sprintf("%s[%d]", fieldName, i)
		fmt.Printf("field name: %s\n", domainName)
		if fv.Kind() == reflect.Struct {
			visitGetNameFieldStruct(fv, domainName)
		}
		if fv.Kind() == reflect.Slice {
			visitGetNameFieldSlice(fv, domainName)
		}
		if fv.Kind() == reflect.Ptr {
			visitGetNameFieldPtr(fv, domainName)
		}
	}
}

func visitGetNameFieldStruct(w reflect.Value, fieldName string) {
	for i := 0; i < w.NumField(); i++ {
		f := w.Type().Field(i)
		domainName := fmt.Sprintf("%s.%s", fieldName, f.Name)
		fmt.Printf("field name: %s\n", domainName)

		fv := w.Field(i)
		if fv.Kind() == reflect.Struct {
			visitGetNameFieldStruct(fv, domainName)
		}
		if fv.Kind() == reflect.Slice {
			visitGetNameFieldSlice(fv, domainName)
		}
		if fv.Kind() == reflect.Ptr {
			visitGetNameFieldPtr(fv, domainName)
		}
	}
}

func main() {
	c := "VN"
	ll := LoginUserInfo_LastLoginInfo{Country: &c}
	lla := []LoginUserInfo_LastLoginInfo{ll}
	v := &LoginUserInfo{LastLogin: lla, TestByte: []byte{70, 71, 72}}

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)
	var w LoginUserInfo
	xml.Unmarshal(output, &w)
	//fmt.Printf("\n%v %v %v", *w.LastLogin, *w.LastLogin.Country, w.TestByte)

	fmt.Println()
	visitGetNameFieldStruct(reflect.ValueOf(w), "")
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
