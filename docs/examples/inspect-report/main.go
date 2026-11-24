package main

import (
	"fmt"
	"reflect"

	"github.com/qlustra/intel/inspect"
	"github.com/qlustra/intel/inspect/tag/formats"
	"github.com/qlustra/intel/options"
)

type Counters struct {
	Hits int `json:"hits,omitempty" label:""`
}

type Report struct {
	Counters
	Title string `json:"title"`
}

type OtherCounters struct{ Hits int }

type AmbiguousReport struct {
	Counters
	OtherCounters
}

func main() {
	typeInspector := inspect.CaptureType(reflect.TypeFor[Report]())
	fields := inspect.InspectFields(typeInspector, inspect.NewDescriptorRegistry(), nil, nil)
	fmt.Println("Report fields:")
	for name, field := range fields.AllAddressable(options.BreadthFirst) {
		var jsonValue string
		var jsonOptions []string
		if _, present := field.Tags().Lookup("json"); present {
			jsonTag, check := field.Tags().Parse[formats.CSV]("json")
			if !check.OK() {
				panic(check.Why())
			}
			parts := jsonTag.Format.Parts
			if len(parts) > 0 {
				jsonValue = parts[0]
				jsonOptions = parts[1:]
			}
		}
		label, labelPresent := field.Tags().Lookup("label")
		fmt.Printf("  %s path=%v type=%s json-key=json value=%s options=%v label=%q present=%t\n",
			name, field.Index(), field.ReflectedType(), jsonValue, jsonOptions, label, labelPresent)
	}

	ambiguousType := inspect.CaptureType(reflect.TypeFor[AmbiguousReport]())
	ambiguous := inspect.InspectFields(ambiguousType, inspect.NewDescriptorRegistry(), nil, nil)
	if _, check := ambiguous.LookupNamed("Hits"); check.OK() {
		fmt.Println("AmbiguousReport Hits: selected")
	} else {
		fmt.Println("AmbiguousReport Hits: ambiguous or absent")
	}
}
