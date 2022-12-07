package register

import (
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"pkg/eswarm"
)

func Register() *cobra.Command {
	return app.NewSchedulerCommand(
		app.WithPlugin(eswarm.Name, eswarm.New),
	)
}
