package lib

import "testing"

func TestShowStruct(t *testing.T) {
	err := ShowStruct("test")

	if err != nil {
		t.Error(err)
	}
}