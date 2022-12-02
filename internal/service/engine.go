package service

import (
	"BugFind/global"
	"BugFind/internal/handle"
	"BugFind/internal/libs/Glog"
	utils2 "BugFind/internal/utils"
	"BugFind/model/request"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func AgentRegister() (err error) {
	OS := runtime.GOOS
	hostname, _ := os.Hostname()

	ext := utils2.CheckFileIsExist("machine-id")
	if ext {
		fp, err := os.OpenFile("machine-id", os.O_RDONLY, 0666)
		if err != nil {
			Glog.ErrorG("%s", err.Error())
			os.Exit(1)
		}

		buf, err := ioutil.ReadAll(fp)
		if err != nil {
			Glog.ErrorG("%s", err.Error())
			os.Exit(1)
		}
		global.AgentId = string(buf)
	}

	id := xid.New().String()
	agentName := OS + "-" + hostname + "-" + global.Version + "-" + id
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	ips := ""
	privateIPv4 := ""
	privateIPv6 := ""
	publicIPv4 := ""
	publicIPv6 := ""

	for _, i := range interfaces {

		if strings.HasPrefix(i.Name, "docker") || strings.HasPrefix(i.Name, "lo") {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil || !ip.IsGlobalUnicast() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				if (ip4[0] == 10) || (ip4[0] == 192 && ip4[1] == 168) || (ip4[0] == 172 && ip4[1]&0x10 == 0x10) {
					privateIPv4 = privateIPv4 + "," + ip4.String()
				} else {
					publicIPv4 = (publicIPv4 + "," + ip4.String())
				}
			} else if len(ip) == net.IPv6len {
				if ip[0] == 0xfd {
					privateIPv6 = (privateIPv6 + "," + ip.String())
				} else {
					publicIPv6 = (publicIPv6 + "," + ip.String())
				}
			}
		}
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return err
		}
		addresses, err := byName.Addrs()
		for _, v := range addresses {
			if ips == "" {
				ips = "{\"name\":\"" + byName.Name + "\",\"ip\":\"" + v.String() + "\"}"
			} else {
				ips += ",{\"name\":\"" + byName.Name + "\",\"ip\":\"" + v.String() + "\"}"
			}
		}
	}

	pid := os.Getpid()
	envMap := make(map[string]string)
	envs := os.Environ()
	for _, v := range envs {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			continue
		} else {
			envMap[parts[0]] = parts[1]
		}
	}
	envB, err := json.Marshal(envMap)
	if err != nil {
		return err
	}
	encodeEnv := base64.StdEncoding.EncodeToString(envB)

	filePath, err := utils2.GetCurrentPath()
	if err != nil {
		Glog.ErrorG("%s", err.Error())
		return err
	}
	req := request.AgentRegisterReq{
		AgentId:     global.AgentId,
		Name:        agentName,
		Version:     global.Version,
		PrivateIPv4: privateIPv4[1:],
		PublicIPv4:  publicIPv4,
		PrivateIPv6: privateIPv6,
		PublicIPv6:  publicIPv6,
		Hostname:    hostname,
		Network:     ips,
		ServerPath:  filePath,
		ServerEnv:   encodeEnv,
		Pid:         strconv.Itoa(pid),
	}

	go func() {
		agentId, err := handle.AgentRegister(req)
		if err != nil {
			fmt.Println("agent register err", err)
			os.Exit(1)
		}
		global.AgentId = agentId

		if !ext {
			fp, err := os.OpenFile("machine-id", os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				Glog.ErrorG("%s", err.Error())
				os.Exit(1)
			}

			fp.WriteString(agentId)
			if err := fp.Close(); err != nil {
				Glog.ErrorG("%s", err.Error())
			}

		}

	}()

	return nil

}
