package main

import (
	"fmt"

	"github.com/qlustra/intel"
	"github.com/qlustra/intel/logic"
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
	shape, check := intel.CaptureType[Report]().AsStruct()
	if !check.OK() {
		panic(check.Why())
	}
	for _, tc := range []struct {
		name, member string
		nested       field.RuleAPI
	}{
		{"match", "Title", field.TaggedWith("json")},
		{"non-match", "Hits", field.TaggedWith("json")},
		{"failed evaluation", "Title", field.Field()},
	} {
		member, check := shape.Property(tc.member)
		if !check.OK() {
			panic(check.Why())
		}
		outer := property.Property()
		outer.ExpectMatch("custom field check", func(branch rules.ValidationBranch, e property.Encounter) (logic.Check, bool) {
			matched, valid := tc.nested.ValidateRule(branch.Branch("nested field checks"), field.Encounter{
				EncounterContext: e.EncounterContext, Struct: e.Struct, Field: e.Property.FieldDescriptor,
			})
			if !valid {
				return logic.Nope("nested_failed", "nested evaluation failed"), false
			}
			if !matched {
				return logic.Nope("nested_unmatched", "nested rule did not match"), true
			}
			return logic.Ok(), true
		})
		validation := rules.NewValidation("property validation")
		matched, valid := outer.ValidateRule(validation, property.Encounter{Struct: shape, Property: member})
		fmt.Printf("%s: matched=%t valid=%t\n", tc.name, matched, valid)
		if !valid {
			fmt.Println(validation.Trace())
		}
	}
}
