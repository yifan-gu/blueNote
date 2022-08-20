#!/usr/bin/env bash

function db_start() {
    mongod --config /opt/homebrew/etc/mongod.conf --fork
}

function db_shutdown() {
    mongosh --quiet --norc --eval "db.adminCommand({"shutdown":1})"
}

function usage() {
    echo "Usage: $0 start|shutdown"
}

if [ "$#" -ne 1 ]; then
    usage
    exit 1
fi

if [ "$1" == "start" ]; then
    db_start
elif [ "$1" == "shutdown" ]; then
    db_shutdown
fi
