package lifecycle

import (
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/version"
)

// ProjectAction represents a lifecycle event action for projects.
type ProjectAction string

// All supported lifecycle events for projects.
const (
	ProjectCreated = ProjectAction(api.EventLifecycleProjectCreated)
	ProjectDeleted = ProjectAction(api.EventLifecycleProjectDeleted)
	ProjectUpdated = ProjectAction(api.EventLifecycleProjectUpdated)
	ProjectRenamed = ProjectAction(api.EventLifecycleProjectRenamed)
)

// Event creates the lifecycle event for an action on a project.
func (a ProjectAction) Event(name string, requestor *api.EventLifecycleRequestor, ctx map[string]any) api.EventLifecycle {
	u := api.NewURL().Path(version.APIVersion, "projects", name)

	return api.EventLifecycle{
		Action:    string(a),
		Source:    u.String(),
		Context:   ctx,
		Requestor: requestor,
	}
}
