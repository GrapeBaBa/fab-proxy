#!/usr/bin/env bash
set -ex

DIR=$PWD

cd test-network/
echo y |  ./real-network.sh down -i 2.2
echo y |  ./real-network.sh up createChannel -i 2.2