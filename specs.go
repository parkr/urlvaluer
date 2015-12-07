// Copyright 2014 Parker Moore
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
	"fmt"
	"go/ast"
	"log"
	"strings"
)

type StructSpec struct {
	Name   string
	Fields []fieldSpec
}

type fieldSpec struct {
	Name          string
	SnakeCaseName string
	resolvedType  string
	isPointer     bool
}

func NewFieldSpec(name, snakeCase string, astType ast.Expr) fieldSpec {
	spec := fieldSpec{Name: name, SnakeCaseName: snakeCase}
	return *spec.withType(astType)
}

func (f *fieldSpec) withType(astType ast.Expr) *fieldSpec {
	var fieldType ast.Expr
	fieldType = astType
	if starExpr, ok := fieldType.(*ast.StarExpr); ok {
		fieldType = starExpr.X
		f.isPointer = true
	}
	if ident, ok := fieldType.(*ast.Ident); ok {
		f.resolvedType = ident.String()
	} else {
		f.resolvedType = fmt.Sprintf("%T", fieldType)
	}
	return f
}

func (f fieldSpec) Accessor(owner string) string {
	if f.isPointer {
		return fmt.Sprintf("*%s.%s", owner, f.Name)
	} else {
		return fmt.Sprintf("%s.%s", owner, f.Name)
	}
}

func (f fieldSpec) Zero() string {
	log.Printf("%v is a %v and isPointer=%t", f.Name, f.resolvedType, f.isPointer)
	if f.isPointer {
		return `nil`
	}

	switch f.resolvedType {
	case "string":
		return `""`
	case "int", "int64", "float", "float64":
		return `0`
	case "bool":
		return `false`
	default:
		return `0`
	}
}

func (f fieldSpec) IsString() bool {
	return f.resolvedType == "string"
}

func (f fieldSpec) IsInt() bool {
	return strings.HasPrefix(f.resolvedType, "int")
}

func (f fieldSpec) IsFloat() bool {
	return strings.HasPrefix(f.resolvedType, "float")
}

func (f fieldSpec) IsStruct() bool {
	return !(f.IsString() || f.IsInt() || f.IsFloat())
}
