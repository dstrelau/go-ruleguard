// +build generate

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"
)

type opInfo struct {
	name      string
	comment   string
	valueType string
	flags     uint64
}

const (
	flagIsBinaryExpr uint64 = 1 << iota
	flagIsBasicLit
)

func main() {
	ops := []opInfo{
		{name: "Invalid"},

		{name: "Not", comment: "!$Args[0]"},

		// Binary expressions.
		{name: "And", comment: "$Args[0] && $Args[1]", flags: flagIsBinaryExpr},
		{name: "Or", comment: "$Args[0] || $Args[1]", flags: flagIsBinaryExpr},
		{name: "Eq", comment: "$Args[0] == $Args[1]", flags: flagIsBinaryExpr},
		{name: "Neq", comment: "$Args[0] != $Args[1]", flags: flagIsBinaryExpr},
		{name: "Gt", comment: "$Args[0] > $Args[1]", flags: flagIsBinaryExpr},
		{name: "Lt", comment: "$Args[0] < $Args[1]", flags: flagIsBinaryExpr},
		{name: "GtEq", comment: "$Args[0] >= $Args[1]", flags: flagIsBinaryExpr},
		{name: "LtEq", comment: "$Args[0] <= $Args[1]", flags: flagIsBinaryExpr},

		{name: "VarAddressable", comment: "m[$Value].Addressable", valueType: "string"},
		{name: "VarPure", comment: "m[$Value].Pure", valueType: "string"},
		{name: "VarConst", comment: "m[$Value].Const", valueType: "string"},
		{name: "VarText", comment: "m[$Value].Text", valueType: "string"},
		{name: "VarLine", comment: "m[$Value].Line", valueType: "string"},
		{name: "VarValueInt", comment: "m[$Value].Value.Int()", valueType: "string"},
		{name: "VarTypeSize", comment: "m[$Value].Type.Size", valueType: "string"},

		{name: "VarFilter", comment: "m[$Value].Filter($Args[0])", valueType: "string"},
		{name: "VarNodeIs", comment: "m[$Value].Node.Is($Args[0])", valueType: "string"},
		{name: "VarObjectIs", comment: "m[$Value].Object.Is($Args[0])", valueType: "string"},
		{name: "VarTypeIs", comment: "m[$Value].Type.Is($Args[0])", valueType: "string"},
		{name: "VarTypeUnderlyingIs", comment: "m[$Value].Type.Underlying().Is($Args[0])", valueType: "string"},
		{name: "VarTypeConvertibleTo", comment: "m[$Value].Type.ConvertibleTo($Args[0])", valueType: "string"},
		{name: "VarTypeAssignableTo", comment: "m[$Value].Type.AssignableTo($Args[0])", valueType: "string"},
		{name: "VarTypeImplements", comment: "m[$Value].Type.Implements($Args[0])", valueType: "string"},
		{name: "VarTextMatches", comment: "m[$Value].Text.Matches($Args[0])", valueType: "string"},

		{name: "FileImports", comment: "m.File.Imports($Value)", valueType: "string"},
		{name: "FilePkgPathMatches", comment: "m.File.PkgPath.Matches($Value)", valueType: "string"},
		{name: "FileNameMatches", comment: "m.File.Name.Matches($Value)", valueType: "string"},

		{name: "FilterFuncRef", comment: "$Value holds a function name", valueType: "string"},

		{name: "String", comment: "$Value holds a string constant", valueType: "string", flags: flagIsBasicLit},
		{name: "Int", comment: "$Value holds an int64 constant", valueType: "int64", flags: flagIsBasicLit},
	}

	var buf bytes.Buffer

	buf.WriteString(`// Code generated "gen_filter_op.go"; DO NOT EDIT.` + "\n")
	buf.WriteString("\n")
	buf.WriteString("package ir\n")
	buf.WriteString("const (\n")

	for i, op := range ops {
		if strings.Contains(op.comment, "$Value") && op.valueType == "" {
			fmt.Printf("missing %s valueType\n", op.name)
		}
		if op.comment != "" {
			buf.WriteString("// " + op.comment + "\n")
		}
		if op.valueType != "" {
			buf.WriteString("// $Value type: " + op.valueType + "\n")
		}
		fmt.Fprintf(&buf, "Filter%sOp FilterOp = %d\n", op.name, i)
		buf.WriteString("\n")
	}
	buf.WriteString(")\n")

	buf.WriteString("var filterOpNames = map[FilterOp]string{\n")
	for _, op := range ops {
		fmt.Fprintf(&buf, "Filter%sOp: `%s`,\n", op.name, op.name)
	}
	buf.WriteString("}\n")

	buf.WriteString("var filterOpFlags = map[FilterOp]uint64{\n")
	for _, op := range ops {
		if op.flags == 0 {
			continue
		}
		parts := make([]string, 0, 1)
		if op.flags&flagIsBinaryExpr != 0 {
			parts = append(parts, "flagIsBinaryExpr")
		}
		if op.flags&flagIsBasicLit != 0 {
			parts = append(parts, "flagIsBasicLit")
		}
		fmt.Fprintf(&buf, "Filter%sOp: %s,\n", op.name, strings.Join(parts, " | "))
	}
	buf.WriteString("}\n")

	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("filter_op.gen.go", pretty, 0644); err != nil {
		panic(err)
	}
}
