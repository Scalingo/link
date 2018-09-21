package publicgen

import (
	"reflect"
)

type GeneratorParams struct {
	OutputFile    string
	OutputPackage string
	Types         []interface{}
}

func GeneratePublicModels(params GeneratorParams) error {
	astFields := make(map[string][]field)
	importedTypes := make([]reflect.Type, 0)
	alreadyDone := make([]string, 0)

	for _, typ := range params.Types {
		st := reflect.TypeOf(typ)
		fields, types := listFields(st)
		alreadyDone = append(alreadyDone, st.Name())
		importedTypes = append(importedTypes, types...)
		astFields[st.Name()] = fields
	}

	containsImported := true
	for containsImported {
		containsImported = false
		newImportedTypes := make([]reflect.Type, 0)
		for _, typ := range importedTypes {
			if !contains(alreadyDone, typ.Name()) {
				containsImported = true
				fields, types := listFields(typ)

				alreadyDone = append(alreadyDone, typ.Name())
				astFields[typ.Name()] = fields
				newImportedTypes = append(newImportedTypes, types...)
			}
		}
		importedTypes = newImportedTypes
	}

	ast := newAST(params.OutputPackage, astFields)
	err := writeAst(ast, params.OutputFile)
	return err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
