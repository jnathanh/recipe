package lib_test

import (
	"testing"

	"github.com/jnathanh/recipe/lib"
)

func TestList(t *testing.T) {
	_, err := lib.List()
	if err != nil {
		t.Error(err)
	}
}
