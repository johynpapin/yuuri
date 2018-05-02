package agrovoc

import (
	"reflect"
	"testing"
)

func TestExactMatch(t *testing.T) {
	tests := []struct {
		query  string
		result []Result
	}{
		{"Oeuf", []Result{{"Oeuf", "http://aims.fao.org/aos/agrovoc/c_2502"}}},
	}

	for _, test := range tests {
		actualRes, actualErr := ExactMatch(test.query)
		if actualErr != nil {
			t.Errorf("ExactMatch(%v) returned an error: %v; want %v", test.query, actualErr, test.result)
		}
		if !reflect.DeepEqual(actualRes, test.result) {
			t.Errorf("ExactMatch(%v) = %v; want %v", test.query, actualRes, test.result)
		}
	}
}
