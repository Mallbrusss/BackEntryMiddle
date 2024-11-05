#!/bin/sh

IP_POSTGRES=&(getent hosts postgres | awk '{print $1}')
IP_REDIS=$(getent hosts redis | awk '{ print $1 }')

echo "POSTGRES_HOST=${IP_POSTGRES}" > .env
echo "REDIS_HOST=${IP_REDIS}" >> .env

exec ./app-mid