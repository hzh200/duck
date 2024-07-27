@echo off

@REM cd scripts/windows
@REM call clean.bat debug

set mode="debug"

rd /s /q "build/%mode%"
md "build/%mode%"

cp src/app/index.html build/%mode%/
call yarn tailwindcss -i src/app/interfaces/styles/globals.css -o build/%mode%/globals.css
call yarn webpack --config webpack.config.%mode%.js

set CGO_ENABLED=1 
call go build -C src/kernel.exe -o ../../build/%mode%/kernel.exe
call yarn tsc src/utils/preload.ts --outdir build/%mode%

set NODE_ENV=debug 
call yarn electron . --enable-logging --inspect
