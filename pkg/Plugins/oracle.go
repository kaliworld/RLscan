package Plugins

import (
	common2 "RLscan/pkg/common"
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"strings"
	"time"
)

func OracleScan(info *common2.HostInfo) (tmperr error) {
	if common2.IsBrute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range common2.Userdict["oracle"] {
		for _, pass := range common2.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := OracleConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] oracle %v:%v %v %v %v", info.Host, info.Ports, user, pass, err)
				common2.LogError(errlog)
				tmperr = err
				if common2.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common2.Userdict["oracle"])*len(common2.Passwords)) * common2.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func OracleConn(info *common2.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s:%s/orcl", Username, Password, Host, Port)
	db, err := sql.Open("oracle", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(common2.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(common2.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("[+] oracle %v:%v:%v %v", Host, Port, Username, Password)
			common2.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}
