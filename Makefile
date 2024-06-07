.PHONY: debug-clean
debug-clean:
	mkdir -p build/debug
	rm -rf build/debug/*

.PHONY: debug-build-kernel
debug-build-kernel: debug-clean
	CGO_ENABLED=1 go build -C src/kernel -o ../../build/debug/kernel

.PHONY: debug-build-app
debug-build-app: debug-clean
	cp src/app/index.html build/debug/
	yarn tailwindcss -i src/app/interfaces/styles/globals.css -o build/debug/globals.css
	yarn webpack --config webpack.config.debug.js
	yarn tsc src/utils/preload.ts --outdir build/debug

.PHONY: debug-build
debug-build: debug-clean debug-build-kernel debug-build-app

.PHONY: debug-run-kernel
debug-run-kernel: debug-build-kernel
	./build/debug/kernel -mode dev -port 9000 -dbPath ./build/debug/duck.db

.PHONY: debug-run
debug-run: debug-build
	export NODE_ENV=development 
	yarn electron . --enable-logging --inspect

.PHONY: test
test: 
	cd src/kernel && go test duck/kernel/server/...
	cd src/kernel && go test duck/kernel/persistence/...
