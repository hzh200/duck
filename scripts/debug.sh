#! /bin/bash
rm -rf build/debug
mkdir -p build/debug
cp src/app/index.html build/debug/
yarn tailwindcss -i src/app/interfaces/styles/globals.css -o build/debug/globals.css
yarn webpack --config webpack.config.debug.js
CGO_ENABLED=1 go build -C src/kernel -o ../../build/debug/kernel
yarn tsc src/utils/preload.ts --outdir build/debug
export NODE_ENV=development 
yarn electron . --enable-logging --inspect
