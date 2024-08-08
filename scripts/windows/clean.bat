@echo off

if "%1" == "" exit

set mode=%1
set target_dir=build/%mode%

if not "%mode%" == "release" if not "%mode%" == "debug" exit

rd /s /q %target_dir%
md %target_dir%
