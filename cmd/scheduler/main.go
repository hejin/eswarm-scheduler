package main

import (
	"eswarm-scheduler/pkg/eswarm"
	//"eswarm-scheduler/pkg/register"
	//"flag"
	"fmt"
	//"k8s.io/component-base/cli"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"math/rand"
	"os"
	_ "sigs.k8s.io/scheduler-plugins/apis/config/scheme"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	//command := register.Register()
	command := app.NewSchedulerCommand(app.WithPlugin(eswarm.Name, eswarm.New))
	logs.InitLogs()
	//flag.Parse()
	defer logs.FlushLogs()

	klog.Infof("Welcome to eSwarm!!!")
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	klog.Infof("Bye from eSwarm...")
}
