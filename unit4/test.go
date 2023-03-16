package main

import (
    "fmt"
    "strings"
)

func main() {
    inputString := "Hello, world, this, is, a, test"
    outputString := strings.ReplaceAll(inputString, ",", "")

    fmt.Println(outputString) // Output: Hello world this is a test
}
