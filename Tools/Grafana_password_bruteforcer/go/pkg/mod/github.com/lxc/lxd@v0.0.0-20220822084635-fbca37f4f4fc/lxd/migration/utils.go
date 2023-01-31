package migration

// IndexHeaderVersion version of the index header to be sent/recv.
const IndexHeaderVersion uint32 = 1

// BTRFSFeatureMigrationHeader indicates a migration header will be sent/recv in data channel after index header.
const BTRFSFeatureMigrationHeader = "migration_header"

// BTRFSFeatureSubvolumes indicates migration can send/recv subvolumes.
const BTRFSFeatureSubvolumes = "header_subvolumes"

// BTRFSFeatureSubvolumeUUIDs indicates that the header will include subvolume UUIDs.
const BTRFSFeatureSubvolumeUUIDs = "header_subvolume_uuids"

// ZFSFeatureMigrationHeader indicates a migration header will be sent/recv in data channel after index header.
const ZFSFeatureMigrationHeader = "migration_header"

// GetRsyncFeaturesSlice returns a slice of strings representing the supported RSYNC features.
func (m *MigrationHeader) GetRsyncFeaturesSlice() []string {
	features := []string{}
	if m == nil {
		return features
	}

	if m.RsyncFeatures != nil {
		if m.RsyncFeatures.Xattrs != nil && *m.RsyncFeatures.Xattrs {
			features = append(features, "xattrs")
		}

		if m.RsyncFeatures.Delete != nil && *m.RsyncFeatures.Delete {
			features = append(features, "delete")
		}

		if m.RsyncFeatures.Compress != nil && *m.RsyncFeatures.Compress {
			features = append(features, "compress")
		}

		if m.RsyncFeatures.Bidirectional != nil && *m.RsyncFeatures.Bidirectional {
			features = append(features, "bidirectional")
		}
	}

	return features
}

// GetZfsFeaturesSlice returns a slice of strings representing the supported ZFS features.
func (m *MigrationHeader) GetZfsFeaturesSlice() []string {
	features := []string{}
	if m == nil {
		return features
	}

	if m.ZfsFeatures != nil {
		if m.ZfsFeatures.Compress != nil && *m.ZfsFeatures.Compress {
			features = append(features, "compress")
		}

		if m.ZfsFeatures.MigrationHeader != nil && *m.ZfsFeatures.MigrationHeader {
			features = append(features, ZFSFeatureMigrationHeader)
		}
	}

	return features
}

// GetBtrfsFeaturesSlice returns a slice of strings representing the supported BTRFS features.
func (m *MigrationHeader) GetBtrfsFeaturesSlice() []string {
	features := []string{}
	if m == nil {
		return features
	}

	if m.BtrfsFeatures != nil {
		if m.BtrfsFeatures.MigrationHeader != nil && *m.BtrfsFeatures.MigrationHeader {
			features = append(features, BTRFSFeatureMigrationHeader)
		}

		if m.BtrfsFeatures.HeaderSubvolumes != nil && *m.BtrfsFeatures.HeaderSubvolumes {
			features = append(features, BTRFSFeatureSubvolumes)
		}

		if m.BtrfsFeatures.HeaderSubvolumeUuids != nil && *m.BtrfsFeatures.HeaderSubvolumeUuids {
			features = append(features, BTRFSFeatureSubvolumeUUIDs)
		}
	}

	return features
}
