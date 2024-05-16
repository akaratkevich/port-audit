package internal

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"time"
)

/*
Set up and return an SSH session for the specified network device.
Parameters:
  - host string: The IP address or hostname of the network device.
  - port int: The port number to connect to on the network device for SSH.
  - username string: The username for SSH authentication.
  - password string: The password for SSH authentication.

Returns:
  - *ssh.Session: A pointer to an ssh.Session which can be used to execute commands on the connected device.
  - error: An error is returned if the connection or session setup fails. The error provides details about the failure.


Execution Flow:
  - The function attempts to dial an SSH connection using the provided configuration. If it fails, it logs and returns
    an error detailing the connection issue.
  - On a successful connection, it attempts to create an SSH session. If session creation fails, it logs the error,
    closes the client connection, and returns an error.
  - If all steps are successful, it returns the created session ready for command execution.
*/

func InitialiseConnection(host string, port int, username, password string) (*ssh.Session, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil // Skip host key verification
		},
		Timeout: 5 * time.Second,
	}
	log.Printf("Attempting SSH connection to %s:%d with user %s", host, port, username)
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial SSH to %s:%d: %v", host, port, err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("Failed to create session for %s: %v", host, err)
	}

	return session, nil
}
