@echo off

@REM cd scripts/windows
@REM call clean.bat debug

set mode="release"

rd /s /q "build/%mode%"
md "build/%mode%"
