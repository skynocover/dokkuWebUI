package resp

import (
	"encoding/json"
	"reflect"
)

// Error ...
type Error struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Err          string `json:"err"`
}

// ToBytes ...
func (Err *Error) ToBytes() []byte {
	jsondata, _ := json.Marshal(Err)
	return jsondata
}

// ToBytesWithErr ...
func (Err *Error) ToBytesWithErr(err error) []byte {
	Err.Err = err.Error()
	jsondata, _ := json.Marshal(Err)
	return jsondata
}

// ToBytesWithStruct ...
func (Err *Error) ToBytesWithStruct(obj interface{}) []byte {
	m := map[string]interface{}{}
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	for i := 0; i < obj1.NumField(); i++ {
		m[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}

	json.Unmarshal(Err.ToBytes(), &m)
	for k, v := range m {
		m[k] = v
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}

// ToBytesWithObject ...
func (Err *Error) ToBytesWithObject(v map[string]interface{}) []byte {
	m := map[string]interface{}{}
	json.Unmarshal(Err.ToBytes(), &m)
	for k, v := range v {
		m[k] = v
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}

// MustWithObject ...
func (Err *Error) MustWithObject(v map[string]interface{}) map[string]interface{} {
	m := map[string]interface{}{}

	json.Unmarshal(Err.ToBytes(), &m)
	for k, v := range v {
		m[k] = v
	}

	return m
}

// SUCCESS ...
var SUCCESS = Error{
	ErrorCode:    0,
	ErrorMessage: "success",
}

var LoginFail = Error{
	ErrorCode:    1,
	ErrorMessage: "Account or password incorrect",
}

var SessionExpired = Error{
	ErrorCode:    2,
	ErrorMessage: "Session expired",
}

var PrivateKeyNotFound = Error{
	ErrorCode:    3,
	ErrorMessage: "Private key not found",
}

var SystemError = Error{
	ErrorCode:    9999,
	ErrorMessage: "SystemError",
}

//////////// api fail

// ErrorParsingJSON ...
var ErrorParsingJSON = Error{
	ErrorCode:    1000,
	ErrorMessage: "json format error",
}

// ErrorParameter ...
var ErrorParameter = Error{
	ErrorCode:    1001,
	ErrorMessage: "parameter input fail",
}

// ErrorAppNotExist ...
var ErrorAppNotExist = Error{
	ErrorCode:    1003,
	ErrorMessage: "App not found",
}

// ErrorAppAlreadyExist ...
var ErrorAppAlreadyExist = Error{
	ErrorCode:    1004,
	ErrorMessage: "App already exists",
}

// ErrorSSHKeyAlreadyExist ...
var ErrorSSHKeyAlreadyExist = Error{
	ErrorCode:    1005,
	ErrorMessage: "SSH key already exists",
}

///////// system fail
var ErrorCloudFlareFail = Error{
	ErrorCode:    1006,
	ErrorMessage: "Cloud flare fail",
}
