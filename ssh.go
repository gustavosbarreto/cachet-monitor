package cachet

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type SSHMonitor struct {
	AbstractMonitor `mapstructure:",squash"`

	// IP:port format
	Server string

	Username string
	Password string

	// Command to execute
	Command string
}

func (monitor *SSHMonitor) Validate() []string {
	errs := monitor.AbstractMonitor.Validate()

	if monitor.Server == "" {
		errs = append(errs, "server cannot be empty")
	}

	if monitor.Username == "" {
		errs = append(errs, "username cannot be empty")
	}

	return errs
}

func (monitor *SSHMonitor) test() bool {
	sshConfig := &ssh.ClientConfig{
		User:            monitor.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(monitor.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", monitor.Server, sshConfig)
	if err != nil {
		monitor.lastFailReason = err.Error()
		return false
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		monitor.lastFailReason = err.Error()
		return false
	}

	_, err = session.CombinedOutput(monitor.Command)
	if err != nil {
		monitor.lastFailReason = errors.Wrap(err, "Failed to execute command").Error()
		return false
	}

	return true
}
