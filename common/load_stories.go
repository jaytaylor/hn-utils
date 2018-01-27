package common

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jaytaylor/hn-utils/domain"

	log "github.com/sirupsen/logrus"
)

// LoadStories loads a array of stories from the named file.
// "-" can be used to signify readying from STDIN.
func LoadStories(filename string) (domain.Stories, error) {
	var (
		stories domain.Stories
		r       io.Reader
	)

	if filename == "-" {
		r = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("opening %v: %s", filename, err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Warnf("Unexpected problem closing %v: %s", filename, err)
			}
		}()
		r = file
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(&stories); err != nil {
		return nil, fmt.Errorf("loading stories from %v: %s", filename, err)
	}
	log.Debugf("Loaded %v stories from %v", len(stories), filename)

	return stories, nil
}
