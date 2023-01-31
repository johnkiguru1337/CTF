package drivers

import (
	"fmt"
	"strconv"

	"golang.org/x/sys/unix"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
)

// Cmd represents a running command for an Qemu VM.
type qemuCmd struct {
	attachedChildPid int
	cmd              lxd.Operation
	dataDone         chan bool
	controlSendCh    chan api.InstanceExecControl
	controlResCh     chan error
	cleanupFunc      func()
}

// PID returns the attached child's process ID.
func (c *qemuCmd) PID() int {
	return c.attachedChildPid
}

// Signal sends a signal to the command.
func (c *qemuCmd) Signal(sig unix.Signal) error {
	command := api.InstanceExecControl{
		Command: "signal",
		Signal:  int(sig),
	}

	// Check handler hasn't finished.
	select {
	case <-c.dataDone:
		return fmt.Errorf("no such process") // Aligns with error retured from unix.Kill in lxc's Signal().
	default:
	}

	c.controlSendCh <- command
	err := <-c.controlResCh
	if err != nil {
		return err
	}

	logger.Debugf(`Forwarded signal "%d" to lxd-agent`, sig)
	return nil
}

// Wait for the command to end and returns its exit code and any error.
func (c *qemuCmd) Wait() (int, error) {
	err := c.cmd.Wait()
	if err != nil {
		return -1, err
	}

	<-c.dataDone

	exitStatus := int(c.cmd.Get().Metadata["return"].(float64))

	if c.cleanupFunc != nil {
		defer c.cleanupFunc()
	}

	return exitStatus, nil
}

// WindowResize resizes the running command's window.
func (c *qemuCmd) WindowResize(fd, winchWidth, winchHeight int) error {
	command := api.InstanceExecControl{
		Command: "window-resize",
		Args: map[string]string{
			"width":  strconv.Itoa(winchWidth),
			"height": strconv.Itoa(winchHeight),
		},
	}

	// Check handler hasn't finished.
	select {
	case <-c.dataDone:
		return fmt.Errorf("no such process") // Aligns with error retured from unix.Kill in lxc's Signal().
	default:
	}

	c.controlSendCh <- command
	err := <-c.controlResCh
	if err != nil {
		return err
	}

	logger.Debugf(`Forwarded window resize "%dx%d" to lxd-agent`, winchWidth, winchHeight)
	return nil
}
