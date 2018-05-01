package common

import (
	"fmt"
	"strings"
)

// ReconstructHNURL turns a relative HN link into a full URL.
func ReconstructHNURL(u string) string {
	if strings.HasPrefix(u, "item?id=") {
		return fmt.Sprintf("%v/%v", BaseURL, u)
	}
	return u
}
