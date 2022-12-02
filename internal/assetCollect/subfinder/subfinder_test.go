package subfinder

import (
	"bytes"
	"context"
	"fmt"
	"github.com/projectdiscovery/subfinder/v2/pkg/passive"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"io"
	"log"
	"strings"
	"testing"
)

func TestColletBySubfinder(t *testing.T) {

	providers := runner.Providers{
		Bufferover:  []string{"w4d4sjjB0y2sJapUatYNK8kmQMJ8eDQ77tT6ytJC"},
		Binaryedge:  []string{"de203f66-e2aa-45ea-b79f-bdb566633463"},
		Censys:      []string{"5bd2a5da-abf7-4f9b-b216-16fa003c9d1d"},
		Certspotter: []string{"k46684_dFiZBAF7HNaSTxhqFU3o"},
		Chaos:       []string{"97f6445660321b71387e088dbbc8bb5e77a0d48240519ed5b5e0e607c378e687"},
		Chinaz:      []string{"4b20375e736a40b181dc380f41c391cd"},
		GitHub:      []string{"p_VdG9J5LR6TBLOPDiVTJLbJGqJA0yjU0J7aHS"},
		Shodan:      []string{"U5ityH3ya9gcGMrKdvHcaKD2xyJEsgn8"},
		ThreatBook:  []string{"f7bff3be38934f778d70c7a5fc6097387638c347625d41379a6f6ca0f2ae89e0"},
		URLScan:     []string{"f40240d3-54a6-4778-a401-7d83e6ee83ac"},
		Virustotal:  []string{"209e09bd70fb39a37987218def266691e0bb834a91b26d56039eed82c02b2e71"},
		Fofa:        []string{"zjgelen@gmail.com:1bd13cc61d22823099fea2a8e26f7478"},
	}

	runnerInstance, err := runner.NewRunner(&runner.Options{
		Threads:            10,                              // Thread controls the number of threads to use for active enumerations
		Timeout:            30,                              // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10,                              // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		Resolvers:          resolve.DefaultResolvers,        // Use the default list of resolvers by marshaling it to the config
		Sources:            passive.DefaultSources,          // Use the default list of passive sources
		AllSources:         passive.DefaultAllSources,       // Use the default list of all passive sources
		Recursive:          passive.DefaultRecursiveSources, // Use the default list of recursive sources
		Providers:          &providers,                      // Use empty api keys for all providers

	})

	buf := bytes.Buffer{}
	err = runnerInstance.EnumerateSingleDomain(context.Background(), "momo.com", []io.Writer{&buf})
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(&buf)
	if err != nil {
		log.Fatal(err)
	}
	domains := strings.Split(string(data), "\n")
	fmt.Printf("%s", domains)
}
