package masscan

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/zan8in/masscan/errors"
	"github.com/zan8in/masscan/tools"
)

type (
	// Scanner ...
	Scanner struct {
		cmd *exec.Cmd

		args       []string
		binaryPath string
		ctx        context.Context

		debug bool `default:false`

		stderr, stdout bufio.Scanner
	}
)

// NewScanner Create a new Scanner, and can take options to apply to the scanner
func NewScanner(options ...func(*Scanner)) (*Scanner, error) {
	scanner := &Scanner{}
	var err error

	for _, opt := range options {
		opt(scanner)
	}

	if scanner.binaryPath == "" {
		scanner.binaryPath, err = exec.LookPath("masscan")
		if err != nil {
			return nil, errors.ErrMasscanNotInstalled
		}
	}

	if scanner.ctx == nil {
		scanner.ctx = context.Background()
	}

	return scanner, err
}

// Run  run masscan and returns the result of the scan
func (s *Scanner) Run() (result *tools.MasscanResult, warnings []string, err error) {
	var stdout, stderr bytes.Buffer

	// Enable JSON output and output in stdout
	s.args = append(s.args, "-oJ")
	s.args = append(s.args, "-")

	// debugFlag is true
	if s.debug == true {
		ss := strings.Join(s.args, " ")
		println(s.binaryPath, ss)
	}

	// Prepare masscan process
	cmd := exec.Command(s.binaryPath, s.args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		return nil, warnings, err
	}

	// Make a goroutine to notify the select when the scan is done
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for masscan process or timout
	select {
	case <-s.ctx.Done():

		// Context was done before the scan was finished
		// Killed process.
		_ = cmd.Process.Kill()

		// return a timeout error
		return nil, warnings, errors.ErrScanTimout
	case <-done:

		// Although it is a warning, it needs to be clear
		if stderr.Len() > 0 {
			warnings = strings.Split(strings.Trim(stderr.String(), "\n"), "\n")
		}

		// Parse masscan JSON ouput
		if stdout.Len() > 0 {
			result, err := tools.ParseJson(stdout.Bytes())
			if err != nil {
				warnings = append(warnings, err.Error())
				return nil, warnings, errors.ErrParseOutput
			}
			// Return result, optional warnings but no error
			return result, warnings, err
		}

	}
	return nil, nil, nil
}

// EnableDebug set debug mode is true
// eg:C:\masscan\masscan.exe 146.56.202.100/24 -p 80,8000-8100 --rate=10000 -oJ -
func EnableDebug() func(*Scanner) {
	return func(s *Scanner) {
		s.debug = true
	}
}

// SetParamTargets set the target of a scanner
// eg: 192.168.88.0/24  192.168.88.0-255  192.168.88.0.255-192.168.88.255
func SetParamTargets(targets ...string) func(*Scanner) {
	return func(s *Scanner) {
		s.args = append(s.args, targets...)
	}
}

// SetParamPorts sets the ports which the scanner should scan on each host.
// eg: -p 80,8000-8100
func SetParamPorts(ports ...string) func(*Scanner) {
	portList := strings.Join(ports, ",")
	return func(s *Scanner) {
		s.args = append(s.args, "-p")
		s.args = append(s.args, portList)
	}
}

// SetParamTopPorts eg: --top-ports
func SetParamTopPorts() func(*Scanner) {
	return func(s *Scanner) {
		s.args = append(s.args, "--top-ports")
	}
}

// SetParamRate set the rate
// masscan -p80,8000-8100 10.0.0.0/8 --rate=10000
// scan some web ports on 10.x.x.x at 10kpps
func SetParamRate(maxRate int) func(*Scanner) {
	return func(s *Scanner) {
		s.args = append(s.args, fmt.Sprintf("--rate=%d", maxRate))
	}
}

// SetParamWait The waiting time after sending the packet, the default is 10 seconds
// --wait=10s   default is 10 seconds
func SetParamWait(delay int) func(*Scanner) {
	return func(s *Scanner) {
		s.args = append(s.args, fmt.Sprintf("--wait=%d", delay))
	}
}
