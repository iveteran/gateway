#!/bin/sh

TARGET=../../fmx-deployment/binaries

echo "copying ..."
cp -v fgw $TARGET
cp -v ../conf/fgw.toml $TARGET
echo "done"
