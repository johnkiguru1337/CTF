package cluster

// DO NOT EDIT BY HAND
//
// This code was generated by the schema.DotGo function. If you need to
// modify the database schema, please add a new schema update to update.go
// and the run 'make update-schema'.
const freshSchema = `
CREATE TABLE certificates (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    fingerprint TEXT NOT NULL,
    type INTEGER NOT NULL,
    name TEXT NOT NULL,
    certificate TEXT NOT NULL,
    restricted INTEGER NOT NULL DEFAULT 0,
    UNIQUE (fingerprint)
);
CREATE TABLE "certificates_projects" (
	certificate_id INTEGER NOT NULL,
	project_id INTEGER NOT NULL,
	FOREIGN KEY (certificate_id) REFERENCES certificates (id) ON DELETE CASCADE,
	FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE,
	UNIQUE (certificate_id, project_id)
);
CREATE TABLE "cluster_groups" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    UNIQUE (name)
);
CREATE TABLE config (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    UNIQUE (key)
);
CREATE TABLE "images" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    fingerprint TEXT NOT NULL,
    filename TEXT NOT NULL,
    size INTEGER NOT NULL,
    public INTEGER NOT NULL DEFAULT 0,
    architecture INTEGER NOT NULL,
    creation_date DATETIME,
    expiry_date DATETIME,
    upload_date DATETIME NOT NULL,
    cached INTEGER NOT NULL DEFAULT 0,
    last_use_date DATETIME,
    auto_update INTEGER NOT NULL DEFAULT 0,
    project_id INTEGER NOT NULL,
    type INTEGER NOT NULL DEFAULT 0,
    UNIQUE (project_id, fingerprint),
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE TABLE "images_aliases" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    image_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    UNIQUE (project_id, name),
    FOREIGN KEY (image_id) REFERENCES "images" (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE INDEX images_aliases_project_id_idx ON images_aliases (project_id);
CREATE TABLE "images_nodes" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    image_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    UNIQUE (image_id, node_id),
    FOREIGN KEY (image_id) REFERENCES "images" (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE
);
CREATE TABLE "images_profiles" (
	image_id INTEGER NOT NULL,
	profile_id INTEGER NOT NULL,
	FOREIGN KEY (image_id) REFERENCES "images" (id) ON DELETE CASCADE,
	FOREIGN KEY (profile_id) REFERENCES "profiles" (id) ON DELETE CASCADE,
	UNIQUE (image_id, profile_id)
);
CREATE INDEX images_project_id_idx ON images (project_id);
CREATE TABLE "images_properties" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    image_id INTEGER NOT NULL,
    type INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    FOREIGN KEY (image_id) REFERENCES "images" (id) ON DELETE CASCADE
);
CREATE TABLE "images_source" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    image_id INTEGER NOT NULL,
    server TEXT NOT NULL,
    protocol INTEGER NOT NULL,
    certificate TEXT NOT NULL,
    alias TEXT NOT NULL,
    FOREIGN KEY (image_id) REFERENCES "images" (id) ON DELETE CASCADE
);
CREATE TABLE "instances" (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    node_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    architecture INTEGER NOT NULL,
    type INTEGER NOT NULL,
    ephemeral INTEGER NOT NULL DEFAULT 0,
    creation_date DATETIME NOT NULL DEFAULT 0,
    stateful INTEGER NOT NULL DEFAULT 0,
    last_use_date DATETIME,
    description TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    expiry_date DATETIME,
    UNIQUE (project_id, name),
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE TABLE "instances_backups" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    instance_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    creation_date DATETIME,
    expiry_date DATETIME,
    container_only INTEGER NOT NULL default 0,
    optimized_storage INTEGER NOT NULL default 0,
    FOREIGN KEY (instance_id) REFERENCES "instances" (id) ON DELETE CASCADE,
    UNIQUE (instance_id, name)
);
CREATE TABLE "instances_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    instance_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (instance_id) REFERENCES "instances" (id) ON DELETE CASCADE,
    UNIQUE (instance_id, key)
);
CREATE TABLE "instances_devices" (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    instance_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type INTEGER NOT NULL default 0,
    FOREIGN KEY (instance_id) REFERENCES "instances" (id) ON DELETE CASCADE,
    UNIQUE (instance_id, name)
);
CREATE TABLE "instances_devices_config" (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    instance_device_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (instance_device_id) REFERENCES "instances_devices" (id) ON DELETE CASCADE,
    UNIQUE (instance_device_id, key)
);
CREATE INDEX instances_node_id_idx ON instances (node_id);
CREATE TABLE "instances_profiles" (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    instance_id INTEGER NOT NULL,
    profile_id INTEGER NOT NULL,
    apply_order INTEGER NOT NULL default 0,
    UNIQUE (instance_id, profile_id),
    FOREIGN KEY (instance_id) REFERENCES "instances" (id) ON DELETE CASCADE,
    FOREIGN KEY (profile_id) REFERENCES "profiles"(id) ON DELETE CASCADE
);
CREATE INDEX instances_project_id_and_name_idx ON instances (project_id,
    name);
CREATE INDEX instances_project_id_and_node_id_and_name_idx ON instances (project_id,
    node_id,
    name);
CREATE INDEX instances_project_id_and_node_id_idx ON instances (project_id,
    node_id);
CREATE INDEX instances_project_id_idx ON instances (project_id);
CREATE TABLE "instances_snapshots" (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    instance_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    creation_date DATETIME NOT NULL DEFAULT 0,
    stateful INTEGER NOT NULL DEFAULT 0,
    description TEXT NOT NULL,
    expiry_date DATETIME,
    UNIQUE (instance_id, name),
    FOREIGN KEY (instance_id) REFERENCES "instances" (id) ON DELETE CASCADE
);
CREATE TABLE "instances_snapshots_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    instance_snapshot_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (instance_snapshot_id) REFERENCES "instances_snapshots" (id) ON DELETE CASCADE,
    UNIQUE (instance_snapshot_id, key)
);
CREATE TABLE "instances_snapshots_devices" (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    instance_snapshot_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type INTEGER NOT NULL default 0,
    FOREIGN KEY (instance_snapshot_id) REFERENCES "instances_snapshots" (id) ON DELETE CASCADE,
    UNIQUE (instance_snapshot_id, name)
);
CREATE TABLE "instances_snapshots_devices_config" (
    id INTEGER primary key AUTOINCREMENT NOT NULL,
    instance_snapshot_device_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (instance_snapshot_device_id) REFERENCES "instances_snapshots_devices" (id) ON DELETE CASCADE,
    UNIQUE (instance_snapshot_device_id, key)
);
CREATE TABLE "networks" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    state INTEGER NOT NULL DEFAULT 0,
    type INTEGER NOT NULL DEFAULT 0,
    UNIQUE (project_id, name),
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_acls" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    ingress TEXT NOT NULL,
    egress TEXT NOT NULL,
    UNIQUE (project_id, name),
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_acls_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    network_acl_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE (network_acl_id, key),
    FOREIGN KEY (network_acl_id) REFERENCES "networks_acls" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    network_id INTEGER NOT NULL,
    node_id INTEGER,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE (network_id, node_id, key),
    FOREIGN KEY (network_id) REFERENCES "networks" (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_forwards" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_id INTEGER NOT NULL,
	node_id INTEGER,
	listen_address TEXT NOT NULL,
	description TEXT NOT NULL,
	ports TEXT NOT NULL,
	UNIQUE (network_id, node_id, listen_address),
	FOREIGN KEY (network_id) REFERENCES "networks" (id) ON DELETE CASCADE,
	FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_forwards_config" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_forward_id INTEGER NOT NULL,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	UNIQUE (network_forward_id, key),
	FOREIGN KEY (network_forward_id) REFERENCES "networks_forwards" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_load_balancers" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_id INTEGER NOT NULL,
	node_id INTEGER,
	listen_address TEXT NOT NULL,
	description TEXT NOT NULL,
	backends TEXT NOT NULL,
	ports TEXT NOT NULL,
	UNIQUE (network_id, node_id, listen_address),
	FOREIGN KEY (network_id) REFERENCES "networks" (id) ON DELETE CASCADE,
	FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_load_balancers_config" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_load_balancer_id INTEGER NOT NULL,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	UNIQUE (network_load_balancer_id, key),
	FOREIGN KEY (network_load_balancer_id) REFERENCES "networks_load_balancers" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_nodes" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    network_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    state INTEGER NOT NULL DEFAULT 0,
    UNIQUE (network_id, node_id),
    FOREIGN KEY (network_id) REFERENCES "networks" (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_peers" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	target_network_project TEXT NULL,
	target_network_name TEXT NULL,
	target_network_id INTEGER NULL,
	UNIQUE (network_id, name),
	UNIQUE (network_id, target_network_project, target_network_name),
	UNIQUE (network_id, target_network_id),
	FOREIGN KEY (network_id) REFERENCES "networks" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_peers_config" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_peer_id INTEGER NOT NULL,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	UNIQUE (network_peer_id, key),
	FOREIGN KEY (network_peer_id) REFERENCES "networks_peers" (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX networks_unique_network_id_node_id_key ON "networks_config" (network_id, IFNULL(node_id, -1), key);
CREATE TABLE "networks_zones" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	project_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	UNIQUE (name),
	FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE TABLE "networks_zones_config" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_zone_id INTEGER NOT NULL,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	UNIQUE (network_zone_id, key),
	FOREIGN KEY (network_zone_id) REFERENCES "networks_zones" (id) ON DELETE CASCADE
);
CREATE TABLE networks_zones_records (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_zone_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	entries TEXT NOT NULL,
	UNIQUE (name),
	FOREIGN KEY (network_zone_id) REFERENCES networks_zones (id) ON DELETE CASCADE
);
CREATE TABLE "networks_zones_records_config" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	network_zone_record_id INTEGER NOT NULL,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	UNIQUE (network_zone_record_id, key),
	FOREIGN KEY (network_zone_record_id) REFERENCES networks_zones_records (id) ON DELETE CASCADE
);
CREATE TABLE "nodes" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    address TEXT NOT NULL,
    schema INTEGER NOT NULL,
    api_extensions INTEGER NOT NULL,
    heartbeat DATETIME DEFAULT CURRENT_TIMESTAMP,
    state INTEGER NOT NULL DEFAULT 0,
    arch INTEGER NOT NULL DEFAULT 0 CHECK (arch > 0),
    failure_domain_id INTEGER DEFAULT NULL REFERENCES nodes_failure_domains (id) ON DELETE SET NULL,
    UNIQUE (name),
    UNIQUE (address)
);
CREATE TABLE "nodes_cluster_groups" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    node_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES cluster_groups (id) ON DELETE CASCADE,
    UNIQUE (node_id, group_id)
);
CREATE TABLE "nodes_config" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	node_id INTEGER NOT NULL,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE,
	UNIQUE (node_id, key)
);
CREATE TABLE nodes_failure_domains (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    UNIQUE (name)
);
CREATE TABLE "nodes_roles" (
    node_id INTEGER NOT NULL,
    role INTEGER NOT NULL,
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE,
    UNIQUE (node_id, role)
);
CREATE TABLE "operations" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    uuid TEXT NOT NULL,
    node_id TEXT NOT NULL,
    type INTEGER NOT NULL DEFAULT 0,
    project_id INTEGER,
    UNIQUE (uuid),
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE TABLE "profiles" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    UNIQUE (project_id, name),
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE TABLE "profiles_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    profile_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE (profile_id, key),
    FOREIGN KEY (profile_id) REFERENCES "profiles"(id) ON DELETE CASCADE
);
CREATE TABLE "profiles_devices" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    profile_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type INTEGER NOT NULL default 0,
    UNIQUE (profile_id, name),
    FOREIGN KEY (profile_id) REFERENCES "profiles" (id) ON DELETE CASCADE
);
CREATE TABLE "profiles_devices_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    profile_device_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE (profile_device_id, key),
    FOREIGN KEY (profile_device_id) REFERENCES "profiles_devices" (id) ON DELETE CASCADE
);
CREATE INDEX profiles_project_id_idx ON profiles (project_id);
CREATE TABLE "projects" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    UNIQUE (name)
);
CREATE TABLE "projects_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    project_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE,
    UNIQUE (project_id, key)
);
CREATE TABLE "storage_buckets" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name TEXT NOT NULL,
	storage_pool_id INTEGER NOT NULL,
	node_id INTEGER,
	description TEXT NOT NULL,
	project_id INTEGER NOT NULL,
	UNIQUE (node_id, name),
	FOREIGN KEY (storage_pool_id) REFERENCES "storage_pools" (id) ON DELETE CASCADE,
	FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE,
	FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE TABLE "storage_buckets_config" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	storage_bucket_id INTEGER NOT NULL,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	UNIQUE (storage_bucket_id, key),
	FOREIGN KEY (storage_bucket_id) REFERENCES "storage_buckets" (id) ON DELETE CASCADE
);
CREATE TABLE "storage_buckets_keys" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	storage_bucket_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	access_key TEXT NOT NULL,
	secret_key TEXT NOT NULL,
	role TEXT NOT NULL,
	UNIQUE (storage_bucket_id, name),
	FOREIGN KEY (storage_bucket_id) REFERENCES "storage_buckets" (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX storage_buckets_unique_storage_pool_id_node_id_name ON "storage_buckets" (storage_pool_id, IFNULL(node_id, -1), name);
CREATE TABLE "storage_pools" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    driver TEXT NOT NULL,
    description TEXT NOT NULL,
    state INTEGER NOT NULL DEFAULT 0,
    UNIQUE (name)
);
CREATE TABLE "storage_pools_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    storage_pool_id INTEGER NOT NULL,
    node_id INTEGER,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE (storage_pool_id, node_id, key),
    FOREIGN KEY (storage_pool_id) REFERENCES "storage_pools" (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE
);
CREATE TABLE "storage_pools_nodes" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    storage_pool_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    state INTEGER NOT NULL DEFAULT 0,
    UNIQUE (storage_pool_id, node_id),
    FOREIGN KEY (storage_pool_id) REFERENCES "storage_pools" (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX storage_pools_unique_storage_pool_id_node_id_key ON storage_pools_config (storage_pool_id, IFNULL(node_id, -1), key);
CREATE TABLE "storage_volumes" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    storage_pool_id INTEGER NOT NULL,
    node_id INTEGER,
    type INTEGER NOT NULL,
    description TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    content_type INTEGER NOT NULL DEFAULT 0,
    UNIQUE (storage_pool_id, node_id, project_id, name, type),
    FOREIGN KEY (storage_pool_id) REFERENCES "storage_pools" (id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES "nodes" (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE VIEW storage_volumes_all (
         id,
         name,
         storage_pool_id,
         node_id,
         type,
         description,
         project_id,
         content_type) AS
  SELECT id,
         name,
         storage_pool_id,
         node_id,
         type,
         description,
         project_id,
         content_type
    FROM storage_volumes UNION
  SELECT storage_volumes_snapshots.id,
         printf('%s/%s',
    storage_volumes.name,
    storage_volumes_snapshots.name),
         storage_volumes.storage_pool_id,
         storage_volumes.node_id,
         storage_volumes.type,
         storage_volumes_snapshots.description,
         storage_volumes.project_id,
         storage_volumes.content_type
    FROM storage_volumes
    JOIN storage_volumes_snapshots ON storage_volumes.id = storage_volumes_snapshots.storage_volume_id;
CREATE TABLE "storage_volumes_backups" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    storage_volume_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    creation_date DATETIME,
    expiry_date DATETIME,
    volume_only INTEGER NOT NULL default 0,
    optimized_storage INTEGER NOT NULL default 0,
    FOREIGN KEY (storage_volume_id) REFERENCES "storage_volumes" (id) ON DELETE CASCADE,
    UNIQUE (storage_volume_id, name)
);
CREATE TRIGGER storage_volumes_check_id
  BEFORE INSERT ON storage_volumes
  WHEN NEW.id IN (SELECT id FROM storage_volumes_snapshots)
  BEGIN
    SELECT RAISE(FAIL,
    "invalid ID");
  END;
CREATE TABLE "storage_volumes_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    storage_volume_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE (storage_volume_id, key),
    FOREIGN KEY (storage_volume_id) REFERENCES "storage_volumes" (id) ON DELETE CASCADE
);
CREATE TABLE "storage_volumes_snapshots" (
    id INTEGER NOT NULL,
    storage_volume_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    expiry_date DATETIME,
    UNIQUE (id),
    UNIQUE (storage_volume_id, name),
    FOREIGN KEY (storage_volume_id) REFERENCES "storage_volumes" (id) ON DELETE CASCADE
);
CREATE TRIGGER storage_volumes_snapshots_check_id
  BEFORE INSERT ON storage_volumes_snapshots
  WHEN NEW.id IN (SELECT id FROM storage_volumes)
  BEGIN
    SELECT RAISE(FAIL,
    "invalid ID");
  END;
CREATE TABLE "storage_volumes_snapshots_config" (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    storage_volume_snapshot_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY (storage_volume_snapshot_id) REFERENCES "storage_volumes_snapshots" (id) ON DELETE CASCADE,
    UNIQUE (storage_volume_snapshot_id, key)
);
CREATE UNIQUE INDEX storage_volumes_unique_storage_pool_id_node_id_project_id_name_type ON "storage_volumes" (storage_pool_id, IFNULL(node_id, -1), project_id, name, type);
CREATE TABLE "warnings" (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	node_id INTEGER,
	project_id INTEGER,
	entity_type_code INTEGER,
	entity_id INTEGER,
	uuid TEXT NOT NULL,
	type_code INTEGER NOT NULL,
	status INTEGER NOT NULL,
	first_seen_date DATETIME NOT NULL,
	last_seen_date DATETIME NOT NULL,
	updated_date DATETIME,
	last_message TEXT NOT NULL,
	count INTEGER NOT NULL,
	UNIQUE (uuid),
	FOREIGN KEY (node_id) REFERENCES "nodes"(id) ON DELETE CASCADE,
	FOREIGN KEY (project_id) REFERENCES "projects" (id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX warnings_unique_node_id_project_id_entity_type_code_entity_id_type_code ON warnings(IFNULL(node_id, -1), IFNULL(project_id, -1), entity_type_code, entity_id, type_code);

INSERT INTO schema (version, updated_at) VALUES (65, strftime("%s"))
`
