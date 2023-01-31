#!/usr/bin/env bats

load helpers

function setup() {
	setup_busybox
}

function teardown() {
	teardown_bundle
}

@test "runc run [bind mount]" {
	update_config '	  .mounts += [{
					source: ".",
					destination: "/tmp/bind",
					options: ["bind"]
				}]
			| .process.args |= ["ls", "/tmp/bind/config.json"]'

	runc run test_busybox
	[ "$status" -eq 0 ]
	[[ "${lines[0]}" == *'/tmp/bind/config.json'* ]]
}

# https://github.com/opencontainers/runc/issues/2246
@test "runc run [ro tmpfs mount]" {
	update_config '	  .mounts += [{
					source: "tmpfs",
					destination: "/mnt",
					type: "tmpfs",
					options: ["ro", "nodev", "nosuid", "mode=755"]
				}]
			| .process.args |= ["grep", "^tmpfs /mnt", "/proc/mounts"]'

	runc run test_busybox
	[ "$status" -eq 0 ]
	[[ "${lines[0]}" == *'ro,'* ]]
}

# https://github.com/opencontainers/runc/issues/3248
@test "runc run [ro /dev mount]" {
	update_config '   .mounts |= map((select(.destination == "/dev") | .options += ["ro"]) // .)
			| .process.args |= ["grep", "^tmpfs /dev", "/proc/mounts"]'

	runc run test_busybox
	[ "$status" -eq 0 ]
	[[ "${lines[0]}" == *'ro,'* ]]
}

# https://github.com/opencontainers/runc/issues/2683
@test "runc run [tmpfs mount with absolute symlink]" {
	# in container, /conf -> /real/conf
	mkdir -p rootfs/real/conf
	ln -s /real/conf rootfs/conf
	update_config '	  .mounts += [{
					type: "tmpfs",
					source: "tmpfs",
					destination: "/conf/stack",
					options: ["ro", "nodev", "nosuid"]
				}]
			| .process.args |= ["true"]'
	runc run test_busybox
	[ "$status" -eq 0 ]
}
