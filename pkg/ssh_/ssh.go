package ssh_

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func ConnectToHost(host string, user string, pass string, cipher []string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.Ciphers = append(sshConfig.Ciphers, "3des-cbc", "twofish128-cbc")

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func RunCommand(client *ssh.Client, command string, sleep int) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return "", err
	}

	defer session.Close()

	var buffer bytes.Buffer
	session.Stdout = &buffer
	session.Stderr = &buffer

	stdin, err := session.StdinPipe()
	if err != nil {
		return "", err
	}

	err = session.Shell()
	if err != nil {
		return "", err
	}

	stdin.Write([]byte(fmt.Sprintf("%s\n", command)))

	// XXX: This is a ugly workaround to wait than the command returns his result.
	// We should instead find a way to send correctly a EOF (CTRL+D) and use `session.wait()`.
	if sleep != 0 {
		time.Sleep(time.Duration(sleep) * time.Second)
	}

	return buffer.String(), nil
}
