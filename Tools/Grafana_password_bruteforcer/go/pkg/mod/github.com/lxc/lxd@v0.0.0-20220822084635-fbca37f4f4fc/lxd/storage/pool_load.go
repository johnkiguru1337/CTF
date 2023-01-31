package storage

import (
	"fmt"

	"github.com/lxc/lxd/lxd/instance"
	"github.com/lxc/lxd/lxd/project"
	"github.com/lxc/lxd/lxd/response"
	"github.com/lxc/lxd/lxd/state"
	"github.com/lxc/lxd/lxd/storage/drivers"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/logger"
)

// PoolIDTemporary is used to indicate a temporary pool instance that is not in the database.
const PoolIDTemporary = -1

// volIDFuncMake returns a function that can be supplied to the underlying storage drivers allowing
// them to lookup the volume ID for a specific volume type and volume name. This function is tied
// to the Pool ID that it is generated for, meaning the storage drivers do not need to know the ID
// of the pool they belong to, or do they need access to the database.
func volIDFuncMake(state *state.State, poolID int64) func(volType drivers.VolumeType, volName string) (int64, error) {
	// Return a function to retrieve a volume ID for a volume Name for use in driver.
	return func(volType drivers.VolumeType, volName string) (int64, error) {
		volTypeID, err := VolumeTypeToDBType(volType)
		if err != nil {
			return -1, err
		}

		// It is possible for the project name to be encoded into the volume name in the
		// format <project>_<volume>. However not all volume types currently use this
		// encoding format, so if there is no underscore in the volume name then we assume
		// the project is default.
		projectName := project.Default

		// Currently only Containers, VMs and custom volumes support project level volumes.
		// This means that other volume types may have underscores in their names that don't
		// indicate the project name.
		if volType == drivers.VolumeTypeContainer || volType == drivers.VolumeTypeVM {
			projectName, volName = project.InstanceParts(volName)
		} else if volType == drivers.VolumeTypeCustom {
			projectName, volName = project.StorageVolumeParts(volName)
		}

		volID, _, err := state.DB.Cluster.GetLocalStoragePoolVolume(projectName, volName, volTypeID, poolID)
		if err != nil {
			if response.IsNotFoundError(err) {
				return -1, fmt.Errorf("Failed to get volume ID for project %q, volume %q, type %q: Volume doesn't exist", projectName, volName, volType)
			}

			return -1, err
		}

		return volID, nil
	}
}

// commonRules returns a set of common validators.
func commonRules() *drivers.Validators {
	return &drivers.Validators{
		PoolRules:   validatePoolCommonRules,
		VolumeRules: validateVolumeCommonRules,
	}
}

// NewTemporary instantiates a temporary pool from config supplied and returns a Pool interface.
// Not all functionality will be available due to the lack of Pool ID.
// If the pool's driver is not recognised then drivers.ErrUnknownDriver is returned.
func NewTemporary(state *state.State, info *api.StoragePool) (Pool, error) {
	// Handle mock requests.
	if state.OS.MockMode {
		pool := mockBackend{}
		pool.name = info.Name
		pool.state = state
		pool.logger = logger.AddContext(logger.Log, logger.Ctx{"driver": "mock", "pool": pool.name})
		driver, err := drivers.Load(state, "mock", "", nil, pool.logger, nil, nil)
		if err != nil {
			return nil, err
		}

		pool.driver = driver

		return &pool, nil
	}

	var poolID int64 = PoolIDTemporary // Temporary as not in DB. Not all functionality will be available.

	// Ensure a config map exists.
	if info.Config == nil {
		info.Config = map[string]string{}
	}

	logger := logger.AddContext(logger.Log, logger.Ctx{"driver": info.Driver, "pool": info.Name})

	// Load the storage driver.
	driver, err := drivers.Load(state, info.Driver, info.Name, info.Config, logger, volIDFuncMake(state, poolID), commonRules())
	if err != nil {
		return nil, err
	}

	// Setup the pool struct.
	pool := lxdBackend{}
	pool.driver = driver
	pool.id = poolID
	pool.db = *info
	pool.name = info.Name
	pool.state = state
	pool.logger = logger
	pool.nodes = nil // TODO support clustering.

	return &pool, nil
}

// LoadByType loads a network by driver type.
func LoadByType(state *state.State, driverType string) (Type, error) {
	logger := logger.AddContext(logger.Log, logger.Ctx{"driver": driverType})

	driver, err := drivers.Load(state, driverType, "", nil, logger, nil, commonRules())
	if err != nil {
		return nil, err
	}

	// Setup the pool struct.
	pool := lxdBackend{}
	pool.state = state
	pool.driver = driver
	pool.id = PoolIDTemporary
	pool.logger = logger

	return &pool, nil
}

// LoadByName retrieves the pool from the database by its name and returns a Pool interface.
// If the pool's driver is not recognised then drivers.ErrUnknownDriver is returned.
func LoadByName(state *state.State, name string) (Pool, error) {
	// Handle mock requests.
	if state.OS.MockMode {
		pool := mockBackend{}
		pool.name = name
		pool.state = state
		pool.logger = logger.AddContext(logger.Log, logger.Ctx{"driver": "mock", "pool": pool.name})
		driver, err := drivers.Load(state, "mock", "", nil, pool.logger, nil, nil)
		if err != nil {
			return nil, err
		}

		pool.driver = driver

		return &pool, nil
	}

	// Load the database record.
	poolID, dbPool, poolNodes, err := state.DB.Cluster.GetStoragePoolInAnyState(name)
	if err != nil {
		return nil, err
	}

	// Ensure a config map exists.
	if dbPool.Config == nil {
		dbPool.Config = map[string]string{}
	}

	logger := logger.AddContext(logger.Log, logger.Ctx{"driver": dbPool.Driver, "pool": dbPool.Name})

	// Load the storage driver.
	driver, err := drivers.Load(state, dbPool.Driver, dbPool.Name, dbPool.Config, logger, volIDFuncMake(state, poolID), commonRules())
	if err != nil {
		return nil, err
	}

	// Setup the pool struct.
	pool := lxdBackend{}
	pool.driver = driver
	pool.id = poolID
	pool.db = *dbPool
	pool.name = dbPool.Name
	pool.state = state
	pool.logger = logger
	pool.nodes = poolNodes

	return &pool, nil
}

// LoadByInstance retrieves the pool from the database using the instance's pool.
// If the pool's driver is not recognised then drivers.ErrUnknownDriver is returned. If the pool's
// driver does not support the instance's type then drivers.ErrNotSupported is returned.
func LoadByInstance(s *state.State, inst instance.Instance) (Pool, error) {
	poolName, err := inst.StoragePool()
	if err != nil {
		return nil, fmt.Errorf("Failed getting instance storage pool name: %w", err)
	}

	pool, err := LoadByName(s, poolName)
	if err != nil {
		return nil, fmt.Errorf("Failed loading storage pool %q: %w", poolName, err)
	}

	volType, err := InstanceTypeToVolumeType(inst.Type())
	if err != nil {
		return nil, err
	}

	for _, supportedType := range pool.Driver().Info().VolumeTypes {
		if supportedType == volType {
			return pool, nil
		}
	}

	// Return drivers not supported error for consistency with predefined errors returned by
	// LoadByName (which can return drivers.ErrUnknownDriver).
	return nil, drivers.ErrNotSupported
}

// IsAvailable checks if a pool is available.
func IsAvailable(poolName string) bool {
	unavailablePoolsMu.Lock()
	defer unavailablePoolsMu.Unlock()

	_, found := unavailablePools[poolName]
	return !found
}

// Patch applies specified patch to all storage pools.
// All storage pools must be available locally before any storage pools are patched.
func Patch(s *state.State, patchName string) error {
	unavailablePoolsMu.Lock()

	if len(unavailablePools) > 0 {
		unavailablePoolNames := make([]string, 0, len(unavailablePools))
		for unavailablePoolName := range unavailablePools {
			unavailablePoolNames = append(unavailablePoolNames, unavailablePoolName)
		}

		unavailablePoolsMu.Unlock()
		return fmt.Errorf("Unvailable storage pools: %v", unavailablePoolNames)
	}

	unavailablePoolsMu.Unlock()

	// Load all the pools.
	pools, err := s.DB.Cluster.GetStoragePoolNames()
	if err != nil {
		if response.IsNotFoundError(err) {
			return nil
		}

		return fmt.Errorf("Failed loading storage pool names: %w", err)
	}

	for _, poolName := range pools {
		pool, err := LoadByName(s, poolName)
		if err != nil {
			return fmt.Errorf("Failed loading storage pool %q: %w", poolName, err)
		}

		err = pool.ApplyPatch(patchName)
		if err != nil {
			return fmt.Errorf("Failed applying patch to pool %q: %w", poolName, err)
		}
	}

	return nil
}
