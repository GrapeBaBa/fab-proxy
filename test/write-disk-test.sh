#!/usr/bin/env bash
set -ex

DIR=$PWD
docker image rm -f tape
docker build -t tape:latest .

case $1 in
 2_2)
#    curl -vsS https://raw.githubusercontent.com/hyperledger/fabric/release-2.2/scripts/bootstrap.sh | bash
    cd ./fabric-samples/test-network
#    echo y |  ./network.sh down -i 2.2
#    echo y |  ./network.sh up createChannel -i 2.2
#    cp -r organizations "$DIR"

    CONFIG_FILE=/config/test/config20org1andorg2.yaml

    if [ $2 == "ORLogic" ]; then
      CONFIG_FILE=/config/test/config20selectendorser.yaml
      ARGS=(-ccep "OR('Org1.member','Org2.member')")
    fi

#    echo y |  ./network.sh deployCC "${ARGS[@]}"
    ;;
esac

cd "$DIR"
docker run  -e TAPE_LOGLEVEL=info --network host -v $PWD:/config tape tape -c $CONFIG_FILE -n 400000 -p diskWrite