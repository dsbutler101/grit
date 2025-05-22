# Systemd Containers internal module

Generates [systemd](https://systemd.io/) service files and corresponding
[`systemctl`](https://www.freedesktop.org/software/systemd/man/latest/systemctl.html)
commands to run containers on host boot.

This can be used with [cloud-init](https://cloudinit.readthedocs.io/en/latest/index.html)
directly for `write_files` and `runcmd` modules to provide VM hosts
a systemd managed set of containers.
