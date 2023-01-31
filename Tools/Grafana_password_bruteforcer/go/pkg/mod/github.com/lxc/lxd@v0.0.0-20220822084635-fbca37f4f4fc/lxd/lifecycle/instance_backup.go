package lifecycle

import (
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/version"
)

// InstanceBackupAction represents a lifecycle event action for instance backups.
type InstanceBackupAction string

// All supported lifecycle events for instance backups.
const (
	InstanceBackupCreated   = InstanceBackupAction(api.EventLifecycleInstanceBackupCreated)
	InstanceBackupDeleted   = InstanceBackupAction(api.EventLifecycleInstanceBackupDeleted)
	InstanceBackupRenamed   = InstanceBackupAction(api.EventLifecycleInstanceBackupRenamed)
	InstanceBackupRetrieved = InstanceBackupAction(api.EventLifecycleInstanceBackupRetrieved)
)

// Event creates the lifecycle event for an action on an instance backup.
func (a InstanceBackupAction) Event(fullBackupName string, inst instance, ctx map[string]any) api.EventLifecycle {
	_, backupName, _ := api.GetParentAndSnapshotName(fullBackupName)

	u := api.NewURL().Path(version.APIVersion, "instances", inst.Name(), "backups", backupName).Project(inst.Project())

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
