package errors

import "errors"

var (
	// ErrMasscanNotInstalled  masscan binary was not found
	ErrMasscanNotInstalled = errors.New("masscan binary was not found")

	// ErrScanTimout  masscan scan timed out
	ErrScanTimout = errors.New("masscan scan timed out")

	// ErrParseOutput  unable to parse masscan output, see warning for details
	ErrParseOutput = errors.New("unable to parse masscan output, see warning for details")

	// ErrResolveName  masscan could not resolve a name
	ErrResolveName = errors.New("masscan could not resolve a name")
)
