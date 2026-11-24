package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/qlustra/intel"
	"github.com/qlustra/intel/capture"
)

type Counters struct {
	Hits   int `metric:"hits"`
	Misses int `metric:"misses"`
}

type Report struct {
	Counters
	InternalCount int // Not tagged; left alone by metric updates.
}

type OtherReport Report

type metricField struct {
	key  string
	path []int
}

type updatePlan struct {
	root   *capture.TypeInterrogator
	fields []metricField
}

// buildPlan defines this application's policy: a metric tag names an input
// key, and each key must identify one exported int property.
func buildPlan[T any]() (updatePlan, error) {
	root := intel.CaptureType[T]()
	shape, check := root.AsStruct()
	if !check.OK() {
		return updatePlan{}, fmt.Errorf("inspect report shape: %s", check.Why())
	}
	plan := updatePlan{root: root}
	seen := make(map[string]bool)
	for name, field := range shape.Properties.AllAddressable() {
		key, tagged := field.Tags().Lookup("metric")
		if !tagged {
			continue
		}
		if key == "" || seen[key] {
			return updatePlan{}, fmt.Errorf("%s: empty or duplicate metric key %q", name, key)
		}
		if !field.IsExported() || !field.Is[int]().OK() {
			return updatePlan{}, fmt.Errorf("%s: metric field must be an exported int", name)
		}
		seen[key] = true
		plan.fields = append(plan.fields, metricField{key: key, path: slices.Clone(field.Index())})
	}
	if len(plan.fields) == 0 {
		return updatePlan{}, fmt.Errorf("report shape has no metric fields")
	}
	return plan, nil
}

func (plan updatePlan) apply(report any, readings map[string]int) error {
	type write struct {
		address *int
		reading int
	}
	writes := make([]write, 0, len(plan.fields))
	for _, field := range plan.fields {
		reading, ok := readings[field.key]
		if !ok {
			return fmt.Errorf("missing reading for %q", field.key)
		}
		value, check := capture.ResolveField(report, plan.root, field.path)
		if !check.OK() {
			return fmt.Errorf("resolve %q: %s", field.key, check.Why())
		}
		address, check := value.Address()
		if !check.OK() {
			return fmt.Errorf("address for %q: %s", field.key, check.Why())
		}
		// The plan verified each field is int, and ResolveField checked the
		// exact struct type before following the saved path.
		writes = append(writes, write{address: address.(*int), reading: reading})
	}
	// This is the application's rule: resolve the whole batch before writing.
	for _, w := range writes {
		*w.address = w.reading
	}
	return nil
}

func main() {
	plan, err := buildPlan[Report]()
	if err != nil {
		log.Fatal(err)
	}
	for _, field := range plan.fields {
		fmt.Printf("plan: %s at %v\n", field.key, field.path)
	}

	first, second := Report{InternalCount: 99}, Report{InternalCount: 88}
	firstReadings := map[string]int{"hits": 12, "misses": 3}
	if err := plan.apply(&first, firstReadings); err != nil {
		log.Fatal(err)
	}
	if err := plan.apply(&second, map[string]int{"hits": 4, "misses": 1}); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("first: hits=%d misses=%d untouched=%d\n", first.Hits, first.Misses, first.InternalCount)
	fmt.Printf("second: hits=%d misses=%d untouched=%d\n", second.Hits, second.Misses, second.InternalCount)

	before := first
	err = plan.apply(&first, map[string]int{"hits": 100})
	if err == nil || first != before {
		log.Fatal("expected an incomplete update to be rejected without a write")
	}
	fmt.Println("incomplete update rejected; first unchanged")
	if err := plan.apply(&OtherReport{}, firstReadings); err == nil {
		log.Fatal("expected a different root type to be rejected")
	}
	fmt.Println("different root type rejected")
}
