package main

import (
	"fmt"

	"github.com/qlustra/intel/compile"
	"github.com/qlustra/intel/define"
	"github.com/qlustra/intel/logic"
	"github.com/qlustra/intel/options"
	"github.com/qlustra/intel/rules/property"
)

type Report struct {
	Hits int `intake:"hits"`
}
type Missing struct{}

type declaration struct {
	define.ScopeDeclaration
	Field string
}

type definition struct {
	define.StructDefinition[declaration]
}

func (d *definition) Define() {
	hits := property.Property()
	hits.Expect("property tagged for intake", func(e property.Encounter) bool {
		_, ok := e.Property.Tags().Lookup("intake")
		return ok
	})
	d.Property(options.Required, hits).Encounter(func(e property.Encounter) logic.Check {
		decl, check := e.Access[declaration](define.DeclarationKey)
		if !check.OK() {
			return check
		}
		decl.Field, _ = e.Property.Tags().Lookup("intake")
		return logic.Ok()
	})
}

func main() {
	d, err := define.Struct[definition, declaration](nil)
	if err != nil {
		panic(err)
	}
	good := compile.Struct[Report](d)
	if !good.IsOkay() {
		panic(good.Trace())
	}
	decl, check := good.GetDeclaredValue()
	if !check.OK() {
		panic(check.Why())
	}
	// The service can reuse the declaration to publish its intake inventory.
	fmt.Println("Intake inventory:", decl.Field)

	bad := compile.Struct[Missing](d)
	fmt.Println("Missing compiled:", bad.IsOkay())
	if !bad.IsOkay() {
		fmt.Println("Missing trace:")
		fmt.Println(bad.Trace())
	}
	_, check = bad.GetDeclaredValue()
	fmt.Println("Missing declaration available:", check.OK())
}
