#!/bin/sh

i=0
while [ $i -lt 5000000 ]; do echo "$URL"; i=$(( i + 1 )); done | xargs -P 100 curl -Ss
