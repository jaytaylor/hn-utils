package common

import (
	"fmt"
	"testing"
)

func TestReconstructHNURL(t *testing.T) {
	var (
		fragment = "item?id=16838460"
		expected = fmt.Sprintf("%v/%v", BaseURL, fragment)
	)

	actual := ReconstructHNURL(fragment)

	if actual != expected {
		t.Errorf("Expected fixed URL=%v but actual result=%v", expected, actual)
	}
}
