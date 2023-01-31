package config

import (
	"github.com/lxc/lxd/shared/api"
)

// Config represents the config of a backup that can be stored in a backup.yaml file (or embedded in index.yaml).
type Config struct {
	Container       *api.Instance                `yaml:"container,omitempty"` // Used by VM backups too.
	Snapshots       []*api.InstanceSnapshot      `yaml:"snapshots,omitempty"`
	Pool            *api.StoragePool             `yaml:"pool,omitempty"`
	Volume          *api.StorageVolume           `yaml:"volume,omitempty"`
	VolumeSnapshots []*api.StorageVolumeSnapshot `yaml:"volume_snapshots,omitempty"`
}
