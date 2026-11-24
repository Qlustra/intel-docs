package main

import (
	"fmt"

	"github.com/qlustra/intel"
	"github.com/qlustra/intel/logic"
)

type Counter struct{ Hits int }

func (c *Counter) Read() int      { return c.Hits }
func (c *Counter) Increment() int { c.Hits++; return c.Hits }

type incrementer interface{ Increment() int }
type stringIncrementer interface{ Increment() string }

func must(check logic.Check) {
	if !check.OK() {
		panic(check.Why())
	}
}

func main() {
	receiver, check := intel.CaptureType[Counter]().ZeroReceiver()
	must(check)
	method, check := receiver.Method("Read")
	must(check)
	getter, check := receiver.Getter(method)
	must(check)
	result, check := getter.Read()
	must(check)
	fmt.Println("zero receiver:", result.Get())

	counter := Counter{Hits: 4}
	value, check := intel.CaptureValue(&counter)
	must(check)
	structure, check := value.AsStruct()
	must(check)
	hits, check := structure.Invoke[incrementer, int]("Increment", func(c incrementer) (int, logic.Check) {
		return c.Increment(), logic.Ok()
	})
	must(check)
	fmt.Printf("invoked: %d, live hits: %d\n", hits, counter.Hits)

	_, check = structure.Invoke[stringIncrementer, string]("Increment", func(c stringIncrementer) (string, logic.Check) {
		return c.Increment(), logic.Ok()
	})
	fmt.Println("incompatible call rejected:", !check.OK())
}
