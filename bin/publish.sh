#!/usr/bin/env bash
SCRIPT_DIR="$(dirname "$0")"
cd ${SCRIPT_DIR}/..

rm -rf dist

goreleaser release --snapshot
