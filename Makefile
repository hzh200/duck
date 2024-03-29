


debug: 
	bash scripts/debug.sh

de:
	export NODE_ENV=development
	rm -rf build/debug
	mkdir -p build/debug
	cp src/app/index.html build/debug/
	yarn webpack --config webpack.config.debug.js
	go build -C src/kernel -o ../../build/debug/kernel
	yarn electron . --enable-logging --inspect --start-mode=normal