package example

import (
	"fmt"
)

// always capitalise the functions you want to import
// it makes them export implicitly in golang
func Hello(){
	fmt.Println("example imported");
}