package globalstore

import (
	"testing"
)

func TestSG(t *testing.T) {
	err := S("key0", []byte("value0"))
	if err != nil {
		t.Error(err)
	}

	err = S("key1", []byte("value1"))
	if err != nil {
		t.Error(err)
	}

	value, err := G("key0")
	if err != nil {
		t.Error(err)
	}

	if string(value) != "value0" {
		t.Error("value is not equal (key0)")
	}

	value, err = G("key1")
	if err != nil {
		t.Error(err)
	}

	if string(value) != "value1" {
		t.Error("value is not equal (key1)")
	}
}
