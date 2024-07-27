@echo off

if "%1" == "" exit

set mode=%1

if not "%mode%" == "release" if not "%mode%" == "debug" exit

rd /s /q "build/%mode%"
md "build/%mode%"
