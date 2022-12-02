package env

import (
	"bugfind/internal/libs/Glog"
	"fmt"
	"log"
	"os/exec"
)

//subfinder、redis、nuclei

func EnvCheck() {
	cmd := exec.Command("nuclei", "-help")
	_, err := cmd.Output()
	if err != nil {
		//go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest
		Glog.InfoG("[!] [nuclei] not install.  [go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest]")
		log.Println("[!] [nuclei] not install.  [go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest]")
	}
	Glog.InfoG("nuclei installed.")
	log.Println("[+] nuclei installed.")
	//
	//cmd = exec.Command("subfinder", "-help")
	//_, err = cmd.Output()
	//if err != nil {
	//	fmt.Println("[!] can not found [nuclei]. \n[*] please run [go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest] to install.")
	//	os.Exit(0)
	//}
	//log.Println("[+] subfinder installed.")

	cmd = exec.Command("httpx", "-help")
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("[!] can not found [httpx]. \n[*] please run [go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest] to install.")

	}
	log.Println("[+] httpx installed.")

	//_os := runtime.GOOS
	//existRedis := false
	//if _os == "linux" {
	//	cmd := exec.Command("ps", "-ef")
	//	buf, err := cmd.Output()
	//	if err != nil {
	//		log.Println("[!] ps -ef error.")
	//		os.Exit(0)
	//	}
	//	scanner := bufio.NewScanner(strings.NewReader(string(buf)))
	//	for scanner.Scan() {
	//		content := scanner.Text()
	//		if strings.Contains(content, "redis") {
	//			existRedis = true
	//		}
	//	}
	//	if !existRedis {
	//		log.Println("[!] redis is not running or not to installed . please run after installation")
	//		os.Exit(0)
	//	}
	//} else if _os == "windows" {
	//	cmd := exec.Command("tasklist")
	//	buf, err := cmd.Output()
	//	if err != nil {
	//		log.Println("[!] tasklist error.")
	//		os.Exit(0)
	//	}
	//	scanner := bufio.NewScanner(strings.NewReader(string(buf)))
	//	for scanner.Scan() {
	//		content := scanner.Text()
	//		if strings.Contains(content, "redis") {
	//			existRedis = true
	//		}
	//	}
	//	if !existRedis {
	//		log.Println("[!] redis is not running or not to installed  . please run after installation ")
	//		os.Exit(0)
	//	}
	//}
	//log.Println("[+] reids is running.")

}
