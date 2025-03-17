#!/usr/bin/env bash

#删除同名镜像
docker rmi wsdev/server-gstep:1.0 --force

#-t: 镜像名称
docker build -f ./Dockerfile -t wsdev/server-gstep:1.0 .