#!/bin/bash

usage() {
    echo "Usage: $0 <height> <width> <fill(0/1)> <character>"
	exit -1
}

if (( $# < 4 )); then
    usage
fi

height=$1
width=$2
isFilling=$3
char=$4

for (( i = 0; i < height; i++ )); do
    for (( j = 0; j < width; j++ )); do
        if (( isFilling || i == 0 || j == 0 || i == (height-1) || j == (width-1) )); then
            echo -n "$char"
        else
            echo -n " "
        fi
    done
    echo
done