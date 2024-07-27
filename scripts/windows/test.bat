@echo off 

cd src/kernel
call go test duck/kernel/server/... -v
call go test duck/kernel/persistence/... -v
