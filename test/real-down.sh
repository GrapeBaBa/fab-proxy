#!/usr/bin/env bash
set -ex

DIR=$PWD

cd test-network
echo y |  ./network.sh down -i 2.2