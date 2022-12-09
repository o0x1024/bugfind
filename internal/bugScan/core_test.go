package bugScan

import (
	"bytes"
	"fmt"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"io"
	"log"
	"regexp"
	"strings"
	"testing"
)

func TestRunBugScan(t *testing.T) {
	runnerInstance, err := runner.NewRunner(&runner.Options{
		Threads:            10,                       // Thread controls the number of threads to use for active enumerations
		Timeout:            30,                       // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10,                       // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		Resolvers:          resolve.DefaultResolvers, // Use the default list of resolvers by marshaling it to the config
		//ResultCallback: func(s *resolve.HostEntry) { // Callback function to execute for available host
		//	log.Println(s.Host, s.Source)
		//},
		Silent: true,
	})

	buf := bytes.Buffer{}
	err = runnerInstance.EnumerateSingleDomain("pingan.com", []io.Writer{&buf})
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(&buf)
	if err != nil {
		log.Fatal(err)
	}
	domains := strings.Split(string(data), "\n")
	log.Println("[+]  find ", len(domains), "assets by subfinder.")
}

func TestRunBugScans(t *testing.T) {

	var keywords = []string{`[.*(low.*]`, `\[.*medium.*\]`, `\[.*high.*\]`, `\[.*critical.*\]`}

	data := "[2022-03-15 04:19:39] [CVE-2016-6210] [network] [medium] fkysr2.iqiyi.com:22 [SSH-2.0-OpenSSH_5.3] [Hostname=fkysr2.iqiyi.com]"

	for _, v := range keywords {
		reg := regexp.MustCompile(v)
		match := reg.FindAllString(data, -1)
		if len(match) > 0 {
			//if strings.Contains(data, v) {
			fmt.Println(match)
		}
	}
}
