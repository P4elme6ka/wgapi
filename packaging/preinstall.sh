#!/usr/bin/env bash
set -eu

if systemctl status wgapi &> /dev/null; then
    systemctl stop wgapi.service
    systemctl disable wgapi.service
fi
