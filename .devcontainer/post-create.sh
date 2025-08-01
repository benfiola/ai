#!/bin/bash -e
apt -y update
apt -y install curl git libarchive-tools make

go mod download
BIN=/usr/bin make install-tools