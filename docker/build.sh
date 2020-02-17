#! /bin/bash
set -x
source .env

########################### GITLAB CREDENTIALS ###############################
# git config --global credential.helper store && echo "https://$GIT_USERNAME:$GIT_TOKEN@$GIT_URL" > ~/.git-credentials
##########################DO NOT TOUCH THESE##################################

DISABLE_HOSTKEY_VERIFICATION='
Host *
    StrictHostKeyChecking no
'

mkdir -p ~/.ssh
printf "$DISABLE_HOSTKEY_VERIFICATION" | tee ~/.ssh/config

go build -o aws-api  main.go

if [[ $? -ne 0 ]]; then
    echo "BUILD FAILED"
    exit
fi
echo 'BUILD OK'

chmod 755 aws-api

./aws-api