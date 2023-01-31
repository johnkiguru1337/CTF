//go:build linux && cgo && !agent

package cluster

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
)

var _ = api.ServerEnvironment{}

var clusterGroupObjects = RegisterStmt(`
SELECT cluster_groups.id, cluster_groups.name, coalesce(cluster_groups.description, '')
  FROM cluster_groups
  ORDER BY cluster_groups.name
`)

var clusterGroupObjectsByName = RegisterStmt(`
SELECT cluster_groups.id, cluster_groups.name, coalesce(cluster_groups.description, '')
  FROM cluster_groups
  WHERE cluster_groups.name = ? ORDER BY cluster_groups.name
`)

var clusterGroupID = RegisterStmt(`
SELECT cluster_groups.id FROM cluster_groups
  WHERE cluster_groups.name = ?
`)

var clusterGroupCreate = RegisterStmt(`
INSERT INTO cluster_groups (name, description)
  VALUES (?, ?)
`)

var clusterGroupRename = RegisterStmt(`
UPDATE cluster_groups SET name = ? WHERE name = ?
`)

var clusterGroupDeleteByName = RegisterStmt(`
DELETE FROM cluster_groups WHERE name = ?
`)

var clusterGroupUpdate = RegisterStmt(`
UPDATE cluster_groups
  SET name = ?, description = ?
 WHERE id = ?
`)

// GetClusterGroups returns all available cluster_groups.
// generator: cluster_group GetMany
func GetClusterGroups(ctx context.Context, tx *sql.Tx, filter ClusterGroupFilter) ([]ClusterGroup, error) {
	var err error

	// Result slice.
	objects := make([]ClusterGroup, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	var args []any

	if filter.Name != nil && filter.ID == nil {
		sqlStmt, err = Stmt(tx, clusterGroupObjectsByName)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"clusterGroupObjectsByName\" prepared statement: %w", err)
		}

		args = []any{
			filter.Name,
		}
	} else if filter.ID == nil && filter.Name == nil {
		sqlStmt, err = Stmt(tx, clusterGroupObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"clusterGroupObjects\" prepared statement: %w", err)
		}

		args = []any{}
	} else {
		return nil, fmt.Errorf("No statement exists for the given Filter")
	}

	// Dest function for scanning a row.
	dest := func(scan func(dest ...any) error) error {
		c := ClusterGroup{}
		err := scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	// Select.
	err = query.SelectObjects(sqlStmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"clusters_groups\" table: %w", err)
	}

	return objects, nil
}

// GetClusterGroup returns the cluster_group with the given key.
// generator: cluster_group GetOne
func GetClusterGroup(ctx context.Context, tx *sql.Tx, name string) (*ClusterGroup, error) {
	filter := ClusterGroupFilter{}
	filter.Name = &name

	objects, err := GetClusterGroups(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"clusters_groups\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "ClusterGroup not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"clusters_groups\" entry matches")
	}
}

// GetClusterGroupID return the ID of the cluster_group with the given key.
// generator: cluster_group ID
func GetClusterGroupID(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	stmt, err := Stmt(tx, clusterGroupID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"clusterGroupID\" prepared statement: %w", err)
	}

	rows, err := stmt.Query(name)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"clusters_groups\" ID: %w", err)
	}

	defer func() { _ = rows.Close() }()

	// Ensure we read one and only one row.
	if !rows.Next() {
		return -1, api.StatusErrorf(http.StatusNotFound, "ClusterGroup not found")
	}

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("Failed to scan ID: %w", err)
	}

	if rows.Next() {
		return -1, fmt.Errorf("More than one row returned")
	}

	err = rows.Err()
	if err != nil {
		return -1, fmt.Errorf("Result set failure: %w", err)
	}

	return id, nil
}

// ClusterGroupExists checks if a cluster_group with the given key exists.
// generator: cluster_group Exists
func ClusterGroupExists(ctx context.Context, tx *sql.Tx, name string) (bool, error) {
	_, err := GetClusterGroupID(ctx, tx, name)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// RenameClusterGroup renames the cluster_group matching the given key parameters.
// generator: cluster_group Rename
func RenameClusterGroup(ctx context.Context, tx *sql.Tx, name string, to string) error {
	stmt, err := Stmt(tx, clusterGroupRename)
	if err != nil {
		return fmt.Errorf("Failed to get \"clusterGroupRename\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(to, name)
	if err != nil {
		return fmt.Errorf("Rename ClusterGroup failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows failed: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query affected %d rows instead of 1", n)
	}

	return nil
}

// CreateClusterGroup adds a new cluster_group to the database.
// generator: cluster_group Create
func CreateClusterGroup(ctx context.Context, tx *sql.Tx, object ClusterGroup) (int64, error) {
	// Check if a cluster_group with the same key exists.
	exists, err := ClusterGroupExists(ctx, tx, object.Name)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"clusters_groups\" entry already exists")
	}

	args := make([]any, 2)

	// Populate the statement arguments.
	args[0] = object.Name
	args[1] = object.Description

	// Prepared statement to use.
	stmt, err := Stmt(tx, clusterGroupCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"clusterGroupCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"clusters_groups\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"clusters_groups\" entry ID: %w", err)
	}

	return id, nil
}

// UpdateClusterGroup updates the cluster_group matching the given key parameters.
// generator: cluster_group Update
func UpdateClusterGroup(ctx context.Context, tx *sql.Tx, name string, object ClusterGroup) error {
	id, err := GetClusterGroupID(ctx, tx, name)
	if err != nil {
		return err
	}

	stmt, err := Stmt(tx, clusterGroupUpdate)
	if err != nil {
		return fmt.Errorf("Failed to get \"clusterGroupUpdate\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(object.Name, object.Description, id)
	if err != nil {
		return fmt.Errorf("Update \"clusters_groups\" entry failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query updated %d rows instead of 1", n)
	}

	return nil
}

// DeleteClusterGroup deletes the cluster_group matching the given key parameters.
// generator: cluster_group DeleteOne-by-Name
func DeleteClusterGroup(ctx context.Context, tx *sql.Tx, name string) error {
	stmt, err := Stmt(tx, clusterGroupDeleteByName)
	if err != nil {
		return fmt.Errorf("Failed to get \"clusterGroupDeleteByName\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(name)
	if err != nil {
		return fmt.Errorf("Delete \"clusters_groups\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "ClusterGroup not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d ClusterGroup rows instead of 1", n)
	}

	return nil
}
