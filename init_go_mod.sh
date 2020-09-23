#!/bin/sh
go mod init matrix.works/fmx-gateway

echo -en "\nreplace matrix.works/fmx-common => ../fmx-common" >> go.mod
