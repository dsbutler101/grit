package tf

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type ModuleTypes struct {
	Name      string
	Provider  string
	Variables []*VariableType
	Outputs   []*OutputType
}

type VariableType struct {
	Name string
	Type string
}

type OutputType struct {
	Name      string
	Type      string
	Sensitive bool
}

type ModuleDependency struct {
	// The module with the dependency.
	DependantModule *ModuleTypes
	// The variable in that dependency which should be named after
	// the dependency module and typed as an object with fields
	// matching names and types of the dependencies outputs.
	DependantVariable *VariableType
	// A module that provides outputs.
	DependencyModule *ModuleTypes
	// Type string representing the module.
	DependencyOutputsType string
}

func LoadModuleDependencies(modulesDir string) ([]*ModuleDependency, error) {
	dirs, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil, err
	}
	modulesByProvider := map[string]map[string]*ModuleTypes{}
	for _, dir := range dirs {
		providerName := dir.Name()
		if providerName == "internal" {
			// We don't make any assertions about internal
			// because it isn't user facing.
			continue
		}
		if providerName == "gitlab" {
			// The GitLab module doesn't yet conform to
			// the provider/module structure.
			continue
		}
		provider, ok := modulesByProvider[providerName]
		if !ok {
			provider = map[string]*ModuleTypes{}
			modulesByProvider[providerName] = provider
		}
		dirs, err := os.ReadDir(filepath.Join(modulesDir, providerName))
		if err != nil {
			return nil, err
		}
		for _, dir := range dirs {
			moduleName := dir.Name()
			m, err := LoadModuleTypes(modulesDir, providerName, moduleName)
			if err != nil {
				return nil, err
			}
			provider[moduleName] = m
		}
	}
	dependencies := []*ModuleDependency{}
	for providerName, provider := range modulesByProvider {
		for _, dependency := range provider {
			for _, dependant := range modulesByProvider[providerName] {
				for _, variable := range dependant.Variables {
					if dependant.Name == dependency.Name {
						continue
					}
					if variable.Name == dependency.Name {
						outputTypes := []string{}
						for _, o := range dependency.Outputs {
							outputTypes = append(outputTypes, o.Name+":"+o.Type)
						}
						dependencyType := "object({" + strings.Join(outputTypes, ",") + "})"
						dependencies = append(dependencies, &ModuleDependency{
							DependantModule:       dependant,
							DependantVariable:     variable,
							DependencyModule:      dependency,
							DependencyOutputsType: dependencyType,
						})
					}
				}
			}
		}
	}
	return dependencies, nil
}

func LoadModuleTypes(modulesDir, providerName, moduleName string) (*ModuleTypes, error) {
	path := filepath.Join(modulesDir, providerName, moduleName)
	parser := hclparse.NewParser()
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	m := &ModuleTypes{
		Name:     moduleName,
		Provider: providerName,
	}
	for _, e := range dirEntry {
		if e.IsDir() {
			continue
		}
		filename := filepath.Join(path, e.Name())
		if !strings.HasSuffix(filename, ".tf") {
			continue
		}
		tfData, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		file, diags := parser.ParseHCL(tfData, filename)
		if diags != nil && diags.HasErrors() {
			return nil, err
		}
		hclSyntaxBody, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			return nil, fmt.Errorf("I can't work with this type: %T", file.Body)
		}
		for _, block := range hclSyntaxBody.Blocks {
			switch block.Type {
			case "variable":
				name := block.Labels[0]
				attrs := block.Body.Attributes
				v := &VariableType{
					Name: name,
				}
				typ, err := extractType(attrs)
				if err != nil {
					return nil, fmt.Errorf("reading %q from %v: %w (add a type hint?)", name, filename, err)
				}
				v.Type = typ
				m.Variables = append(m.Variables, v)
			case "output":
				name := block.Labels[0]
				attrs := block.Body.Attributes
				o := &OutputType{
					Name: name,
				}
				typ, err := extractType(attrs)
				if err != nil {
					return nil, fmt.Errorf("reading %q from %v: %w (add a type hint?)", name, filename, err)
				}
				o.Type = typ
				m.Outputs = append(m.Outputs, o)
			default:
				continue
			}
		}
	}
	// Normalize variables and outputs by sorting by name.
	slices.SortFunc(m.Variables, func(a, b *VariableType) int {
		return strings.Compare(a.Name, b.Name)
	})
	slices.SortFunc(m.Outputs, func(a, b *OutputType) int {
		return strings.Compare(a.Name, b.Name)
	})
	return m, nil
}

func extractType(attrs hclsyntax.Attributes) (string, error) {
	typ, have := attrs["type"]
	if have && typ != nil {
		return typeString(typ.Expr)
	}
	// Infer output types by reading function calls as type hints:3
	// https://github.com/hashicorp/terraform/issues/30579#issuecomment-1051407912
	val, have := attrs["value"]
	if !have || val == nil {
		return "", nil
	}
	fn, ok := val.Expr.(*hclsyntax.FunctionCallExpr)
	if !ok {
		return "", nil
	}
	return inferTypeFromFunctionCall(fn)
}

const (
	toNumberFn = "tonumber"
	numberType = "number"
	toStringFn = "tostring"
	stringType = "string"
	toBoolFn   = "tobool"
	boolType   = "bool"
	toListFn   = "tolist"
	listFn     = "list"
	listType   = "list"
	toMapFn    = "tomap"
	mapFn      = "map"
	mapType    = "map"
)

func inferTypeFromFunctionCall(fn *hclsyntax.FunctionCallExpr) (string, error) {
	switch fn.Name {
	case toNumberFn:
		return numberType, nil
	case toStringFn:
		return stringType, nil
	case toBoolFn:
		return boolType, nil
	case toListFn, listFn:
		return listType, nil
	case toMapFn, mapFn:
		return mapType, nil
	default:
		return "", fmt.Errorf("couldn't infer type of function call %v", fn.Name)
	}
}

func typeString(expr hcl.Expression) (string, error) {
	switch expr := expr.(type) {
	case *hclsyntax.FunctionCallExpr:
		switch expr.Name {
		case "optional":
			if len(expr.Args) == 0 {
				return "", fmt.Errorf("want at least one arg for optional. got %v", len(expr.Args))
			}
			// The default value doesn't affect the type.
			return typeString(expr.Args[0])
		case toNumberFn, toStringFn, toBoolFn, toListFn, listFn, toMapFn, mapFn:
			return inferTypeFromFunctionCall(expr)
		}
		args := []string{}
		for _, arg := range expr.Args {
			arg, err := typeString(arg)
			if err != nil {
				return "", err
			}
			args = append(args, arg)
		}
		return expr.Name + "(" + strings.Join(args, ",") + ")", nil
	case *hclsyntax.ObjectConsExpr:
		args := []string{}
		for _, i := range expr.Items {
			key, err := typeString(i.KeyExpr)
			if err != nil {
				return "", err
			}
			value, err := typeString(i.ValueExpr)
			if err != nil {
				return "", err
			}
			args = append(args, key+":"+value)
		}
		// Normalize object field order by sorting
		slices.Sort(args)
		return "{" + strings.Join(args, ",") + "}", nil
	case *hclsyntax.ObjectConsKeyExpr:
		return typeString(expr.Wrapped)
	case *hclsyntax.ScopeTraversalExpr:
		if len(expr.Traversal) != 1 {
			return "", fmt.Errorf("expecting 1 traversal. got %v", len(expr.Traversal))
		}
		switch trav := expr.Traversal[0].(type) {
		case hcl.TraverseRoot:
			return trav.Name, nil
		default:
			return "", fmt.Errorf("I can't work with this type: %T", expr.Traversal[0])
		}
	case *hclsyntax.TupleConsExpr:
		args := []string{}
		for _, arg := range expr.Exprs {
			arg, err := typeString(arg)
			if err != nil {
				return "", err
			}
			args = append(args, arg)
		}
		return "[" + strings.Join(args, ",") + "]", nil
	case *hclsyntax.TemplateExpr, *hclsyntax.LiteralValueExpr:
		return "<value>", nil
	default:
		return "", fmt.Errorf("I can't work with this type: %T", expr)
	}
}
