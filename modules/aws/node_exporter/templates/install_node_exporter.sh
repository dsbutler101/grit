#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

set -x

arch_suffix=$(uname -m | grep -q 'x86_64' && echo 'amd64' || echo 'arm64')

tmp_dir="${node_exporter_dir}/tmp"
mkdir -p "$${tmp_dir}"
node_exporter_tar_file_name="node_exporter.tar.gz"

cd "$${tmp_dir}"
curl -L "https://github.com/prometheus/node_exporter/releases/download/v${node_exporter_version}/node_exporter-${node_exporter_version}.linux-$${arch_suffix}.tar.gz" -o "$${node_exporter_tar_file_name}"
tar xvfz "$${node_exporter_tar_file_name}"
mv node_exporter*/* "${node_exporter_dir}"

rm -rf "$${tmp_dir}"

cat > /etc/systemd/system/node-exporter.service << EOF
[Install]
WantedBy=multi-user.target

[Unit]
Description=Node exporter

[Service]
Type=exec
ExecStart="${node_exporter_dir}"/node_exporter --web.listen-address=0.0.0.0:${node_exporter_port}
EOF

chmod 644 /etc/systemd/system/node-exporter.service
systemctl daemon-reload
systemctl enable node-exporter.service
systemctl start node-exporter.service
