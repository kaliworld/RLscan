package main

import (
	"RLscan/pkg/AL"
	"RLscan/pkg/Plugins"
	common2 "RLscan/pkg/common"
	"fmt"
	"time"
)

func main() {
	//Q_learning.Train1()

	start := time.Now()
	var Info common2.HostInfo
	common2.Flag(&Info)
	common2.Parse(&Info)
	if common2.AdaptiveLearning {
		AL.Run(Info)
	}
	Plugins.Scan(Info)
	fmt.Printf("[*] 扫描结束,耗时: %s\n", time.Since(start))
}
