#!/bin/sh

s=$(git pull)
if [ "$s" = "Already up to date." ]; then
    echo 'git is already up to date'
    if [ ! -z "$CRONDEPLOY" ]; then
        exit
    fi
elif [ "$?" -ne 0 ]; then
    echo "couldn't git pull:"
    echo "$s"
fi

docker build . -t "usva:local"
docker compose down
docker compose up -d
