package RLScan

import (
	common2 "RLscan/pkg/common"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var result1 string

type Addr struct {
	ip   string
	port int
}

func PortScan(hostslist []string, ports string, timeout int64, score int) ([]string, int, int, bool, string) {
	var AliveAddress []string
	var ret string
	reward := 0
	probePorts := common2.ParsePort(ports)
	if len(probePorts) == 0 {
		fmt.Printf("[-] parse port %s error, please check your port format\n", ports)
		os.Exit(0)
	}
	noPorts := common2.ParsePort(common2.NoPorts)
	if len(noPorts) > 0 {
		temp := map[int]struct{}{}
		for _, port := range probePorts {
			temp[port] = struct{}{}
		}

		for _, port := range noPorts {
			delete(temp, port)
		}

		var newDatas []int
		for port := range temp {
			newDatas = append(newDatas, port)
		}
		probePorts = newDatas
		sort.Ints(probePorts)
	}
	workers := common2.Threads
	Addrs := make(chan Addr, len(hostslist)*len(probePorts))
	results := make(chan string, len(hostslist)*len(probePorts))
	var wg sync.WaitGroup

	//接收结果
	go func() {
		for found := range results {
			AliveAddress = append(AliveAddress, found)
			wg.Done()
		}
	}()

	//多线程扫描
	for i := 0; i < workers; i++ {
		go func() {
			for addr := range Addrs {
				reward, ret = PortConnect(addr, results, timeout, &wg)
				wg.Done()
			}
		}()
	}

	//添加扫描目标
	for _, port := range probePorts {
		for _, host := range hostslist {
			wg.Add(1)
			Addrs <- Addr{host, port}
		}
	}
	wg.Wait()
	score += reward
	done := score >= 80 || score <= -30
	nextState := rand.Intn(10)
	st1 := uniqueLines(ret)
	close(Addrs)
	close(results)
	return AliveAddress, reward, nextState, done, st1
}

// 字符串按行去重
func uniqueLines(input string) string {
	var result strings.Builder
	seen := make(map[string]bool) // 创建一个map以记录出现过的行

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if !seen[line] { // 如果这一行没有出现过
			seen[line] = true               // 标记这一行为出现过
			result.WriteString(line + "\n") // 将其写入结果字符串中
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(&result, "reading input:", err)
	}

	return strings.TrimSuffix(result.String(), "\n") // 移除最后的换行符
}

func PortConnect(addr Addr, respondingHosts chan<- string, adjustedTimeout int64, wg *sync.WaitGroup) (int, string) {
	host, port := addr.ip, addr.port
	conn, err := common2.WrapperTcpWithTimeout("tcp4", fmt.Sprintf("%s:%v", host, port), time.Duration(adjustedTimeout)*time.Second)
	if err == nil {
		defer conn.Close()
		address := host + ":" + strconv.Itoa(port)
		result1 += address + " open\n"
		wg.Add(1)
		respondingHosts <- address
		return 5, result1
	} else {
		return -1, result1
	}
}

func NoPortScan(hostslist []string, ports string) (AliveAddress []string) {
	probePorts := common2.ParsePort(ports)
	noPorts := common2.ParsePort(common2.NoPorts)
	if len(noPorts) > 0 {
		temp := map[int]struct{}{}
		for _, port := range probePorts {
			temp[port] = struct{}{}
		}

		for _, port := range noPorts {
			delete(temp, port)
		}

		var newDatas []int
		for port, _ := range temp {
			newDatas = append(newDatas, port)
		}
		probePorts = newDatas
		sort.Ints(probePorts)
	}
	for _, port := range probePorts {
		for _, host := range hostslist {
			address := host + ":" + strconv.Itoa(port)
			AliveAddress = append(AliveAddress, address)
		}
	}
	return
}
