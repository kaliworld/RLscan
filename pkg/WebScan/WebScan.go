package WebScan

import (
	lib2 "RLscan/pkg/WebScan/lib"
	common2 "RLscan/pkg/common"
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed pocs
var Pocs embed.FS
var once sync.Once
var AllPocs []*lib2.Poc

func WebScan(info *common2.HostInfo) {
	once.Do(initpoc)
	var pocinfo = common2.Pocinfo
	buf := strings.Split(info.Url, "/")
	pocinfo.Target = strings.Join(buf[:3], "/")

	if pocinfo.PocName != "" {
		Execute(pocinfo)
	} else {
		for _, infostr := range info.Infostr {
			pocinfo.PocName = lib2.CheckInfoPoc(infostr)
			Execute(pocinfo)
		}
	}
}

func Execute(PocInfo common2.PocInfo) {
	req, err := http.NewRequest("GET", PocInfo.Target, nil)
	if err != nil {
		errlog := fmt.Sprintf("[-] webpocinit %v %v", PocInfo.Target, err)
		common2.LogError(errlog)
		return
	}
	req.Header.Set("User-agent", common2.UserAgent)
	req.Header.Set("Accept", common2.Accept)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	if common2.Cookie != "" {
		req.Header.Set("Cookie", common2.Cookie)
	}
	pocs := filterPoc(PocInfo.PocName)
	lib2.CheckMultiPoc(req, pocs, common2.PocNum)
}

func initpoc() {
	if common2.PocPath == "" {
		entries, err := Pocs.ReadDir("pocs")
		if err != nil {
			fmt.Printf("[-] init poc error: %v", err)
			return
		}
		for _, one := range entries {
			path := one.Name()
			if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
				if poc, _ := lib2.LoadPoc(path, Pocs); poc != nil {
					AllPocs = append(AllPocs, poc)
				}
			}
		}
	} else {
		fmt.Println("[+] load poc from " + common2.PocPath)
		err := filepath.Walk(common2.PocPath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil || info == nil {
					return err
				}
				if !info.IsDir() {
					if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
						poc, _ := lib2.LoadPocbyPath(path)
						if poc != nil {
							AllPocs = append(AllPocs, poc)
						}
					}
				}
				return nil
			})
		if err != nil {
			fmt.Printf("[-] init poc error: %v", err)
		}
	}
}

func filterPoc(pocname string) (pocs []*lib2.Poc) {
	if pocname == "" {
		return AllPocs
	}
	for _, poc := range AllPocs {
		if strings.Contains(poc.Name, pocname) {
			pocs = append(pocs, poc)
		}
	}
	return
}
