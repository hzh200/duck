@echo off

@REM cd scripts/windows
@REM call clean.bat debug

set mode=debug
set target_dir=build/%mode%

rd /s /q %target_dir%
md %target_dir%

cp src/app/index.html %target_dir%/
call yarn tailwindcss -i src/app/interfaces/styles/globals.css -o %target_dir%/globals.css
call yarn webpack --config webpack.config.%mode%.js

set CGO_ENABLED=1 

cd src/kernel
call go build

cd ../..

if not exist src/kernel/kernel.exe exit 1

mv src/kernel/kernel.exe %target_dir%/kernel.exe

call yarn tsc src/utils/preload.ts --outdir %target_dir%

call yarn electron . --enable-logging --inspect
