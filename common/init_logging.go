package common

import (
	log "github.com/sirupsen/logrus"
)

func InitLogging(quiet bool, verbose bool) {
	if quiet && verbose {
		log.Fatal("Invalid logging flag combination: cannot turn on both quiet and verbose modes")
	}

	level := log.InfoLevel
	if verbose {
		level = log.DebugLevel
	}
	if quiet {
		level = log.ErrorLevel
	}
	log.SetLevel(level)
}
