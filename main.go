package main

import (
	core "bugfind/internal/agentCore"
	"bugfind/internal/env"
	"bugfind/internal/libs/Glog"
	"bugfind/internal/taskCore"
	"fmt"
	"github.com/kataras/golog"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

//func init() {
//	configfile := flag.String("c", "conf.yaml", "the config file path")
//	flag.Parse()
//
//	cfg, err := config.LoadConfig(*configfile)
//	if err != nil {
//		panic(err)
//	}
//
//	config.Cfg = cfg
//	fmt.Println(config.Cfg)
//
//	//global.WechatKey = config.Cfg.Wechat
//	global.WkgURL = cfg.WkgUrl
//	global.V3Token = cfg.V3Token
//
//}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	sigs := make(chan os.Signal, 1)
	t := time.Now()
	Glog.InitLog()

	//环境检查，工具是否安装好了
	log.Println("--------env check.---------")
	env.EnvCheck()
	log.Println("--------env check done.-------")
	//agent注销启动
	core.Start()
	go taskCore.Run()
	//go assetCollect2.RunAssetsCollectAndBugScan()
	//go notice2.AlarmNoticeAndReport()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigs:
		fmt.Print("\n")
		golog.Info("receive signal:", sig)
		golog.Info("task exit.")
		golog.Info("Wait for the task to end.")
		time.Sleep(3)
	}

	log.Println("[*] run time:", time.Since(t))

}

//func sysvinitStart() error {
//	cmd := exec.Command("BugFind")
//	cmd.Dir = "/autoBufFind"
//	cmd.SysProcAttr = &syscall.SysProcAttr{
//		Setpgid: true,
//	}
//	for k, v := range viper.AllSettings() {
//		cmd.Env = append(cmd.Env, k+"="+v.(string))
//	}
//	cmd.Env = append(cmd.Env, "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
//	err := cmd.Start()
//	if err != nil {
//		return err
//	}
//	// set cgroup
//	quota := int64(10000)
//	memLimit := int64(262144000)
//	cg, err := cgroups.New(cgroups.V1,
//		cgroups.StaticPath("/autobugfind"),
//		&specs.LinuxResources{
//			CPU: &specs.LinuxCPU{
//				Quota: &quota,
//			},
//			Memory: &specs.LinuxMemory{
//				Limit: &memLimit,
//			},
//		})
//	if err == nil {
//		return cg.AddProc(uint64(cmd.Process.Pid))
//	}
//	return err
//}
