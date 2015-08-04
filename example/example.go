package main

import "fmt"

func main() {
	values := []MyStruct{MyStruct{1}, MyStruct{2}, MyStruct{3}}
	output1 := []*MyStruct{}
	for _, v := range values {
		output1 = append(output1, &v)
	}
	fmt.Println("output1: ", output1)

	output2 := []*MyStruct{}
	for i := range values {
		v := values[i]
		output2 = append(output2, &v)
	}
	fmt.Println("output2: ", output2)

	badPointers := []*int{}
	for i := range values {
		badPointers = append(badPointers, &i)
	}
	fmt.Println("badPointers: ", badPointers)

	stillWrong := []*int{}
	for _, v := range values {
		stillWrong = append(stillWrong, &v.number)
		stillWrong = append(stillWrong, &(v.number))
		stillWrong = append(stillWrong, &(v).number)
	}
	fmt.Println("stillWrong: ", stillWrong)

	// TODO: False positive: This is actually *not* a problem and should not be flagged
	actuallyFine := []*MyStruct{}
	for i, v := range values {
		actuallyFine = append(actuallyFine, &MyStruct{v.number})
		actuallyFine = append(actuallyFine, &values[i])
	}
	fmt.Println("actuallyFine: ", actuallyFine)

	// This requires a more complicated analysis to flag as a bug
	alsoWrong := []*MyStruct{}
	var v2 MyStruct
	for _, v := range values {
		v2 = v
		alsoWrong = append(alsoWrong, &v2)
	}
	fmt.Println("alsoWrong: ", alsoWrong)
}

type MyStruct struct {
	number int
}

func (m *MyStruct) String() string {
	return fmt.Sprintf("MyStruct{%d}", m.number)
}
