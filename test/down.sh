#!/usr/bin/env bash
set -ex

DIR=$PWD
docker image rm -f tape

cd ./fabric-samples/test-network
echo y |  ./real-network.sh down -i 2.2
