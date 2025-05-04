#!/usr/bin/env bash

while getopts "hrfds" OPTS; do
case $OPTS in
h)
echo " -r       : docker-compose up --build"
echo " -f       : docker-compose down && docker-compose up --build --force-recreate"
echo " -d       : docker-compose down"
echo " -s       : docker-compose up"
echo " -h       : this msg"
echo " NO flags : help"
;;
r)
docker compose up --build
;;
f)
docker compose down && docker compose up --build --force-recreate
;;
d)
docker compose down
;;
s)
docker compose up
;;
\?)
echo " -r       : docker-compose up --build"
echo " -f       : docker-compose down && docker-compose up --build --force-recreate"
echo " -d       : docker-compose down"
echo " -s       : docker-compose up"
echo " NO flags : docker-compose up --build"
;;
esac
done