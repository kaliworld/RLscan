package Plugins

import (
	common2 "RLscan/pkg/common"
	"errors"
	"fmt"
	"github.com/stacktitan/smb/smb"
	"strings"
	"time"
)

func SmbScan(info *common2.HostInfo) (tmperr error) {
	if common2.IsBrute {
		return nil
	}
	starttime := time.Now().Unix()
	for _, user := range common2.Userdict["smb"] {
		for _, pass := range common2.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := doWithTimeOut(info, user, pass)
			if flag == true && err == nil {
				var result string
				if common2.Domain != "" {
					result = fmt.Sprintf("[+] SMB %v:%v:%v\\%v %v", info.Host, info.Ports, common2.Domain, user, pass)
				} else {
					result = fmt.Sprintf("[+] SMB %v:%v:%v %v", info.Host, info.Ports, user, pass)
				}
				common2.LogSuccess(result)
				return err
			} else {
				errlog := fmt.Sprintf("[-] smb %v:%v %v %v %v", info.Host, 445, user, pass, err)
				errlog = strings.Replace(errlog, "\n", "", -1)
				common2.LogError(errlog)
				tmperr = err
				if common2.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common2.Userdict["smb"])*len(common2.Passwords)) * common2.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func SmblConn(info *common2.HostInfo, user string, pass string, signal chan struct{}) (flag bool, err error) {
	flag = false
	Host, Username, Password := info.Host, user, pass
	options := smb.Options{
		Host:        Host,
		Port:        445,
		User:        Username,
		Password:    Password,
		Domain:      common2.Domain,
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			flag = true
		}
	}
	signal <- struct{}{}
	return flag, err
}

func doWithTimeOut(info *common2.HostInfo, user string, pass string) (flag bool, err error) {
	signal := make(chan struct{})
	go func() {
		flag, err = SmblConn(info, user, pass, signal)
	}()
	select {
	case <-signal:
		return flag, err
	case <-time.After(time.Duration(common2.Timeout) * time.Second):
		return false, errors.New("time out")
	}
}
