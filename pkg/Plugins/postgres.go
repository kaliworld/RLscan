package Plugins

import (
	common2 "RLscan/pkg/common"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

func PostgresScan(info *common2.HostInfo) (tmperr error) {
	if common2.IsBrute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range common2.Userdict["postgresql"] {
		for _, pass := range common2.Passwords {
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag, err := PostgresConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] psql %v:%v %v %v %v", info.Host, info.Ports, user, pass, err)
				common2.LogError(errlog)
				tmperr = err
				if common2.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common2.Userdict["postgresql"])*len(common2.Passwords)) * common2.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func PostgresConn(info *common2.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", Username, Password, Host, Port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(common2.Timeout) * time.Second)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("[+] Postgres:%v:%v:%v %v", Host, Port, Username, Password)
			common2.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}
