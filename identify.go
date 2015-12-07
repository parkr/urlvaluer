// Copyright 2014 Brett Slatkin, Parker Moore
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func loadFile(inputPath string) (string, []GeneratedType) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, inputPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Could not parse file: %s", err)
	}

	packageName := identifyPackage(f)
	if packageName == "" {
		log.Fatalf("Could not determine package name of %s", inputPath)
	}

	valuertypes := map[string]*StructSpec{}
	valuers := map[string]bool{}
	for _, decl := range f.Decls {
		structSpec, ok := identifyUrlValuerType(decl)
		if ok {
			valuertypes[structSpec.Name] = structSpec
			continue
		}

		typeName, ok := identifyUrlValuer(decl)
		if ok {
			valuers[typeName] = true
			continue
		}
	}

	types := []GeneratedType{}
	for typeName, typeSpec := range valuertypes {
		_, isUrlValuer := valuers[typeName]
		urlvaluer := GeneratedType{*typeSpec, isUrlValuer}
		types = append(types, urlvaluer)
	}

	return packageName, types
}

func identifyPackage(f *ast.File) string {
	if f.Name == nil {
		return ""
	}
	return f.Name.Name
}

func identifyUrlValuerType(decl ast.Decl) (structSpec *StructSpec, match bool) {
	genDecl, ok := decl.(*ast.GenDecl)
	if !ok {
		return
	}

	structSpec = &StructSpec{}

	for _, spec := range genDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				if typeSpec.Name != nil {
					structSpec.Name = typeSpec.Name.Name
				}
				if structType.Fields != nil {
					structSpec.Fields = getFieldData(structType.Fields.List)
				}
				break
			}
		}
	}
	if structSpec.Name == "" || len(structSpec.Fields) == 0 {
		return
	}

	match = true
	return
}

func identifyUrlValuer(decl ast.Decl) (typeName string, match bool) {
	funcDecl, ok := decl.(*ast.FuncDecl)
	if !ok {
		return
	}

	// Method name should match fmt.Stringer
	if funcDecl.Name == nil {
		return
	}
	if funcDecl.Name.Name != "UrlValues" {
		return
	}

	// Should have no arguments
	if funcDecl.Type == nil {
		return
	}
	if funcDecl.Type.Params == nil {
		return
	}
	if len(funcDecl.Type.Params.List) != 0 {
		return
	}

	// Return value should be a string
	if funcDecl.Type.Results == nil {
		return
	}
	if len(funcDecl.Type.Results.List) != 1 {
		return
	}
	result := funcDecl.Type.Results.List[0]
	if result.Type == nil {
		return
	}
	if ident, ok := result.Type.(*ast.Ident); !ok {
		return
	} else if ident.Name != "string" {
		return
	}

	// Receiver type
	if funcDecl.Recv == nil {
		return
	}
	if len(funcDecl.Recv.List) != 1 {
		return
	}
	recv := funcDecl.Recv.List[0]
	if recv.Type == nil {
		return
	}
	if ident, ok := recv.Type.(*ast.Ident); !ok {
		return
	} else {
		typeName = ident.Name
	}

	match = true
	return
}

func getFieldData(fieldsList []*ast.Field) []fieldSpec {
	fields := []fieldSpec{}
	for _, field := range fieldsList {
		if !field.Names[0].IsExported() || field.Names[0].Name == "XXX_unrecognized" {
			continue
		}

		var snakeCaseName string
		tag, err := strconv.Unquote(field.Tag.Value)
		if err == nil {
			st := reflect.StructTag(tag)
			snakeCaseName = strings.SplitN(st.Get("json"), ",", 2)[0]
		}

		fields = append(fields, NewFieldSpec(
			field.Names[0].Name,
			snakeCaseName,
			field.Type,
		))
	}
	return fields
}
