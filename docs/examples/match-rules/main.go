package main

import (
	"fmt"

	"github.com/qlustra/intel"
	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules"
	"github.com/qlustra/intel/rules/field"
	"github.com/qlustra/intel/rules/property"
)

type Counters struct{ Hits int }
type Report struct {
	Counters
	Title string `json:"title"`
	Note  bool
}

func main() {
	namedTitle := property.Property()
	namedTitle.Expect("named Title", func(e property.Encounter) bool {
		return e.Property.FieldName() == "Title"
	})
	stringType := property.Property()
	stringType.ExpectCheck("string type", func(e property.Encounter) logic.Check {
		return e.Property.Is[string]()
	})
	var checks property.Rules
	checks.Match(options.Optional, namedTitle)
	checks.Expect(stringType, property.Field(field.TaggedWith("json")))
	for _, candidate := range []struct {
		name string
		root *capture.TypeInterrogator
	}{
		{"tagged title and unrelated fields", intel.CaptureType[Report]()},
		{"no title", intel.CaptureType[struct{ Note bool }]()},
		{"non-string title", intel.CaptureType[struct {
			Title int `json:"title"`
		}]()},
		{"untagged title", intel.CaptureType[struct{ Title string }]()},
	} {
		shape, check := candidate.root.AsStruct()
		if !check.OK() {
			panic(check.Why())
		}
		accepted, valid := checkTitlePolicy(shape, &checks)
		fmt.Printf("%s: accepted=%t valid=%t\n", candidate.name, accepted, valid)
	}
}

func checkTitlePolicy(shape *capture.StructInterrogator, checks *property.Rules) (accepted, valid bool) {
	validation := rules.NewValidation("title policy")
	for _, member := range shape.Properties.AllAddressable() {
		satisfied, valid := checks.ValidateRules(validation.Branch(member.FieldName()), property.Encounter{
			Struct: shape, Property: member,
		})
		if !satisfied || !valid {
			return false, valid
		}
	}
	return true, true
}
