package AL

import (
	"RLscan/pkg/common"
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const (
	// Ubuntu 的平均值
	DEFAULT_FILE_DESCRIPTORS_LIMIT = 8000 // 举例值
	// 基于实验的最安全批量大小
	AVERAGE_BATCH_SIZE = 3000
)

func Run(info common.HostInfo) {
	if runtime.GOOS == "windows" {
		// issues#4056: "runtime: limit number of operating system threads"
		common.Threads = 10000
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("bash", "-c", "ulimit -n")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("命令执行错误: %s\n", err)
			return
		}
		ulimit := strings.TrimSpace(out.String())
		ulimit1, err := strconv.Atoi(ulimit)
		if err != nil {
			fmt.Println("转换错误:", err)
			return
		}
		if ulimit1 > common.Threads {
			// 若操作系统支持较高的文件限制，例如8000，但用户设置的批处理大小超出此值，则需要减少
			if ulimit1 > DEFAULT_FILE_DESCRIPTORS_LIMIT {
				common.Threads = ulimit1 / 2
			} else if ulimit1 < AVERAGE_BATCH_SIZE {
				common.Threads = AVERAGE_BATCH_SIZE
			} else {
				common.Threads = ulimit1 - 100
			}
		}
	} else {
		fmt.Printf("你在运行其他系统: %s\n，不支持除了wiodows外的系统", runtime.GOOS)
	}
}
