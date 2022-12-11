/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package eswarm

import (
	"context"
	"fmt"
	"math"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	schedulerconfig "k8s.io/kubernetes/pkg/scheduler/apis/config"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	"sigs.k8s.io/scheduler-plugins/apis/config"
)

// eSwarm is a score plugin that favors nodes based on their allocatable
// resources.
type eSwarm struct {
	handle framework.Handle
	resourceAllocationScorer
}

var _ = framework.ScorePlugin(&eSwarm{})

// eSwarmName is the name of the plugin used in the Registry and configurations.
const Name = "eswarm"

// the M1 pro single thread performance, comparing to Intel i7 gen8 - both in mobile
// Intel i7 gen8 as the baseline
const CpuCorePerfRatio = 4

// Name returns name of the plugin. It is used in logs, etc.
func (alloc *eSwarm) Name() string {
	return Name
}

func validateResources(resources []schedulerconfig.ResourceSpec) error {
	for _, resource := range resources {
		if resource.Weight <= 0 {
			return fmt.Errorf("resource Weight of %v should be a positive value, got %v", resource.Name, resource.Weight)
		}
		// No upper bound on weight.
	}
	return nil
}

// Score invoked at the score extension point.
func (alloc *eSwarm) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := alloc.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}

	node := nodeInfo.Node()
	nodeNameArch := fmt.Sprintf("(name: %v, arch: %v)", nodeName, node.Status.NodeInfo.Architecture)
	klog.V(0).InfoS("scheduler : Score() for node", "node details", nodeNameArch)

	// alloc.score favors nodes with least allocatable or most allocatable resources.
	// It calculates the sum of the node's weighted allocatable resources.
	//
	// Note: the returned "score" is negative for least allocatable, and positive for most allocatable.
	return alloc.score(pod, nodeInfo)
}

// ScoreExtensions of the Score plugin.
func (alloc *eSwarm) ScoreExtensions() framework.ScoreExtensions {
	return alloc
}

// NeweSwarm initializes a new plugin and returns it.
func New(allocArgs runtime.Object, h framework.Handle) (framework.Plugin, error) {
	// Start with default values.
	mode := config.Least
	resToWeightMap := defaultResourcesToWeightMap

	// Update values from args, if specified.
	if allocArgs != nil {
		args, ok := allocArgs.(*config.NodeResourcesAllocatableArgs)
		if !ok {
			return nil, fmt.Errorf("want args to be of type, got NodeResourcesAllocatableArgs %T", allocArgs)
		}
		if args.Mode != "" {
			mode = args.Mode
			if mode != config.Least && mode != config.Most {
				return nil, fmt.Errorf("invalid mode, got %s", mode)
			}
		}

		if len(args.Resources) > 0 {
			if err := validateResources(args.Resources); err != nil {
				return nil, err
			}

			resToWeightMap = make(resourceToWeightMap)
			for _, resource := range args.Resources {
				resToWeightMap[v1.ResourceName(resource.Name)] = resource.Weight
			}
		}
	}

	return &eSwarm{
		handle: h,
		resourceAllocationScorer: resourceAllocationScorer{
			Name:                Name,
			scorer:              resourceScorer(resToWeightMap, mode),
			resourceToWeightMap: resToWeightMap,
		},
	}, nil
}

func resourceScorer(resToWeightMap resourceToWeightMap, mode config.ModeType) func(bool, resourceToValueMap, resourceToValueMap) int64 {
	return func(useCpuPerfRatio bool, requested, allocable resourceToValueMap) int64 {
		// TODO: consider volumes in scoring.
		var nodeScore, weightSum int64
		for resource, weight := range resToWeightMap {
			resourceScore := score(allocable[resource], mode)
			nodeScore += resourceScore * weight
			weightSum += weight
		}
		if useCpuPerfRatio {
			nodeScore *= CpuCorePerfRatio
		}
		return nodeScore / weightSum
	}
}

func score(capacity int64, mode config.ModeType) int64 {
	switch config.ModeType(mode) {
	case config.Least:
		return -1 * capacity
	case config.Most:
		return capacity
	}

	klog.V(10).InfoS("No match for mode", "mode", mode)
	return 0
}

// NormalizeScore invoked after scoring all nodes.
func (alloc *eSwarm) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	// Find highest and lowest scores.
	var highest int64 = -math.MaxInt64
	var lowest int64 = math.MaxInt64
	for _, nodeScore := range scores {
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}

	// Transform the highest to lowest score range to fit the framework's min to max node score range.
	oldRange := highest - lowest
	newRange := framework.MaxNodeScore - framework.MinNodeScore
	for i, nodeScore := range scores {
		if oldRange == 0 {
			scores[i].Score = framework.MinNodeScore
		} else {
			scores[i].Score = ((nodeScore.Score - lowest) * newRange / oldRange) + framework.MinNodeScore
		}
	}

	return nil
}
