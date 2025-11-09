package model

import (
	"go/ast"
	"reflect"
	"testing"

	"golang.org/x/tools/go/packages"
	"gorm.io/gorm"
)

// TestAllGormModelsRegistered scans the model package for all structs that
// embed gorm.Model and ensures they are all included in the Models variable.
func TestAllGormModelsRegistered(t *testing.T) {
	// Credit for this code goes to Claude Sonnet 4.5.

	// Step 1: Load the current package using go/packages
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedName,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		t.Fatalf("Failed to load package: %v", err)
	}
	if len(pkgs) == 0 {
		t.Fatal("No packages found")
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Fatal("Package has errors")
	}

	pkg := pkgs[0]

	// Step 2: Collect all struct names that embed gorm.Model
	gormModelStructs := make(map[string]bool)
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			if typeSpec, ok := n.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					if embedsGormModel(structType) {
						gormModelStructs[typeSpec.Name.Name] = true
					}
				}
			}
			return true
		})
	}

	if len(gormModelStructs) == 0 {
		t.Log("No structs embedding gorm.Model found - this test may need adjustment")
	}

	// Step 3: Check which structs are registered in Models
	registeredModels := make(map[string]bool)
	for _, model := range Models {
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Pointer {
			modelType = modelType.Elem()
		}
		registeredModels[modelType.Name()] = true
	}

	// Step 4: Find missing registrations
	var missingModels []string
	for structName := range gormModelStructs {
		if !registeredModels[structName] {
			missingModels = append(missingModels, structName)
		}
	}

	if len(missingModels) > 0 {
		t.Errorf("The following structs embed gorm.Model but are not registered in Models: %v", missingModels)
	}

	// Step 5: Verify that all registered models actually embed gorm.Model
	var extraModels []string
	for modelName := range registeredModels {
		if !gormModelStructs[modelName] {
			extraModels = append(extraModels, modelName)
		}
	}

	if len(extraModels) > 0 {
		t.Errorf("The following models are registered in Models but don't embed gorm.Model: %v", extraModels)
	}
}

// embedsGormModel checks if a struct type embeds gorm.Model
func embedsGormModel(structType *ast.StructType) bool {
	for _, field := range structType.Fields.List {
		// Embedded fields have no names
		if len(field.Names) == 0 {
			// Check if the embedded type is a selector expression (e.g., gorm.Model)
			if selectorExpr, ok := field.Type.(*ast.SelectorExpr); ok {
				// Check if it's gorm.Model
				if ident, ok := selectorExpr.X.(*ast.Ident); ok {
					if ident.Name == "gorm" && selectorExpr.Sel.Name == "Model" {
						return true
					}
				}
			}
		}
	}
	return false
}

// TestModelsContainValidGormModels verifies that all items in Models
// actually embed gorm.Model at runtime using reflection
func TestModelsContainValidGormModels(t *testing.T) {
	gormModelType := reflect.TypeOf(gorm.Model{})

	for i, model := range Models {
		modelType := reflect.TypeOf(model)
		if modelType.Kind() == reflect.Ptr {
			modelType = modelType.Elem()
		}

		if modelType.Kind() != reflect.Struct {
			t.Errorf("Models[%d] (%v) is not a struct", i, modelType)
			continue
		}

		// Check if the struct has an embedded gorm.Model field
		hasGormModel := false
		for j := 0; j < modelType.NumField(); j++ {
			field := modelType.Field(j)
			if field.Anonymous && field.Type == gormModelType {
				hasGormModel = true
				break
			}
		}

		if !hasGormModel {
			t.Errorf("Models[%d] (%s) does not embed gorm.Model", i, modelType.Name())
		}
	}
}
