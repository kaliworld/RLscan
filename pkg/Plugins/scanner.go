package Plugins

import (
	Q_learning "RLscan/pkg/RL/Q-learning"
	"RLscan/pkg/RL/utlis"
	"RLscan/pkg/WebScan/lib"
	common2 "RLscan/pkg/common"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

func Scan(info common2.HostInfo) {
	fmt.Println("start infoscan")
	// IP 解析
	Hosts, err := common2.ParseIP(info.Host, common2.HostFile, common2.NoHosts)
	if err != nil {
		fmt.Println("len(hosts)==0", err)
		return
	}
	lib.Inithttp()
	var ch = make(chan struct{}, common2.Threads)
	var wg = sync.WaitGroup{}
	// 根据扫描到的端口在做对应检测
	web := strconv.Itoa(common2.PORTList["web"])
	ms17010 := strconv.Itoa(common2.PORTList["ms17010"])
	if len(Hosts) > 0 || len(common2.HostPort) > 0 {
		if common2.NoPing == false && len(Hosts) > 1 || common2.Scantype == "icmp" {
			// 存活探测
			Hosts = CheckLive(Hosts, common2.Ping)
			fmt.Println("[*] Icmp alive hosts len is:", len(Hosts))
		}
		if common2.Scantype == "icmp" {
			common2.LogWG.Wait()
			return
		}
		var AlivePorts []string
		if common2.Scantype == "webonly" || common2.Scantype == "webpoc" {
			AlivePorts = NoPortScan(Hosts, common2.Ports)
		} else if common2.Scantype == "hostname" {
			common2.Ports = "139"
			AlivePorts = NoPortScan(Hosts, common2.Ports)
		} else if common2.ReinforcementLearning {
			AlivePorts = PortScan(Hosts, common2.Ports, common2.Timeout)
			fmt.Println("[+] RL扫描中")
			AlivePorts1 := Q_learning.Run(Hosts, common2.Timeout)
			AlivePorts = utlis.MergeSlicesExcludeDuplicates(AlivePorts, AlivePorts1)
			fmt.Println("[*] alive ports len is:", len(AlivePorts))
			if common2.Scantype == "portscan" {
				common2.LogWG.Wait()
				return
			}
		} else if common2.Scantype != "RL" {
			AlivePorts = PortScan(Hosts, common2.Ports, common2.Timeout)
			fmt.Println("[*] alive ports len is:", len(AlivePorts))
			if common2.Scantype == "portscan" {
				common2.LogWG.Wait()
				return
			}
		}
		if len(common2.HostPort) > 0 {
			AlivePorts = append(AlivePorts, common2.HostPort...)
			AlivePorts = common2.RemoveDuplicate(AlivePorts)
			common2.HostPort = nil
			fmt.Println("[*] AlivePorts len is:", len(AlivePorts))
		}
		var severports []string //severports := []string{"21","22","135"."445","1433","3306","5432","6379","9200","11211","27017"...}
		for _, port := range common2.PORTList {
			severports = append(severports, strconv.Itoa(port))
		}
		fmt.Println("start vulscan")
		for _, targetIP := range AlivePorts {
			host, _, _ := net.SplitHostPort(targetIP)
			if ip := net.ParseIP(host); ip != nil {
				if ip.To4() != nil {
					info.Host, info.Ports = strings.Split(targetIP, ":")[0], strings.Split(targetIP, ":")[1]
				} else if ip.To16() != nil {
					info.Host, info.Ports, err = net.SplitHostPort(targetIP)
				}
			}
			if common2.Scantype == "all" || common2.Scantype == "main" || common2.Scantype == "RL" {
				switch {
				case info.Ports == "135":
					AddScan(info.Ports, info, &ch, &wg) //findnet
					if common2.IsWmi {
						AddScan("1000005", info, &ch, &wg) //wmiexec
					}
				case info.Ports == "445":
					AddScan(ms17010, info, &ch, &wg)    //ms17010
					AddScan(info.Ports, info, &ch, &wg) //smb
					AddScan("1000002", info, &ch, &wg)  //smbghost
				case info.Ports == "9000":
					AddScan(web, info, &ch, &wg)        //http
					AddScan(info.Ports, info, &ch, &wg) //fcgiscan
				case IsContain(severports, info.Ports):
					AddScan(info.Ports, info, &ch, &wg) //plugins scan
				default:
					AddScan(web, info, &ch, &wg) //webtitle
				}
			} else {
				scantype := strconv.Itoa(common2.PORTList[common2.Scantype])
				AddScan(scantype, info, &ch, &wg)
			}
		}
	}
	for _, url := range common2.Urls {
		info.Url = url
		AddScan(web, info, &ch, &wg)
	}
	wg.Wait()
	common2.LogWG.Wait()
	close(common2.Results)
	fmt.Printf("已完成 %v/%v\n", common2.End, common2.Num)
}

var Mutex = &sync.Mutex{}

func AddScan(scantype string, info common2.HostInfo, ch *chan struct{}, wg *sync.WaitGroup) {
	*ch <- struct{}{}
	wg.Add(1)
	go func() {
		Mutex.Lock()
		common2.Num += 1
		Mutex.Unlock()
		ScanFunc(&scantype, &info)
		Mutex.Lock()
		common2.End += 1
		Mutex.Unlock()
		wg.Done()
		<-*ch
	}()
}

func ScanFunc(name *string, info *common2.HostInfo) {
	if name == nil {
		return
	}
	f, exists := PluginList[*name]
	if !exists {
		return
	}
	reflectValue := reflect.ValueOf(f)
	if reflectValue.Kind() == reflect.Func {
		in := []reflect.Value{reflect.ValueOf(info)}
		reflectValue.Call(in)
	} else {
		fmt.Println("[-] 不存在的函数调用f")
	}
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
