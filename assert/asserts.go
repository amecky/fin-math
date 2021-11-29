package assert

import (
	"fmt"
	"reflect"
	"testing"
)

func Equals(t *testing.T, message string, expected, got float64) {
	es := fmt.Sprintf("%.2f", expected)
	gs := fmt.Sprintf("%.2f", got)
	if es != gs {
		t.Fatalf("%s - expected: %s but got %s", message, es, gs)
	}
}

func NotNull(t *testing.T, message string, d interface{}) {
	if d == nil {
		t.Fatal(message)
	}
	value := reflect.ValueOf(d)
	if value.IsNil() {
		t.Fatal(message)
	}
}

func Null(t *testing.T, message string, d interface{}) {
	if d != nil {
		t.Fatal(message)
	}
	value := reflect.ValueOf(d)
	if value.IsNil() == false {
		t.Fatal(message)
	}
}
