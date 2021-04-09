package provider

import (
	"reflect"
	"testing"
)

func TestExpandStringList(t *testing.T) {
	expanded := []interface{}{
		"log(\"hello 1\");",
		"log(\"hello 2\");",
	}
	stringList := expandStringList(expanded)
	expected := []string{
		"log(\"hello 1\");",
		"log(\"hello 2\");",
	}

	if !reflect.DeepEqual(stringList, expected) {
		t.Fatalf(
			"Got:\n\n%#v\n\nExpected:\n\n%#v\n",
			stringList,
			expected)
	}
}
