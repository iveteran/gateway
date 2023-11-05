#!/bin/bash

TARGET=debian-12/home/fimatrix/programs/

echo "pushing ..."
lxc file push fgw $TARGET
lxc file push ../conf/fgw.toml $TARGET
echo "done"
