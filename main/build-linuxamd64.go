//go:build linux && amd64

package main

// these are gonna be populated at build
var (
	Version = "2.0.0"
	GitHash = "unknown"
	Target  = "linux-amd64"
	Build   = Version + "-" + GitHash + "-" + Target
)
