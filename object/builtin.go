package object

import (
	"fmt"
)

var Builtins = []*BuiltIn{
	{
		Name: "console.log",
		Function: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.String())
			}
			return nil
		},
	},
}
