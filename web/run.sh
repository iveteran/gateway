#!/bin/sh

if [ ! -e ./fgw ]; then
  make
fi

./fgw -c ../conf/fgw.toml
