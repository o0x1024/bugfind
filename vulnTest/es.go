package vulnTest

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Es_Data_Scan()  {
	fp ,err := os.OpenFile("F:\\code\\BugFind\\es.txt",os.O_RDONLY,0777)
	if err != nil{
		log.Fatalln(err)
	}

	tr := &http.Transport{TLSClientConfig:&tls.Config{InsecureSkipVerify: true} }
	client := http.Client{Transport: tr,Timeout: 2*time.Second}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		var url string
		line:= scanner.Text()
		if !strings.Contains(line,"http") {
			url = "http://"+line + "/_cat/indices"
		}else{
			url = line + "/_cat/indices"
		}
		resp ,err := client.Get(url)
		if err != nil{
			//fmt.Println(err)
			continue
		}


		es_scanner :=bufio.NewScanner(resp.Body)
		for es_scanner.Scan() {
			esline := es_scanner.Text()
			if !strings.Contains(esline,"open"){
				break
			}
			arr := strings.Fields(esline)
			if len(arr) <= 6 {
				continue
			}
			count ,err := strconv.Atoi(arr[6])
			if err != nil{
				fmt.Println(err)
				continue
			}

			if  count > 1{
				fmt.Println("url:",url,"  index:",arr[2], "count:",arr[6])
			}
		}

	}
}