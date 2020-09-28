package util

import (
    "fmt"
    "testing"
)

func TestGenPass(t *testing.T) {
    fmt.Println(GenPass("base3", 16))
}
