//go:build linux && cgo && !agent

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
)

// CreateNetworkPeer creates a new Network Peer and returns its ID.
// If there is a mutual peering on the target network side the both peer entries are upated to link to each other's
// repspective network ID.
// Returns the local peer ID and true if a mutual peering has been created.
func (c *Cluster) CreateNetworkPeer(networkID int64, info *api.NetworkPeersPost) (int64, bool, error) {
	var err error
	var localPeerID int64
	var targetPeerNetworkID int64 = int64(-1) // -1 means no mutual peering exists.

	err = c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		// Insert a new Network pending peer record.
		result, err := tx.tx.Exec(`
		INSERT INTO networks_peers
		(network_id, name, description, target_network_project, target_network_name)
		VALUES (?, ?, ?, ?, ?)
		`, networkID, info.Name, info.Description, info.TargetProject, info.TargetNetwork)
		if err != nil {
			return err
		}

		localPeerID, err = result.LastInsertId()
		if err != nil {
			return err
		}

		// Save config.
		err = networkPeerConfigAdd(tx.tx, localPeerID, info.Config)
		if err != nil {
			return err
		}

		// Check if we are creating a mutual peering of an existing peer and if so then update both sides
		// with the respective network IDs. This query looks up our network peer's network name and project
		// name and then checks if there are any unlinked (target_network_id IS NULL) peers that have
		// matching target network and project names for the network this peer belongs to. If so then it
		// returns the target peer's ID and network ID. This can then be used to update both our local peer
		// and the target peer itself with the respective network IDs of each side.
		q := `
		SELECT
			target_peer.id,
			target_peer.network_id
		FROM networks_peers AS local_peer
		JOIN networks AS local_network
			ON local_network.id = local_peer.network_id
		JOIN projects AS local_project
			ON local_project.id = local_network.project_id
		JOIN networks_peers AS target_peer
			ON target_peer.target_network_name = local_network.name
			AND target_peer.target_network_project = local_project.name
		JOIN networks AS target_peer_network
			ON target_peer.network_id = target_peer_network.id
			AND target_peer_project.name = ?
		JOIN projects AS target_peer_project
			ON target_peer_network.project_id = target_peer_project.id
			AND target_peer_network.name = ?
		WHERE
			local_peer.network_id = ?
			AND local_peer.id = ?
			AND target_peer.target_network_id IS NULL
		LIMIT 1
		`

		var targetPeerID int64 = int64(-1)

		err = tx.tx.QueryRow(q, info.TargetProject, info.TargetNetwork, networkID, localPeerID).Scan(&targetPeerID, &targetPeerNetworkID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Failed looking up mutual peering: %w", err)
		} else if err == nil {
			peerNetworkMap := map[int64]struct {
				localNetworkID      int64
				targetPeerNetworkID int64
			}{
				localPeerID: {
					localNetworkID:      networkID,
					targetPeerNetworkID: targetPeerNetworkID,
				},
				targetPeerID: {
					localNetworkID:      targetPeerNetworkID,
					targetPeerNetworkID: networkID,
				},
			}

			// A mutual peering has been found, update both sides with their respective network IDs
			// and clear the joining target project and network names.
			for peerID, peerMap := range peerNetworkMap {
				_, err := tx.tx.Exec(`
				UPDATE networks_peers SET
					target_network_id = ?,
					target_network_project = NULL,
					target_network_name = NULL
				WHERE networks_peers.network_id = ? AND networks_peers.id = ?
				`, peerMap.targetPeerNetworkID, peerMap.localNetworkID, peerID)
				if err != nil {
					return fmt.Errorf("Failed updating mutual peering: %w", err)
				}
			}
		}

		return nil
	})
	if err != nil {
		return -1, false, err
	}

	return localPeerID, targetPeerNetworkID > -1, err
}

// networkPeerConfigAdd inserts Network peer config keys.
func networkPeerConfigAdd(tx *sql.Tx, peerID int64, config map[string]string) error {
	stmt, err := tx.Prepare(`
	INSERT INTO networks_peers_config
	(network_peer_id, key, value)
	VALUES(?, ?, ?)
	`)
	if err != nil {
		return err
	}

	defer func() { _ = stmt.Close() }()

	for k, v := range config {
		if v == "" {
			continue
		}

		_, err = stmt.Exec(peerID, k, v)
		if err != nil {
			return fmt.Errorf("Failed inserting config: %w", err)
		}
	}

	return nil
}

// GetNetworkPeer returns the Network Peer ID and info for the given network ID and peer name.
func (c *Cluster) GetNetworkPeer(networkID int64, peerName string) (int64, *api.NetworkPeer, error) {
	// This query loads the specified local peer as well as trying to ascertain whether there is a mutual
	// target peer, and if so what are it's project and network names. This is used to populate the
	// TargetProject, TargetNetwork fields and indicates the Status is api.NetworkStatusCreated if available.
	// If the peer is not mutually configured, then the local target_network_project and target_network_name
	// fields will be used to populate TargetProject and TargetNetwork and the Status will be set to
	// api.NetworkStatusPending.
	q := `
	SELECT
		local_peer.id,
		local_peer.name,
		local_peer.description,
		IFNULL(local_peer.target_network_project, ""),
		IFNULL(local_peer.target_network_name, ""),
		IFNULL(target_peer_network.name, "") AS target_peer_network_name,
		IFNULL(target_peer_project.name, "") AS target_peer_network_project
	FROM networks_peers AS local_peer
	LEFT JOIN networks_peers AS target_peer
		ON target_peer.network_id = local_peer.target_network_id
		AND target_peer.target_network_id = local_peer.network_id
	LEFT JOIN networks AS target_peer_network
		ON target_peer.network_id = target_peer_network.id
	LEFT JOIN projects AS target_peer_project
		ON target_peer_network.project_id = target_peer_project.id
	WHERE local_peer.network_id = ? AND local_peer.name = ?
	LIMIT 1
	`

	var err error
	var peerID int64 = int64(-1)
	var peer api.NetworkPeer
	var targetPeerNetworkName string
	var targetPeerNetworkProject string

	err = c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		err = tx.tx.QueryRow(q, networkID, peerName).Scan(&peerID, &peer.Name, &peer.Description, &peer.TargetProject, &peer.TargetNetwork, &targetPeerNetworkName, &targetPeerNetworkProject)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return api.StatusErrorf(http.StatusNotFound, "Network peer not found")
			}

			return err
		}

		err = networkPeerConfig(tx, peerID, &peer)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return -1, nil, err
	}

	networkPeerPopulatePeerInfo(&peer, targetPeerNetworkProject, targetPeerNetworkName)

	return peerID, &peer, nil
}

// networkPeerPopulatePeerInfo populates the supplied peer's Status, TargetProject and TargetNetwork fields.
// It uses the state of the targetPeerNetworkProject and targetPeerNetworkName arguments to decide whether the
// peering is mutually created and whether to use those values rather than the values contained in the peer.
func networkPeerPopulatePeerInfo(peer *api.NetworkPeer, targetPeerNetworkProject string, targetPeerNetworkName string) {
	// Peer has mutual peering from target network.
	if targetPeerNetworkName != "" && targetPeerNetworkProject != "" {
		if peer.TargetNetwork != "" || peer.TargetProject != "" {
			// Peer is in a conflicting state with both the peer network ID and net/project names set.
			// Peer net/project names should only be populated before the peer is linked with a peer
			// network ID.
			peer.Status = api.NetworkStatusErrored
		} else {
			// Peer is linked to an mutual peer on the target network.
			peer.TargetNetwork = targetPeerNetworkName
			peer.TargetProject = targetPeerNetworkProject
			peer.Status = api.NetworkStatusCreated
		}
	} else {
		if peer.TargetNetwork != "" || peer.TargetProject != "" {
			// Peer isn't linked to a mutual peer on the target network yet but has joining details.
			peer.Status = api.NetworkStatusPending
		} else {
			// Peer isn't linked to a mutual peer on the target network yet and has no joining details.
			// Perhaps it was formely joined (and had its joining details cleared) and subsequently
			// the target peer removed its peering entry.
			peer.Status = api.NetworkStatusErrored
		}
	}
}

// networkPeerConfig populates the config map of the Network Peer with the given ID.
func networkPeerConfig(tx *ClusterTx, peerID int64, peer *api.NetworkPeer) error {
	q := `
	SELECT
		key,
		value
	FROM networks_peers_config
	WHERE network_peer_id=?
	`

	peer.Config = make(map[string]string)
	return query.Scan(tx.Tx(), q, func(scan func(dest ...any) error) error {
		var key, value string

		err := scan(&key, &value)
		if err != nil {
			return err
		}

		_, found := peer.Config[key]
		if found {
			return fmt.Errorf("Duplicate config row found for key %q for network peer ID %d", key, peerID)
		}

		peer.Config[key] = value

		return nil
	}, peerID)
}

// GetNetworkPeers returns map of Network Peers for the given network ID keyed on Peer ID.
func (c *Cluster) GetNetworkPeers(networkID int64) (map[int64]*api.NetworkPeer, error) {
	// This query loads the local peers for the network as well as trying to ascertain whether there is a
	// mutual target peer, and if so what are it's project and network names. This is used to populate the
	// TargetProject, TargetNetwork fields and indicates the Status is api.NetworkStatusCreated if available.
	// If the peer is not mutually configured, then the local target_network_project and target_network_name
	// fields will be used to populate TargetProject and TargetNetwork and the Status will be set to
	// api.NetworkStatusPending.
	q := `
	SELECT
		local_peer.id,
		local_peer.name,
		local_peer.description,
		IFNULL(local_peer.target_network_project, ""),
		IFNULL(local_peer.target_network_name, ""),
		IFNULL(target_peer_network.name, "") AS target_peer_network_name,
		IFNULL(target_peer_project.name, "") AS target_peer_network_project
	FROM networks_peers AS local_peer
	LEFT JOIN networks_peers AS target_peer
		ON target_peer.network_id = local_peer.target_network_id
		AND target_peer.target_network_id = local_peer.network_id
	LEFT JOIN networks AS target_peer_network
		ON target_peer.network_id = target_peer_network.id
	LEFT JOIN projects AS target_peer_project
		ON target_peer_network.project_id = target_peer_project.id
	WHERE local_peer.network_id = ?
	`

	var err error
	peers := make(map[int64]*api.NetworkPeer)

	err = c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		err = query.Scan(tx.Tx(), q, func(scan func(dest ...any) error) error {
			var peerID int64 = int64(-1)
			var peer api.NetworkPeer
			var targetPeerNetworkName string
			var targetPeerNetworkProject string

			err := scan(&peerID, &peer.Name, &peer.Description, &peer.TargetProject, &peer.TargetNetwork, &targetPeerNetworkName, &targetPeerNetworkProject)
			if err != nil {
				return err
			}

			networkPeerPopulatePeerInfo(&peer, targetPeerNetworkProject, targetPeerNetworkName)

			peers[peerID] = &peer

			return nil
		}, networkID)
		if err != nil {
			return err
		}

		// Populate config.
		for peerID := range peers {
			err = networkPeerConfig(tx, peerID, peers[peerID])
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// GetNetworkPeerNames returns map of Network Peer names for the given network ID keyed on Peer ID.
func (c *Cluster) GetNetworkPeerNames(networkID int64) (map[int64]string, error) {
	q := `
	SELECT
		id,
		name
	FROM networks_peers
	WHERE networks_peers.network_id = ?
	`

	peers := make(map[int64]string)

	err := c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		return query.Scan(tx.Tx(), q, func(scan func(dest ...any) error) error {
			var peerID int64 = int64(-1)
			var peerName string

			err := scan(&peerID, &peerName)
			if err != nil {
				return err
			}

			peers[peerID] = peerName

			return nil
		}, networkID)
	})
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// UpdateNetworkPeer updates an existing Network Peer.
func (c *Cluster) UpdateNetworkPeer(networkID int64, peerID int64, info *api.NetworkPeerPut) error {
	err := c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		// Update existing Network peer record.
		res, err := tx.tx.Exec(`
		UPDATE networks_peers
		SET description = ?
		WHERE network_id = ? and id = ?
		`, info.Description, networkID, peerID)
		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected <= 0 {
			return api.StatusErrorf(http.StatusNotFound, "Network peer not found")
		}

		// Save config.
		_, err = tx.tx.Exec("DELETE FROM networks_peers_config WHERE network_peer_id=?", peerID)
		if err != nil {
			return err
		}

		err = networkPeerConfigAdd(tx.tx, peerID, info.Config)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteNetworkPeer deletes an existing Network Peer.
func (c *Cluster) DeleteNetworkPeer(networkID int64, peerID int64) error {
	return c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		// Delete existing Network peer record.
		res, err := tx.tx.Exec(`
			DELETE FROM networks_peers
			WHERE network_id = ? and id = ?
		`, networkID, peerID)
		if err != nil {
			return err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected <= 0 {
			return api.StatusErrorf(http.StatusNotFound, "Network peer not found")
		}

		return nil
	})
}

// NetworkPeer represents a peer connection.
type NetworkPeer struct {
	NetworkName string
	PeerName    string
}

// GetNetworkPeersTargetNetworkIDs returns a map of peer connections to target network IDs for networks in the
// specified project and network type.
func (c *Cluster) GetNetworkPeersTargetNetworkIDs(projectName string, networkType NetworkType) (map[NetworkPeer]int64, error) {
	var err error
	peerTargetNetIDs := make(map[NetworkPeer]int64)

	// Build a mapping of network and peer names to target network IDs.
	q := `SELECT p.name, n.name, p.target_network_id
		FROM networks_peers AS p
		JOIN networks AS n ON n.id = p.network_id
		JOIN projects AS pr ON pr.id = n.project_id
		WHERE pr.name = ?
		AND n.type = ?
		AND p.target_network_id > 0
	`

	err = c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		return query.Scan(tx.Tx(), q, func(scan func(dest ...any) error) error {
			var peerName string
			var networkName string
			var targetNetworkID int64 = int64(-1)

			err := scan(&peerName, &networkName, &targetNetworkID)
			if err != nil {
				return err
			}

			peer := NetworkPeer{
				PeerName:    peerName,
				NetworkName: networkName,
			}

			peerTargetNetIDs[peer] = targetNetworkID

			return nil
		}, projectName, networkType)
	})
	if err != nil {
		return nil, err
	}

	return peerTargetNetIDs, nil
}
