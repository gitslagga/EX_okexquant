#!/bin/bash

if [ $# -eq 0 ];then
    echo "command: ./pack.sh [prd|test|dev]"
    exit
fi
pack_env="$1"

#update
git pull

#build
go build .

#pack
rm -rf release
mkdir -v release

mkdir release/config
cp -v EX_okexquant release/EX_okexquant_9801
cp -v run.sh release/
if [ "$pack_env" == "prd" ];then
    cp -rv config/config_prd.toml release/config/config.toml
elif [ "$pack_env" == "test" ];then
    cp -rv config/config_test.toml release/config/config.toml
elif [ "$pack_env" == "dev" ];then
    cp -rv config/config_dev.toml release/config/config.toml
fi

dd=$(date +%Y%m%d%H%M%S)
cd release/
commitid=`git log --pretty=format:"%h" -1`
tar -czvf okexquant-$pack_env-$commitid.tar.gz *
cd ../
mv release/okexquant-$pack_env-$commitid.tar.gz .
rm -rvf release
echo okexquant-$pack_env-$commitid.tar.gz
#sz okexquant-$pack_env-$commitid.tar.gz
