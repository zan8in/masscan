package masscan

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

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

		pid int // os.getpid()

		debug bool `default:false`

		stderr, stdout bufio.Scanner
	}

	ScannerResult struct {
		IP   string `json:"ip"`
		Port string `json:"port"`
	}
)

type Option func(*Scanner)

// NewScanner Create a new Scanner, and can take options to apply to the scanner
func NewScanner(options ...Option) (*Scanner, error) {
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

	// 去掉自动检测网卡，让用户自己控制检测行为
	// dev, err := tools.AutoGetDevices()
	// if err == nil {
	// 	scanner.args = append(scanner.args, fmt.Sprintf("--interface=%s", dev.Device))
	// }

	if scanner.ctx == nil {
		scanner.ctx = context.Background()
	}

	return scanner, err
}

// Run  run masscan and returns the result of the scan
func (s *Scanner) Run() (result *tools.MasscanResult, warnings []string, err error) {
	var (
		stdout, stderr bytes.Buffer
	)

	// Enable JSON output and output in stdout
	s.args = append(s.args, "-oJ")
	s.args = append(s.args, "-")

	// debugFlag is true
	if s.debug {
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

	s.pid = cmd.Process.Pid

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

// runs asynchronously with specified arguments
func (s *Scanner) runAsync(args []string) error {
	s.cmd = exec.Command(s.binaryPath, args...)

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to get error output from asynchronous masscan run: %v", err)
	}

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to get standard output from asynchronous masscan run: %v", err)
	}

	s.stdout = *bufio.NewScanner(stdout)
	s.stderr = *bufio.NewScanner(stderr)

	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("unable to execute asynchronous masscan run: %v", err)
	}

	s.pid = s.cmd.Process.Pid

	go func() {
		<-s.ctx.Done()
		s.cmd.Process.Kill()
	}()

	return nil
}

func (s *Scanner) RunAsync() error {

	// debugFlag is true
	if s.debug {
		ss := strings.Join(s.args, " ")
		println(s.binaryPath, ss)
	}
	return s.runAsync(s.args)

}

// Pauses the masscan proccess by sending a control-c signal to the proccess
// the proccess can then be resumed by calling Resume
func (s *Scanner) PauseAsync(resumefp string) error {
	err := s.cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		return fmt.Errorf("Unable to send interrupt signal to masscan: %v", err)
	}
	s.Wait()
	err = os.Rename("paused.conf", resumefp)
	if err != nil {
		return fmt.Errorf("Unable to move resume file to new location: %v", err)
	}
	return err
}

// NB: there is a bug in the latest release of masscan that will lead to this feature not working
// but if you build masscan from source, you shouldn't experience any issues.
func (s *Scanner) ResumeAsync(resumefp string) error {
	return s.runAsync([]string{"--resume", resumefp})
}

// Wait waits for the cmd to finish and returns error.
func (s *Scanner) Wait() error {
	return s.cmd.Wait()
}

// GetStdout returns stdout variable for scanner.
func (s *Scanner) GetStdout() bufio.Scanner {
	return s.stdout
}

// GetStderr returns stderr variable for scanner.
func (s *Scanner) GetStderr() bufio.Scanner {
	return s.stderr
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

// SetConfigPath set the scanner config-file path
// eg: --conf /etc/masscan/masscan.conf
func SetConfigPath(config string) func(*Scanner) {
	return func(s *Scanner) {
		s.args = append(s.args, "--conf")
		s.args = append(s.args, config)
	}
}

// SetParamExclude sets the targets which to exclude from the scan, this also allows scanning of range 0.0.0.0/0
// eg: 127.0.0.1,255.255.255.255
func SetParamExclude(excludes ...string) func(*Scanner) {
	excludeList := strings.Join(excludes, ",")
	return func(s *Scanner) {
		s.args = append(s.args, "--exclude")
		s.args = append(s.args, excludeList)
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

// SetParamPorts sets the ports which the scanner should scan on each host.
// eg: -p 80,8000-8100
func SetParamInterface(eth string) func(*Scanner) {
	return func(s *Scanner) {
		s.args = append(s.args, fmt.Sprintf("--interface=%s", eth))
	}
}

// SetShard sets the shard number (x) and the total shard amount (y) for distributed scanning
// eg: --shard 1/2
func SetShard(x int, y int) func(*Scanner) {
	return func(s *Scanner) {
		s.args = append(s.args, fmt.Sprintf("--shard=%d/%d", x, y))
	}
}

// WithContext adds a context to a scanner, to make it cancellable and able to timeout.
func WithContext(ctx context.Context) Option {
	return func(s *Scanner) {
		s.ctx = ctx
	}
}

func ParseResult(content []byte) (sr ScannerResult) {
	result := strings.Split(string(content), " ")
	sr.IP = result[5]
	p := strings.Split(result[3], "/")
	sr.Port = p[0]
	return sr
}

// 获得 pid
func (s *Scanner) GetPid() int {
	return s.pid
}
