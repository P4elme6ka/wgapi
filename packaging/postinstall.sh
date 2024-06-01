#!/usr/bin/env bash
set -eu

adduser --system wgrest --home /var/lib/wgapi

systemctl enable "/etc/systemd/system/wgapi.service"
