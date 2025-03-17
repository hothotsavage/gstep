#!/usr/bin/env bash
docker rm -f server-gstep

docker run -d \
-p 8216:8216 \
--name server-gstep \
wsdev/server-gstep:1.0