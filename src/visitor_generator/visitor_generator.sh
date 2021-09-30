#!/bin/bash

DIR=$(git rev-parse --show-toplevel)

function preprocess() {
    local __return_value_1=$2
    local __return_value_2=$3

    local TYPE=$(echo $1 | cut -d "." -f2)
    local PKG=$(echo $1 | cut -d "." -f1 | awk '{ print toupper(substr($0, 1, 1)) substr($0, 2) }')

    eval $__return_value_1=$TYPE
    eval $__return_value_2=$PKG
}

preprocess $1 TYPE_1 PKG_1
preprocess $2 TYPE_2 PKG_2

../../lib/wacc-genny -in="$DIR/src/ast/visitor.go" -out="$DIR/src/ast/${TYPE_1}_${TYPE_2}_visitor.go" gen "Something=$1 Another=$2 Ctx=$3"
../../lib/wacc-genny -in="$DIR/src/ast/acceptor.go" -out="$DIR/src/ast/${TYPE_1}_${TYPE_2}_acceptor.go" gen "Something=$1 Another=$2 Ctx=$3"
sed -i -e "s/${PKG_1}${TYPE_1}/$TYPE_1/g" "$DIR/src/ast/${TYPE_1}_${TYPE_2}_visitor.go" "$DIR/src/ast/${TYPE_1}_${TYPE_2}_acceptor.go"
sed -i -e "s/${PKG_2}${TYPE_2}/$TYPE_2/g" "$DIR/src/ast/${TYPE_1}_${TYPE_2}_visitor.go" "$DIR/src/ast/${TYPE_1}_${TYPE_2}_acceptor.go"
../../lib/wacc-goimports -w ../ast