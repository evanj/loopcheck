# Loop Check: Check go loops for escaping variables

Analyzes Go source code to find places where a pointer to a loop variable is created. This is a common bug. See my [blog post for details](http://www.evanjones.ca/go-gotcha-loop-variables.html).


## Usage

1. `go install github.com/evanj/loopcheck`
2. `$GOPATH/bin/loopcheck (source code or package)`


## Example

Code ([go playground](http://play.golang.org/p/RFNUHJ8eyy)):

```go
package main

import "fmt"

func main() {
  values := []MyStruct{MyStruct{1}, MyStruct{2}, MyStruct{3}}
  output := []*MyStruct{}
  for _, v := range values {
    output = append(output, &v)
  }
  fmt.Println("output:", output)
}

type MyStruct struct {
  number int
}

func (m *MyStruct) String() string {
  return fmt.Sprintf("MyStruct{%d}", m.number)
}
```

Output:

```
example1.go:9: takes address of loop variable: &v
  range at line 8: for _, v := range values
```


## See Also

Inspired by errcheck: https://github.com/kisielk/errcheck

Also inspired by the range check code in go vet: https://github.com/golang/tools/blob/master/cmd/vet/rangeloop.go
