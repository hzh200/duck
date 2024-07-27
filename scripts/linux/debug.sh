#! /bin/bash

mode=debug

source ./clean.sh $mode

cp src/app/index.html build/$mode/
yarn tailwindcss -i src/app/interfaces/styles/globals.css -o build/$mode/globals.css
yarn webpack --config webpack.config.$mode.js
CGO_ENABLED=1 go build -C src/kernel -o ../../build/$mode/kernel
yarn tsc src/utils/preload.ts --outdir build/$mode

export NODE_ENV=debug 
yarn electron . --enable-logging --inspect
