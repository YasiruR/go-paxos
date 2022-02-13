#!/bin/bash

for node in $(cat leaders.txt); do
  echo "${node}/internal/terminate"
  curl -X POST --retry 3 "http://${node}/internal/terminate"
done
rm leaders.txt

for node in $(cat replicas.txt); do
  curl -X POST --retry 3 "http://${node}/internal/terminate"
done
rm replicas.txt