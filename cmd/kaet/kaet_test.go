package cmd

import (
	"reflect"
	"strings"
	"testing"
)

func TestOptionsIsRightType(t *testing.T) {
	if !strings.Contains(reflect.TypeOf(options).String(), "Options") {
		t.Fatal("options variable is not *types.Options")
	}
}
