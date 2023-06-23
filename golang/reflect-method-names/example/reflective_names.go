package reflectnames

import (
	"reflect"
)

func GetMethodNames(arg any) []string {
	typeOfArg := reflect.TypeOf(arg)
	valueOfArg := reflect.ValueOf(arg)
	methodsNames := make([]string, valueOfArg.NumMethod())

	for i := 0; i < valueOfArg.NumMethod(); i++ {
		methodsNames[i] = typeOfArg.Method(i).Name
	}

	return methodsNames
}
