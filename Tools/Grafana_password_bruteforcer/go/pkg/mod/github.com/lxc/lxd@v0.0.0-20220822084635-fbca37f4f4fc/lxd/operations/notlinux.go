//go:build !linux || !cgo || agent

package operations

import (
	"fmt"

	"github.com/lxc/lxd/lxd/db/operationtype"
	"github.com/lxc/lxd/shared/api"
)

func registerDBOperation(op *Operation, opType operationtype.Type) error {
	if op.state != nil {
		return fmt.Errorf("registerDBOperation not supported on this platform")
	}

	return nil
}

func removeDBOperation(op *Operation) error {
	if op.state != nil {
		return fmt.Errorf("registerDBOperation not supported on this platform")
	}

	return nil
}

func (op *Operation) sendEvent(eventMessage any) {
	if op.events == nil {
		return
	}

	op.events.Send(op.projectName, api.EventTypeOperation, eventMessage)
}
