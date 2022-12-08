package eswarm

import (
	"context"
	"errors"
	//"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	//"k8s.io/apimachinery/pkg/types"
	//scv "github.com/NJUPT-ISL/SCV/api/v1"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

const (
	// Name is plugin name
	Name = "eswarm"
)

var (
	_ framework.FilterPlugin  = &eSwarm{}
	_ framework.PreBindPlugin = &eSwarm{}
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "core.run-linux.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
	schemeLocal = runtime.NewScheme()
)

type eSwarm struct {
	handle framework.Handle
	cache  cache.Cache
}

func New(_ runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	mgrConfig := ctrl.GetConfigOrDie()
	mgrConfig.QPS = 1000
	mgrConfig.Burst = 1000

	//if err := scv.AddToScheme(schemeLocal); err != nil {
	if err := AddToScheme(schemeLocal); err != nil {
		klog.Error(err)
		return nil, err
	}

	mgr, err := ctrl.NewManager(mgrConfig, ctrl.Options{
		Scheme:             schemeLocal,
		MetricsBindAddress: "",
		LeaderElection:     false,
		Port:               9443,
	})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	go func() {
		if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			klog.Error(err)
			panic(err)
		}
	}()
	scvCache := mgr.GetCache()

	if scvCache.WaitForCacheSync(context.TODO()) {
		return &eSwarm{
			handle: handle,
			cache:  scvCache,
		}, nil
	} else {
		return nil, errors.New("Cache Not Sync! ")
	}
}

func (s *eSwarm) Name() string {
	return Name
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
