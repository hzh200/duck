#! /bin/bash

mode=debug

source ./clean.sh $mode

CGO_ENABLED=1 go build -C src/kernel -o ../../build/$mode/kernel
./build/debug/kernel -mode dev -port 9000 -dbPath ./build/debug/duck.db
