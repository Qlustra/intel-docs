package main

import (
	"fmt"

	"github.com/qlustra/intel/capture"
	"github.com/qlustra/intel/compile"
	"github.com/qlustra/intel/define"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules/property"
)

type Counters struct{ Hits int }
type Report struct {
	First  Counters `label:"first"`
	Second Counters `label:"second"`
}

type body struct {
	define.ScopeDeclaration
	Field string
}
type reportDeclaration struct {
	define.ScopeDeclaration
	Uses []*define.Carrier[body]
}

type countersDefinition struct{ define.StructDefinition[body] }

func (d *countersDefinition) Define() {
	hits := property.Property()
	hits.Expect("field named Hits", func(e property.Encounter) bool { return e.Property.FieldName() == "Hits" })
	d.Property(options.Required, hits).Encounter(func(e property.Encounter) logic.Check {
		decl, check := e.Access[body](define.DeclarationKey)
		if !check.OK() {
			return check
		}
		decl.Field = e.Property.FieldName()
		return logic.Ok()
	})
}

type reportDefinition struct {
	define.StructDefinition[reportDeclaration]
}

func (d *reportDefinition) Define() {
	fields := property.Property()
	fields.Expect("counter fields", func(e property.Encounter) bool {
		return e.Property.FieldName() == "First" || e.Property.FieldName() == "Second"
	})
	scope := d.Scope[body, countersDefinition](fields)
	// Select the same body type for both original field uses.
	scope.Target(func(property.Encounter) (*capture.StructInterrogator, logic.Check) {
		return capture.Type[Counters]().AsStruct()
	})
	scope.Encounter(func(_ property.Encounter, host *reportDeclaration, use *define.Carrier[body]) logic.Check {
		host.Uses = append(host.Uses, use)
		return logic.Ok()
	})
}

func main() {
	d, err := define.Struct[reportDefinition, reportDeclaration](nil)
	if err != nil {
		panic(err)
	}
	c := compile.Struct[Report](d)
	if !c.IsOkay() {
		panic(c.Trace())
	}
	decl, check := c.GetDeclaredValue()
	if !check.OK() {
		panic(check.Why())
	}
	if len(decl.Uses) != 2 || decl.Uses[0] == decl.Uses[1] || decl.Uses[0].Declaration == nil || decl.Uses[0].Declaration != decl.Uses[1].Declaration {
		panic("expected two uses of one body")
	}
	for _, use := range decl.Uses {
		label, ok := use.Property.Tags().Lookup("label")
		if !ok {
			panic("missing label")
		}
		fmt.Printf("%s.%s label=%s\n", use.Property.FieldName(), use.Declaration.Field, label)
	}
}
