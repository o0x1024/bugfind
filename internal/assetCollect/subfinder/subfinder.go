package subfinder

import (
	"bufio"
	"bugfind/global"
	"bugfind/internal/utils"
	"bugfind/model/response"
	"log"
	"os/exec"
	"strings"
)

func ColletDomainByInterface(task *response.Task) (res []string) {
	log.Println("[*] subfinder start.")
	for _, vd := range task.Domains {
		for _, v := range vd.RootDomains {

			cmd := exec.Command("subfinder", "-d", v)
			buf, err := cmd.Output()
			if err != nil {
				global.Alarm <- utils.GenReportInfo(err)
				continue
			}
			var tlist []string
			scanner := bufio.NewScanner(strings.NewReader(string(buf)))
			for scanner.Scan() {
				res := scanner.Text()
				tlist = append(tlist, res)
			}

			log.Println("[+] ", v, " find ", len(tlist), "assets by subfinder.")
			for _, v := range tlist {
				res = append(res, v)
			}
		}
	}
	log.Println("[*] subfinder done.")
	return
}
