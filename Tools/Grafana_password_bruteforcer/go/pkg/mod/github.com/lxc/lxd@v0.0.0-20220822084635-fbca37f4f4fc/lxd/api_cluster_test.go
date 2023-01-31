package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
)

// A LXD node which is already configured for networking can be converted to a
// single-node LXD cluster.
func TestCluster_Bootstrap(t *testing.T) {
	daemon, cleanup := newTestDaemon(t)
	defer cleanup()

	// Simulate what happens when running "lxd init", where a PUT /1.0
	// request is issued to set both core.https_address and
	// cluster.https_address to the same value.
	f := clusterFixture{t: t}
	f.EnableNetworkingWithClusterAddress(daemon, "")

	client := f.ClientUnix(daemon)

	cluster := api.ClusterPut{}
	cluster.ServerName = "buzz"
	cluster.Enabled = true
	op, err := client.UpdateCluster(cluster, "")
	require.NoError(t, err)
	require.NoError(t, op.Wait())

	server, _, err := client.GetServer()
	require.NoError(t, err)
	assert.True(t, client.IsClustered())
	assert.Equal(t, "buzz", server.Environment.ServerName)
}

// Check the cluster API on a non-clustered server.
func TestCluster_Get(t *testing.T) {
	daemon, cleanup := newTestDaemon(t)
	defer cleanup()

	client, err := lxd.ConnectLXDUnix(daemon.UnixSocket(), nil)
	require.NoError(t, err)

	cluster, _, err := client.GetCluster()
	require.NoError(t, err)
	assert.Equal(t, "", cluster.ServerName)
	assert.False(t, cluster.Enabled)
}

// A LXD node can be renamed.
func TestCluster_RenameNode(t *testing.T) {
	daemon, cleanup := newTestDaemon(t)
	defer cleanup()

	f := clusterFixture{t: t}
	f.EnableNetworking(daemon, "")

	client := f.ClientUnix(daemon)

	cluster := api.ClusterPut{}
	cluster.ServerName = "buzz"
	cluster.Enabled = true
	op, err := client.UpdateCluster(cluster, "")
	require.NoError(t, err)
	require.NoError(t, op.Wait())

	node := api.ClusterMemberPost{ServerName: "rusp"}
	err = client.RenameClusterMember("buzz", node)
	require.NoError(t, err)

	_, _, err = client.GetClusterMember("rusp")
	require.NoError(t, err)
}

// Test helper for cluster-related APIs.
type clusterFixture struct {
	t       *testing.T
	clients map[*Daemon]lxd.InstanceServer
}

// Enable networking in the given daemon. The password is optional and can be
// an empty string.
func (f *clusterFixture) EnableNetworking(daemon *Daemon, password string) {
	port, err := shared.AllocatePort()
	require.NoError(f.t, err)

	address := fmt.Sprintf("127.0.0.1:%d", port)

	client := f.ClientUnix(daemon)
	server, _, err := client.GetServer()
	require.NoError(f.t, err)
	serverPut := server.Writable()
	serverPut.Config["core.https_address"] = address
	serverPut.Config["core.trust_password"] = password

	require.NoError(f.t, client.UpdateServer(serverPut, ""))
}

// Enable networking in the given daemon, and set cluster.https_address to the
// same value as core.https address. The password is optional and can be an
// empty string.
func (f *clusterFixture) EnableNetworkingWithClusterAddress(daemon *Daemon, password string) {
	port, err := shared.AllocatePort()
	require.NoError(f.t, err)

	address := fmt.Sprintf("127.0.0.1:%d", port)

	client := f.ClientUnix(daemon)
	server, _, err := client.GetServer()
	require.NoError(f.t, err)
	serverPut := server.Writable()
	serverPut.Config["core.https_address"] = address
	serverPut.Config["core.trust_password"] = password
	serverPut.Config["cluster.https_address"] = address

	require.NoError(f.t, client.UpdateServer(serverPut, ""))
}

// Get a client for the given daemon connected via UNIX socket, creating one if
// needed.
func (f *clusterFixture) ClientUnix(daemon *Daemon) lxd.InstanceServer {
	if f.clients == nil {
		f.clients = make(map[*Daemon]lxd.InstanceServer)
	}

	client, ok := f.clients[daemon]
	if !ok {
		var err error
		client, err = lxd.ConnectLXDUnix(daemon.UnixSocket(), nil)
		require.NoError(f.t, err)
	}

	return client
}
