package lifecycle

import (
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/version"
)

// StoragePoolAction represents a lifecycle event action for storage pools.
type StoragePoolAction string

// All supported lifecycle events for storage pools.
const (
	StoragePoolCreated = StoragePoolAction(api.EventLifecycleStoragePoolCreated)
	StoragePoolDeleted = StoragePoolAction(api.EventLifecycleStoragePoolDeleted)
	StoragePoolUpdated = StoragePoolAction(api.EventLifecycleStoragePoolUpdated)
)

// Event creates the lifecycle event for an action on an storage pool.
func (a StoragePoolAction) Event(name string, requestor *api.EventLifecycleRequestor, ctx map[string]any) api.EventLifecycle {
	u := api.NewURL().Path(version.APIVersion, "storage-pools", name)

	return api.EventLifecycle{
		Action:    string(a),
		Source:    u.String(),
		Context:   ctx,
		Requestor: requestor,
	}
}
