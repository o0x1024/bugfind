package domain

import (
	"BugFind/internal/libs/ksubdomain/core"
	"log"
	"strconv"
	"testing"
)

func TestColletDomainByKSubdomain(t *testing.T) {
	var rate int64
	bandwith := "1M"
	suffix := string([]rune(bandwith)[len(bandwith)-1])
	rate, _ = strconv.ParseInt(string([]rune(bandwith)[0:len(bandwith)-1]), 10, 64)
	switch suffix {
	case "G":
		fallthrough
	case "g":
		rate *= 1000000000
	case "M":
		fallthrough
	case "m":
		rate *= 1000000
	case "K":
		fallthrough
	case "k":
		rate *= 1000
	}
	packSize := int64(100) // 一个DNS包大概有74byte
	rate = rate / packSize

	option := &core.Options{
		Resolvers:    []string{"223.5.5.5", "223.6.6.6", "180.76.76.76", "119.29.29.29", "182.254.116.116", "114.114.114.115"},
		Rate:         rate,
		Domain:       []string{"momo.com"},
		Silent:       true,
		SkipWildCard: true,
		Summary:      true,
	}
	subdomins, err := core.Start(option)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("subdomins:", subdomins)
	//for _, v := range vd.RootDomains {
	//
	//}

}
