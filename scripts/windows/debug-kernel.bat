@echo off

@REM cd scripts/windows
@REM call clean.bat debug

set mode="debug"

rd /s /q "build/%mode%"
md "build/%mode%"

set CGO_ENABLED=1 
cd src/kernel
call go build

cd ../..

cp src/kernel/kernel.exe build/%mode%/kernel.exe

cd build/debug
call kernel.exe -mode dev -port 9000 -dbPath ./duck.db
