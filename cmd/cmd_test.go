package cmd

import (
	"fmt"
	"testing"
)

func TestReconstructHNURL(t *testing.T) {
	var (
		in       = "item?id=16838460"
		expected = fmt.Sprintf("%v/%v", baseURL, in)
	)

	actual := reconstructHNURL(in)

	if actual != expected {
		t.Errorf("Expected fixed URL=%v but actual result=%v", expected, actual)
	}
}
