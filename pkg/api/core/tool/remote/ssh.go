package remote

import (
	"bytes"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

type Auth struct {
	Config config.SSHHost
}

func (h *Auth) SSHClientExecCmd(command string) (string, error) {
	conn, err := ssh.Dial("tcp", h.Config.HostName+":22", &ssh.ClientConfig{
		User:            h.Config.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Auth:            []ssh.AuthMethod{ssh.Password(h.Pass)},
		Auth: []ssh.AuthMethod{PublicKeyFile(h.Config.KeyPath)},
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err = session.Run(command); err != nil {
		log.Println("Failed to run: " + err.Error())
		return "", err
	}
	log.Println(command + ":" + b.String())

	return b.String(), nil
}

func (h *Auth) SSHClient() (*ssh.Client, error) {
	conn, err := ssh.Dial("tcp", h.Config.HostName+":22", &ssh.ClientConfig{
		User:            h.Config.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Auth:            []ssh.AuthMethod{ssh.Password(h.Pass)},
		Auth: []ssh.AuthMethod{PublicKeyFile(h.Config.KeyPath)},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Println(err)
		return nil
	}

	return ssh.PublicKeys(key)
}
