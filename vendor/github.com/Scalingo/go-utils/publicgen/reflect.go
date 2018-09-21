package publicgen

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var typeStruct = map[string]string{
	"bson.ObjectId": "string",
}

func listFields(st reflect.Type) ([]field, []reflect.Type) {
	fields := make([]field, 0)
	imported := make([]reflect.Type, 0)

	for i := 0; i < st.NumField(); i++ {
		curField := st.Field(i)

		if curField.Anonymous {
			newFields, newTypes := listFields(st.Field(i).Type)
			fields = append(fields, newFields...)
			imported = append(imported, newTypes...)
			continue
		}

		name := curField.Name
		if unicode.IsLower([]rune(name)[0]) {
			continue
		}

		json := curField.Tag.Get("json")
		slice := false
		typ := curField.Type

		if curField.Type.Kind() == reflect.Slice {
			slice = true
			typ = curField.Type.Elem()
		}

		pointer := false
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			pointer = true
		}

		pkg := typ.PkgPath()
		typePrefix := ""

		if pkg != "" {
			path := strings.Split(pkg, "/")
			typePrefix = path[len(path)-1] + "."
		}

		if strings.HasPrefix(pkg, st.PkgPath()) {
			imported = append(imported, typ)
			pkg = ""
			typePrefix = ""
		}
		typStr := typ.Name()
		fullType := fmt.Sprintf("%s%s", typePrefix, typStr)

		if newName, ok := typeStruct[fullType]; ok {
			fullType = newName
			typePrefix = ""
			pkg = ""
		}

		fields = append(fields, field{
			Name:    name,
			Type:    fullType,
			JSONTag: json,
			PkgPath: pkg,
			Pointer: pointer,
			Slice:   slice,
		})
	}
	return fields, imported
}
