package main

import "fmt"

const (
	// AppName is the cli name
	AppName = "gist"
	// AppVersion is current version of cli
	AppVersion = "2.0.0"
)

// Version show the cli's current version
func Version() string {
	return fmt.Sprintf("\n%s %s\nCopyright (c) 2017, zcong1993.", AppName, AppVersion)
}
