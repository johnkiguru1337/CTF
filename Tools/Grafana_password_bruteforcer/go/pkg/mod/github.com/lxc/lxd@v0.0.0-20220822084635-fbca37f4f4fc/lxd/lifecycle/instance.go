package lifecycle

import (
	"github.com/lxc/lxd/lxd/operations"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/version"
)

// Internal copy of the instance interface.
type instance interface {
	Name() string
	Project() string
	Operation() *operations.Operation
}

// InstanceAction represents a lifecycle event action for instances.
type InstanceAction string

// All supported lifecycle events for instances.
const (
	InstanceCreated          = InstanceAction(api.EventLifecycleInstanceCreated)
	InstanceStarted          = InstanceAction(api.EventLifecycleInstanceStarted)
	InstanceStopped          = InstanceAction(api.EventLifecycleInstanceStopped)
	InstanceShutdown         = InstanceAction(api.EventLifecycleInstanceShutdown)
	InstanceRestarted        = InstanceAction(api.EventLifecycleInstanceRestarted)
	InstancePaused           = InstanceAction(api.EventLifecycleInstancePaused)
	InstanceResumed          = InstanceAction(api.EventLifecycleInstanceResumed)
	InstanceRestored         = InstanceAction(api.EventLifecycleInstanceRestored)
	InstanceDeleted          = InstanceAction(api.EventLifecycleInstanceDeleted)
	InstanceRenamed          = InstanceAction(api.EventLifecycleInstanceRenamed)
	InstanceUpdated          = InstanceAction(api.EventLifecycleInstanceUpdated)
	InstanceExec             = InstanceAction(api.EventLifecycleInstanceExec)
	InstanceConsole          = InstanceAction(api.EventLifecycleInstanceConsole)
	InstanceConsoleRetrieved = InstanceAction(api.EventLifecycleInstanceConsoleRetrieved)
	InstanceConsoleReset     = InstanceAction(api.EventLifecycleInstanceConsoleReset)
	InstanceFileRetrieved    = InstanceAction(api.EventLifecycleInstanceFileRetrieved)
	InstanceFilePushed       = InstanceAction(api.EventLifecycleInstanceFilePushed)
	InstanceFileDeleted      = InstanceAction(api.EventLifecycleInstanceFileDeleted)
)

// Event creates the lifecycle event for an action on an instance.
func (a InstanceAction) Event(inst instance, ctx map[string]any) api.EventLifecycle {
	url := api.NewURL().Path(version.APIVersion, "instances", inst.Name()).Project(inst.Project())

	var requestor *api.EventLifecycleRequestor
	if inst.Operation() != nil {
		requestor = inst.Operation().Requestor()
	}

	return api.EventLifecycle{
		Action:    string(a),
		Source:    url.String(),
		Context:   ctx,
		Requestor: requestor,
	}
}
