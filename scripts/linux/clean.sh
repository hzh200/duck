#! /bin/bash

if [ $# == 0 ];
then
        exit
fi

mode=$1

if [ $mode != release ] && [ $mode != debug ]
then
        exit
fi

mkdir -p build/$mode
rm -rf build/$mode/*
