package Plugins

import (
	common2 "RLscan/pkg/common"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

func SshScan(info *common2.HostInfo) (tmperr error) {
	if common2.IsBrute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range common2.Userdict["ssh"] {
		for _, pass := range common2.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := SshConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] ssh %v:%v %v %v %v", info.Host, info.Ports, user, pass, err)
				common2.LogError(errlog)
				tmperr = err
				if common2.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common2.Userdict["ssh"])*len(common2.Passwords)) * common2.Timeout) {
					return err
				}
			}
			if common2.SshKey != "" {
				return err
			}
		}
	}
	return tmperr
}

func SshConn(info *common2.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	var Auth []ssh.AuthMethod
	if common2.SshKey != "" {
		pemBytes, err := ioutil.ReadFile(common2.SshKey)
		if err != nil {
			return false, errors.New("read key failed" + err.Error())
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return false, errors.New("parse key failed" + err.Error())
		}
		Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		Auth = []ssh.AuthMethod{ssh.Password(Password)}
	}

	config := &ssh.ClientConfig{
		User:    Username,
		Auth:    Auth,
		Timeout: time.Duration(common2.Timeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", Host, Port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		if err == nil {
			defer session.Close()
			flag = true
			var result string
			if common2.Command != "" {
				combo, _ := session.CombinedOutput(common2.Command)
				result = fmt.Sprintf("[+] SSH %v:%v:%v %v \n %v", Host, Port, Username, Password, string(combo))
				if common2.SshKey != "" {
					result = fmt.Sprintf("[+] SSH %v:%v sshkey correct \n %v", Host, Port, string(combo))
				}
				common2.LogSuccess(result)
			} else {
				result = fmt.Sprintf("[+] SSH %v:%v:%v %v", Host, Port, Username, Password)
				if common2.SshKey != "" {
					result = fmt.Sprintf("[+] SSH %v:%v sshkey correct", Host, Port)
				}
				common2.LogSuccess(result)
			}
		}
	}
	return flag, err

}
