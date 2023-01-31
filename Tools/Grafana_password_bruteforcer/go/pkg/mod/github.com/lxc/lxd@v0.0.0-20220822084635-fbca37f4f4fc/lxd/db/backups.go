//go:build linux && cgo && !agent

package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
)

// InstanceBackup is a value object holding all db-related details about an instance backup.
type InstanceBackup struct {
	ID                   int
	InstanceID           int
	Name                 string
	CreationDate         time.Time
	ExpiryDate           time.Time
	InstanceOnly         bool
	OptimizedStorage     bool
	CompressionAlgorithm string
}

// StoragePoolVolumeBackup is a value object holding all db-related details about a storage volume backup.
type StoragePoolVolumeBackup struct {
	ID                   int
	VolumeID             int64
	Name                 string
	CreationDate         time.Time
	ExpiryDate           time.Time
	VolumeOnly           bool
	OptimizedStorage     bool
	CompressionAlgorithm string
}

// Returns the ID of the instance backup with the given name.
func (c *Cluster) getInstanceBackupID(name string) (int, error) {
	q := "SELECT id FROM instances_backups WHERE name=?"
	id := -1
	arg1 := []any{name}
	arg2 := []any{&id}
	err := dbQueryRowScan(c, q, arg1, arg2)
	if err == sql.ErrNoRows {
		return -1, api.StatusErrorf(http.StatusNotFound, "Instance backup not found")
	}

	return id, err
}

// GetInstanceBackup returns the backup with the given name.
func (c *Cluster) GetInstanceBackup(projectName string, name string) (InstanceBackup, error) {
	args := InstanceBackup{}
	args.Name = name

	instanceOnlyInt := -1
	optimizedStorageInt := -1
	q := `
SELECT instances_backups.id, instances_backups.instance_id,
       instances_backups.creation_date, instances_backups.expiry_date,
       instances_backups.container_only, instances_backups.optimized_storage
    FROM instances_backups
    JOIN instances ON instances.id=instances_backups.instance_id
    JOIN projects ON projects.id=instances.project_id
    WHERE projects.name=? AND instances_backups.name=?
`
	arg1 := []any{projectName, name}
	arg2 := []any{&args.ID, &args.InstanceID, &args.CreationDate,
		&args.ExpiryDate, &instanceOnlyInt, &optimizedStorageInt}
	err := dbQueryRowScan(c, q, arg1, arg2)
	if err != nil {
		if err == sql.ErrNoRows {
			return args, api.StatusErrorf(http.StatusNotFound, "Instance backup not found")
		}

		return args, err
	}

	if instanceOnlyInt == 1 {
		args.InstanceOnly = true
	}

	if optimizedStorageInt == 1 {
		args.OptimizedStorage = true
	}

	return args, nil
}

// GetInstanceBackupWithID returns the backup with the given ID.
func (c *Cluster) GetInstanceBackupWithID(backupID int) (InstanceBackup, error) {
	args := InstanceBackup{}
	args.ID = backupID

	instanceOnlyInt := -1
	optimizedStorageInt := -1
	q := `
SELECT instances_backups.name, instances_backups.instance_id,
       instances_backups.creation_date, instances_backups.expiry_date,
       instances_backups.container_only, instances_backups.optimized_storage
    FROM instances_backups
    JOIN instances ON instances.id=instances_backups.instance_id
    JOIN projects ON projects.id=instances.project_id
    WHERE instances_backups.id=?
`
	arg1 := []any{backupID}
	arg2 := []any{&args.Name, &args.InstanceID, &args.CreationDate,
		&args.ExpiryDate, &instanceOnlyInt, &optimizedStorageInt}
	err := dbQueryRowScan(c, q, arg1, arg2)
	if err != nil {
		if err == sql.ErrNoRows {
			return args, api.StatusErrorf(http.StatusNotFound, "Instance backup not found")
		}

		return args, err
	}

	if instanceOnlyInt == 1 {
		args.InstanceOnly = true
	}

	if optimizedStorageInt == 1 {
		args.OptimizedStorage = true
	}

	return args, nil
}

// GetInstanceBackups returns the names of all backups of the instance with the
// given name.
func (c *Cluster) GetInstanceBackups(projectName string, name string) ([]string, error) {
	var result []string

	q := `SELECT instances_backups.name FROM instances_backups
JOIN instances ON instances_backups.instance_id=instances.id
JOIN projects ON projects.id=instances.project_id
WHERE projects.name=? AND instances.name=?`
	inargs := []any{projectName, name}
	outfmt := []any{name}
	dbResults, err := queryScan(c, q, inargs, outfmt)
	if err != nil {
		return nil, err
	}

	for _, r := range dbResults {
		result = append(result, r[0].(string))
	}

	return result, nil
}

// CreateInstanceBackup creates a new backup.
func (c *Cluster) CreateInstanceBackup(args InstanceBackup) error {
	_, err := c.getInstanceBackupID(args.Name)
	if err == nil {
		return ErrAlreadyDefined
	}

	err = c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		instanceOnlyInt := 0
		if args.InstanceOnly {
			instanceOnlyInt = 1
		}

		optimizedStorageInt := 0
		if args.OptimizedStorage {
			optimizedStorageInt = 1
		}

		str := "INSERT INTO instances_backups (instance_id, name, creation_date, expiry_date, container_only, optimized_storage) VALUES (?, ?, ?, ?, ?, ?)"
		stmt, err := tx.tx.Prepare(str)
		if err != nil {
			return err
		}

		defer func() { _ = stmt.Close() }()
		result, err := stmt.Exec(args.InstanceID, args.Name,
			args.CreationDate.Unix(), args.ExpiryDate.Unix(), instanceOnlyInt,
			optimizedStorageInt)
		if err != nil {
			return err
		}

		_, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("Error inserting %q into database", args.Name)
		}

		return nil
	})

	return err
}

// DeleteInstanceBackup removes the instance backup with the given name from the database.
func (c *Cluster) DeleteInstanceBackup(name string) error {
	id, err := c.getInstanceBackupID(name)
	if err != nil {
		return err
	}

	err = exec(c, "DELETE FROM instances_backups WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}

// RenameInstanceBackup renames an instance backup from the given current name
// to the new one.
func (c *Cluster) RenameInstanceBackup(oldName, newName string) error {
	err := c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		str := "UPDATE instances_backups SET name = ? WHERE name = ?"
		stmt, err := tx.tx.Prepare(str)
		if err != nil {
			return err
		}

		defer func() { _ = stmt.Close() }()

		logger.Debug(
			"Calling SQL Query",
			logger.Ctx{
				"query":   "UPDATE instances_backups SET name = ? WHERE name = ?",
				"oldName": oldName,
				"newName": newName})
		_, err = stmt.Exec(newName, oldName)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// GetExpiredInstanceBackups returns a list of expired instance backups.
func (c *Cluster) GetExpiredInstanceBackups() ([]InstanceBackup, error) {
	var result []InstanceBackup
	var name string
	var expiryDate string
	var instanceID int

	q := `SELECT instances_backups.name, instances_backups.expiry_date, instances_backups.instance_id FROM instances_backups`
	outfmt := []any{name, expiryDate, instanceID}
	dbResults, err := queryScan(c, q, nil, outfmt)
	if err != nil {
		return nil, err
	}

	for _, r := range dbResults {
		timestamp := r[1]

		var backupExpiry time.Time
		err = backupExpiry.UnmarshalText([]byte(timestamp.(string)))
		if err != nil {
			return []InstanceBackup{}, err
		}

		// Since zero time causes some issues due to timezones, we check the
		// unix timestamp instead of IsZero().
		if backupExpiry.Unix() <= 0 {
			// Backup doesn't expire
			continue
		}

		// Backup has expired
		if time.Now().Unix()-backupExpiry.Unix() >= 0 {
			result = append(result, InstanceBackup{
				Name:       r[0].(string),
				InstanceID: r[2].(int),
				ExpiryDate: backupExpiry,
			})
		}
	}

	return result, nil
}

// GetStoragePoolVolumeBackups returns a list of volume backups.
func (c *Cluster) GetStoragePoolVolumeBackups(projectName string, volumeName string, poolID int64) ([]StoragePoolVolumeBackup, error) {
	q := `
	SELECT
		backups.id,
		backups.storage_volume_id,
		backups.name,
		backups.creation_date,
		backups.expiry_date,
		backups.volume_only,
		backups.optimized_storage
	FROM storage_volumes_backups AS backups
	JOIN storage_volumes ON storage_volumes.id=backups.storage_volume_id
	JOIN projects ON projects.id=storage_volumes.project_id
	WHERE projects.name=? AND storage_volumes.name=? AND storage_volumes.storage_pool_id=?
	ORDER BY backups.id
	`

	var backups []StoragePoolVolumeBackup

	err := c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		return query.Scan(tx.Tx(), q, func(scan func(dest ...any) error) error {
			var b StoragePoolVolumeBackup
			var expiryTime sql.NullTime

			err := scan(&b.ID, &b.VolumeID, &b.Name, &b.CreationDate, &expiryTime, &b.VolumeOnly, &b.OptimizedStorage)
			if err != nil {
				return err
			}

			b.ExpiryDate = expiryTime.Time // Convert nulls to zero.

			backups = append(backups, b)

			return nil
		}, projectName, volumeName, poolID)
	})
	if err != nil {
		return nil, err
	}

	return backups, nil
}

// GetStoragePoolVolumeBackupsNames returns the names of all backups of the storage volume with the given name.
func (c *Cluster) GetStoragePoolVolumeBackupsNames(projectName string, volumeName string, poolID int64) ([]string, error) {
	var result []string

	q := `SELECT storage_volumes_backups.name FROM storage_volumes_backups
JOIN storage_volumes ON storage_volumes_backups.storage_volume_id=storage_volumes.id
JOIN projects ON projects.id=storage_volumes.project_id
WHERE projects.name=? AND storage_volumes.name=?
ORDER BY storage_volumes_backups.id`
	inargs := []any{projectName, volumeName}
	outfmt := []any{volumeName}
	dbResults, err := queryScan(c, q, inargs, outfmt)
	if err != nil {
		return nil, err
	}

	for _, r := range dbResults {
		result = append(result, r[0].(string))
	}

	return result, nil
}

// CreateStoragePoolVolumeBackup creates a new storage volume backup.
func (c *Cluster) CreateStoragePoolVolumeBackup(args StoragePoolVolumeBackup) error {
	_, err := c.getStoragePoolVolumeBackupID(args.Name)
	if err == nil {
		return ErrAlreadyDefined
	}

	err = c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		volumeOnlyInt := 0
		if args.VolumeOnly {
			volumeOnlyInt = 1
		}

		optimizedStorageInt := 0
		if args.OptimizedStorage {
			optimizedStorageInt = 1
		}

		str := "INSERT INTO storage_volumes_backups (storage_volume_id, name, creation_date, expiry_date, volume_only, optimized_storage) VALUES (?, ?, ?, ?, ?, ?)"
		stmt, err := tx.tx.Prepare(str)
		if err != nil {
			return err
		}

		defer func() { _ = stmt.Close() }()
		result, err := stmt.Exec(args.VolumeID, args.Name,
			args.CreationDate.Unix(), args.ExpiryDate.Unix(), volumeOnlyInt,
			optimizedStorageInt)
		if err != nil {
			return err
		}

		_, err = result.LastInsertId()
		if err != nil {
			return fmt.Errorf("Error inserting %q into database", args.Name)
		}

		return nil
	})

	return err
}

// Returns the ID of the storage volume backup with the given name.
func (c *Cluster) getStoragePoolVolumeBackupID(name string) (int, error) {
	q := "SELECT id FROM storage_volumes_backups WHERE name=?"
	id := -1
	arg1 := []any{name}
	arg2 := []any{&id}
	err := dbQueryRowScan(c, q, arg1, arg2)
	if err == sql.ErrNoRows {
		return -1, api.StatusErrorf(http.StatusNotFound, "Storage volume backup not found")
	}

	return id, err
}

// DeleteStoragePoolVolumeBackup removes the storage volume backup with the given name from the database.
func (c *Cluster) DeleteStoragePoolVolumeBackup(name string) error {
	id, err := c.getStoragePoolVolumeBackupID(name)
	if err != nil {
		return err
	}

	err = exec(c, "DELETE FROM storage_volumes_backups WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}

// GetStoragePoolVolumeBackup returns the volume backup with the given name.
func (c *Cluster) GetStoragePoolVolumeBackup(projectName string, poolName string, backupName string) (StoragePoolVolumeBackup, error) {
	args := StoragePoolVolumeBackup{}
	q := `
SELECT
	backups.id,
	backups.storage_volume_id,
	backups.name,
	backups.creation_date,
	backups.expiry_date,
	backups.volume_only,
	backups.optimized_storage
FROM storage_volumes_backups AS backups
JOIN storage_volumes ON storage_volumes.id=backups.storage_volume_id
JOIN projects ON projects.id=storage_volumes.project_id
WHERE projects.name=? AND backups.name=?
`
	arg1 := []any{projectName, backupName}
	outfmt := []any{&args.ID, &args.VolumeID, &args.Name, &args.CreationDate, &args.ExpiryDate, &args.VolumeOnly, &args.OptimizedStorage}
	err := dbQueryRowScan(c, q, arg1, outfmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return args, api.StatusErrorf(http.StatusNotFound, "Storage volume backup not found")
		}

		return args, err
	}

	return args, nil
}

// GetStoragePoolVolumeBackupWithID returns the volume backup with the given ID.
func (c *Cluster) GetStoragePoolVolumeBackupWithID(backupID int) (StoragePoolVolumeBackup, error) {
	args := StoragePoolVolumeBackup{}
	q := `
SELECT
	backups.id,
	backups.storage_volume_id,
	backups.name,
	backups.creation_date,
	backups.expiry_date,
	backups.volume_only,
	backups.optimized_storage
FROM storage_volumes_backups AS backups
JOIN storage_volumes ON storage_volumes.id=backups.storage_volume_id
JOIN projects ON projects.id=storage_volumes.project_id
WHERE backups.id=?
`
	arg1 := []any{backupID}
	outfmt := []any{&args.ID, &args.VolumeID, &args.Name, &args.CreationDate, &args.ExpiryDate, &args.VolumeOnly, &args.OptimizedStorage}
	err := dbQueryRowScan(c, q, arg1, outfmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return args, api.StatusErrorf(http.StatusNotFound, "Storage volume backup not found")
		}

		return args, err
	}

	return args, nil
}

// RenameVolumeBackup renames a volume backup from the given current name
// to the new one.
func (c *Cluster) RenameVolumeBackup(oldName, newName string) error {
	err := c.Transaction(context.TODO(), func(ctx context.Context, tx *ClusterTx) error {
		str := "UPDATE storage_volumes_backups SET name = ? WHERE name = ?"
		stmt, err := tx.tx.Prepare(str)
		if err != nil {
			return err
		}

		defer func() { _ = stmt.Close() }()

		logger.Debug(
			"Calling SQL Query",
			logger.Ctx{
				"query":   "UPDATE storage_volumes_backups SET name = ? WHERE name = ?",
				"oldName": oldName,
				"newName": newName})
		_, err = stmt.Exec(newName, oldName)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
