test_storage_volume_attach() {
  # Check that we have a big enough range for this test
  if [ ! -e /etc/subuid ] && [ ! -e /etc/subgid ]; then
    UIDs=1000000000
    GIDs=1000000000
    UID_BASE=1000000
    GID_BASE=1000000
  else
    UIDs=0
    GIDs=0
    UID_BASE=0
    GID_BASE=0
    LARGEST_UIDs=0
    LARGEST_GIDs=0

    # shellcheck disable=SC2013
    for entry in $(grep ^root: /etc/subuid); do
      COUNT=$(echo "${entry}" | cut -d: -f3)
      UIDs=$((UIDs+COUNT))

      if [ "${COUNT}" -gt "${LARGEST_UIDs}" ]; then
        LARGEST_UIDs=${COUNT}
        UID_BASE=$(echo "${entry}" | cut -d: -f2)
      fi
    done

    # shellcheck disable=SC2013
    for entry in $(grep ^root: /etc/subgid); do
      COUNT=$(echo "${entry}" | cut -d: -f3)
      GIDs=$((GIDs+COUNT))

      if [ "${COUNT}" -gt "${LARGEST_GIDs}" ]; then
        LARGEST_GIDs=${COUNT}
        GID_BASE=$(echo "${entry}" | cut -d: -f2)
      fi
    done
  fi

  ensure_import_testimage

  # create storage volume
  lxc storage volume create "lxdtest-$(basename "${LXD_DIR}")" testvolume

  # create containers
  lxc launch testimage c1 -c security.privileged=true
  lxc launch testimage c2

  # Attach to a single privileged container
  lxc storage volume attach "lxdtest-$(basename "${LXD_DIR}")" testvolume c1 testvolume
  PATH_TO_CHECK="${LXD_DIR}/storage-pools/lxdtest-$(basename "${LXD_DIR}")/custom/default_testvolume"
  [ "$(stat -c %u:%g "${PATH_TO_CHECK}")" = "0:0" ]

  # make container unprivileged
  lxc config set c1 security.privileged false
  [ "$(stat -c %u:%g "${PATH_TO_CHECK}")" = "0:0" ]

  if [ "${UIDs}" -lt 500000 ] || [ "${GIDs}" -lt 500000 ]; then
    echo "==> SKIP: The storage volume attach test requires at least 500000 uids and gids"
    lxc rm -f c1 c2
    lxc storage volume delete "lxdtest-$(basename "${LXD_DIR}")" testvolume
    return
  fi

  # restart
  lxc restart --force c1
  [ "$(stat -c %u:%g "${PATH_TO_CHECK}")" = "${UID_BASE}:${GID_BASE}" ]

  # give container isolated id mapping
  lxc config set c1 security.idmap.isolated true
  [ "$(stat -c %u:%g "${PATH_TO_CHECK}")" = "${UID_BASE}:${GID_BASE}" ]

  # restart
  lxc restart --force c1

  # get new isolated base ids
  ISOLATED_UID_BASE="$(lxc exec c1 -- cat /proc/self/uid_map | awk '{print $2}')"
  ISOLATED_GID_BASE="$(lxc exec c1 -- cat /proc/self/gid_map | awk '{print $2}')"
  [ "$(stat -c %u:%g "${PATH_TO_CHECK}")" = "${ISOLATED_UID_BASE}:${ISOLATED_GID_BASE}" ]

  ! lxc storage volume attach "lxdtest-$(basename "${LXD_DIR}")" testvolume c2 testvolume || false

  # give container standard mapping
  lxc config set c1 security.idmap.isolated false
  [ "$(stat -c %u:%g "${PATH_TO_CHECK}")" = "${ISOLATED_UID_BASE}:${ISOLATED_GID_BASE}" ]

  # restart
  lxc restart --force c1
  [ "$(stat -c %u:%g "${PATH_TO_CHECK}")" = "${UID_BASE}:${GID_BASE}" ]

  # attach second container
  lxc storage volume attach "lxdtest-$(basename "${LXD_DIR}")" testvolume c2 testvolume

  # check that setting perms on the root of the custom volume persists after a reboot.
  lxc exec c2 -- stat -c '%a' /testvolume | grep 711
  lxc exec c2 -- chmod 0700 /testvolume
  lxc exec c2 -- stat -c '%a' /testvolume | grep 700
  lxc restart --force c2
  lxc exec c2 -- stat -c '%a' /testvolume | grep 700

  # delete containers
  lxc delete -f c1
  lxc delete -f c2
  lxc storage volume delete "lxdtest-$(basename "${LXD_DIR}")" testvolume
}
