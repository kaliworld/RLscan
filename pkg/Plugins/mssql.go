package Plugins

import (
	common2 "RLscan/pkg/common"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"strings"
	"time"
)

func MssqlScan(info *common2.HostInfo) (tmperr error) {
	if common2.IsBrute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range common2.Userdict["mssql"] {
		for _, pass := range common2.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MssqlConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] mssql %v:%v %v %v %v", info.Host, info.Ports, user, pass, err)
				common2.LogError(errlog)
				tmperr = err
				if common2.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common2.Userdict["mssql"])*len(common2.Passwords)) * common2.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func MssqlConn(info *common2.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	dataSourceName := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%v;encrypt=disable;timeout=%v", Host, Username, Password, Port, time.Duration(common2.Timeout)*time.Second)
	db, err := sql.Open("mssql", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(common2.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(common2.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("[+] mssql %v:%v:%v %v", Host, Port, Username, Password)
			common2.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}
