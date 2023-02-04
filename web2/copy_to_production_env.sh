#!/bin/sh
echo "copying ..."
scp -P24 fgw fimatrix@matrixworks.cn:~/programs/
scp -P24 ../conf/fgw.toml fimatrix@matrixworks.cn:~/programs/
echo "done"
