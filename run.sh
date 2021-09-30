#!bin/bash

if test -f "input.s"; then
    rm input.s
fi

./compile $1 || exit $?

inputPath=$(dirname $1)/input.s
mv $inputPath input.s

if test -f "input.s"; then
    arm-linux-gnueabi-gcc -g -o input -mcpu=arm1176jzf-s -mtune=arm1176jzf-s input.s
    qemu-arm  -L /usr/arm-linux-gnueabi input
fi
