package controller

import (
	"github.com/tinyzimmer/android-farm-operator/pkg/controller/androiddevice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, androiddevice.Add)
}
