#!/bin/bash

for i in {10..14}; do
  ifname="eth$i"
  echo "Starting $ifname"
  ip link add $ifname type dummy
  ip link set $ifname up
done
