package eswarm

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	// Name is plugin name
	Name = "eswarm"
)

var (
	_ framework.FilterPlugin  = &eSwarm{}
	_ framework.PreBindPlugin = &eSwarm{}
)

type eSwarm struct {
	handle framework.Handle
}

func (s *eSwarm) Name() string {
	return Name
}
func New(_ runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &eSwarm{
		handle: handle,
	}, nil
}

func (s *eSwarm) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, node *framework.NodeInfo) *framework.Status {
	klog.Infof("eswarm --> filter pod: %v", pod.Name)
	return framework.NewStatus(framework.Success, "")
}

func (s *eSwarm) PreBind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return framework.NewStatus(framework.Error, err.Error())
	}
	klog.Infof("eswarm --> prebind node info: %+v", nodeInfo.Node())
	return framework.NewStatus(framework.Success, "")
}
