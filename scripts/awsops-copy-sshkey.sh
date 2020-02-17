#!/usr/bin/env bash

set -x
source .env

SSHKEY=$1
PUBLIC_IP=$2

KEYDATA=`cat data/keys/$SSHKEY`


#check exists
ssh -o StrictHostKeyChecking=no  -i $SERVER_SSHKEY centos@${PUBLIC_IP} "grep '$KEYDATA' .ssh/authorized_keys" >> /dev/null

if [[ $? -eq 0 ]]; then
    echo "Already exists"
    echo "<FAILED>"
    exit
fi

ssh -o StrictHostKeyChecking=no  -i $SERVER_SSHKEY centos@${PUBLIC_IP} "echo '$KEYDATA' >> .ssh/authorized_keys"

if [[ $? -ne 0 ]]; then
    echo "<FAILED>"
else
    echo "<SUCCESS>"
fi
