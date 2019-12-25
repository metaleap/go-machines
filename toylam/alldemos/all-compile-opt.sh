#!/bin/sh

names="comptryouts env fac fib hello json readln"

for name in $names
do
    tl2atem appdemo.$name.tl ../../../atmo/atem/tmpdummies/
    atem_opt < ../../../atmo/atem/tmpdummies/appdemo.$name.json > ../../../atmo/atem/tmpdummies/appdemo.$name.opt.json
done
