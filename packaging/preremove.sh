#!/usr/bin/env bash
set -eu

systemctl stop wgapi.service || true
systemctl disable wgapi.service || true
