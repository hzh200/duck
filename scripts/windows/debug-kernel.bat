@echo off

@REM cd scripts/windows
@REM call clean.bat debug

set mode="debug"
set target_dir=build/%mode%

rd /s /q %target_dir%
md %target_dir%

set CGO_ENABLED=1 
cd src/kernel
call go build

cd ../..

if not exist src/kernel/kernel.exe exit 1

mv src/kernel/kernel.exe %target_dir%/kernel.exe

cd build/debug
call kernel.exe -mode dev -port 9000 -dbPath ./duck.db
