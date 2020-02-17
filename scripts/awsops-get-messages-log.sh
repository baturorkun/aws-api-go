#!/usr/bin/env bash

set -x
source .env

if [[ -z $LINES ]]; then
    echo "<ERROR> Missing lines parameter"
    exit
fi

if (( $LINES > 10000 )); then
    echo "<ERROR> Max lines must be lower than 1000"
    exit
fi

echo "<SUCCESS>"

ssh -o StrictHostKeyChecking=no  -i $SERVER_SSHKEY centos@${PUBLIC_IP} "sudo tail -n $LINES /var/log/messages"
