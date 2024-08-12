@echo off

@REM cd scripts/windows
@REM call clean.bat debug

set mode=release
set target_dir=build/%mode%

rd /s /q %target_dir%
md %target_dir%
