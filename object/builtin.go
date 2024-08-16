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

var ArrayPush = &BuiltIn{
	Name: "push",
	Function: func(args ...Object) Object {
		if len(args) != 2 {
			panic("built in array push function only accept 2 arguments")
		}
		arr, ok := args[0].(*Array)
		if !ok {
			panic("not an array for push function")
		}
		arr.Body = append(arr.Body, args[1])
		return nil
	},
}
