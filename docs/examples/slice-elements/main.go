package main

import (
	"fmt"

	"github.com/qlustra/intel"
	"github.com/qlustra/intel/logic"
)

type Record struct {
	Name  string
	Count int
}

func must(check logic.Check) {
	if !check.OK() {
		panic(check.Why())
	}
}

func main() {
	batch := []Record{{Name: "Ada", Count: 2}, {Name: "Grace", Count: 3}}
	records := make([]Record, 1, 1)
	records[0] = Record{Name: "Linus", Count: 1}

	value, check := intel.CaptureValue(&records)
	must(check)
	slice, check := value.AsSlice()
	must(check)
	for _, source := range batch {
		added, check := slice.AppendZero()
		must(check)
		address, check := added.Address()
		must(check)
		*address.(*Record) = source
	}

	elements, check := slice.Elements()
	must(check)
	for i, element := range elements {
		address, check := element.Address()
		must(check)
		record := address.(*Record)
		fmt.Printf("%d: %s (%d)\n", i+1, record.Name, record.Count)
	}
	fmt.Println("records:", len(records))
}
