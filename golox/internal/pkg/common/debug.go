package common

import "strings"

var DEBUGLOX = 0

const (
	SCANNING = 1 << iota
)

func SetDebug(settings string) {
	settingsSlice := strings.Split(settings, ",")

	for _, set := range settingsSlice {
		switch strings.ToLower(set) {
		case "scanning":
			DEBUGLOX |= SCANNING
		}
	}
}
