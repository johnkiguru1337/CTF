package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/lxd/sys"
)

// The daemon is started and a client can connect to it via unix socket.
func TestIntegration_UnixSocket(t *testing.T) {
	daemon, cleanup := newTestDaemon(t)
	defer cleanup()
	client, err := lxd.ConnectLXDUnix(daemon.UnixSocket(), nil)
	require.NoError(t, err)

	server, _, err := client.GetServer()
	require.NoError(t, err)
	assert.Equal(t, "trusted", server.Auth)
	assert.False(t, server.Environment.ServerClustered)
	assert.False(t, client.IsClustered())
}

// Create a new daemon for testing.
//
// Return a function that can be used to cleanup every associated state.
func newTestDaemon(t *testing.T) (*Daemon, func()) {
	// OS
	os, osCleanup := sys.NewTestOS(t)

	// Daemon
	daemon := newDaemon(newConfig(), os)
	require.NoError(t, daemon.Init())

	cleanup := func() {
		assert.NoError(t, daemon.Stop(context.Background(), unix.SIGQUIT))
		osCleanup()
	}

	return daemon, cleanup
}

// Create a new DaemonConfig object for testing purposes.
func newConfig() *DaemonConfig {
	return &DaemonConfig{
		RaftLatency:        0.8,
		Trace:              []string{"dqlite"},
		DqliteSetupTimeout: 10 * time.Second,
	}
}
