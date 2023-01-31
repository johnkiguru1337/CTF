package lifecycle

import (
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/version"
)

// InstanceSnapshotAction represents a lifecycle event action for instance snapshots.
type InstanceSnapshotAction string

// All supported lifecycle events for instance snapshots.
const (
	InstanceSnapshotCreated = InstanceSnapshotAction(api.EventLifecycleInstanceSnapshotCreated)
	InstanceSnapshotDeleted = InstanceSnapshotAction(api.EventLifecycleInstanceSnapshotDeleted)
	InstanceSnapshotRenamed = InstanceSnapshotAction(api.EventLifecycleInstanceSnapshotRenamed)
	InstanceSnapshotUpdated = InstanceSnapshotAction(api.EventLifecycleInstanceSnapshotUpdated)
)

// Event creates the lifecycle event for an action on an instance snapshot.
func (a InstanceSnapshotAction) Event(inst instance, ctx map[string]any) api.EventLifecycle {
	parentName, snapName, _ := api.GetParentAndSnapshotName(inst.Name())

	u := api.NewURL().Path(version.APIVersion, "instances", parentName, "snapshots", snapName).Project(inst.Project())

	var requestor *api.EventLifecycleRequestor
	if inst.Operation() != nil {
		requestor = inst.Operation().Requestor()
	}

	return api.EventLifecycle{
		Action:    string(a),
		Source:    u.String(),
		Context:   ctx,
		Requestor: requestor,
	}
}
