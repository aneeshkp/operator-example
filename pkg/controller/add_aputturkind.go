package controller

import (
	"github.com/aneeshkp/operator-example/pkg/controller/aputturkind"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, aputturkind.Add)
}
