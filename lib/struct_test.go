package lib

import "testing"

func TestShowStruct(t *testing.T) {
	err := ShowStruct("articles","form")

	if err != nil {
		t.Error(err)
	}
}