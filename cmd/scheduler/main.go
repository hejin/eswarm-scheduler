package main

import (
	"eswarm-scheduler/pkg/eswarm"
	"k8s.io/component-base/cli"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"os"
	_ "sigs.k8s.io/scheduler-plugins/apis/config/scheme"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(eswarm.Name, eswarm.New),
	)

	code := cli.Run(command)
	os.Exit(code)
}
