package evaluator

import (
	"fmt"

	"github.com/jf550-kent/jsgo/object"
)

var builtin = map[string]*object.BuiltIn{
	"console.log": {
		Name: "console.log",
		Function: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.String())
			}
			return NULL
		},
	},
}
