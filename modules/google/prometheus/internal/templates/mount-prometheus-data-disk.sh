#!/usr/bin/env bash

set -eo pipefail
set -x

if ! blkid "${device_path}" | grep 'TYPE="ext4"'; then
  mkfs.ext4 "${device_path}"
fi

umount "${mount_path}" || echo "Nothing to umount from ${mount_path}"

fsck.ext4 -ftvy "${device_path}"
resize2fs "${device_path}"

mkdir -p "${mount_path}"
mount -t ext4 -o rw,nosuid,nodev,relatime "${device_path}" "${mount_path}"

uid=$(docker run --rm --entrypoint /bin/id ${prometheus_image} -u)
gid=$(docker run --rm --entrypoint /bin/id ${prometheus_image} -g)

mkdir -p "${prometheus_volume}"
chown "$${uid}:$${gid}" "${prometheus_volume}"
