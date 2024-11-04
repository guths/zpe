package main

import (
	"fmt"

	factory "github.com/guths/zpe/factory/factories"
)

func main() {
	f := factory.NewUserFactory()

	u, _ := f.Create()

	fmt.Printf("%+v", u)

}
