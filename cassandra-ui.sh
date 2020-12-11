#!/bin/sh

IP=`ping -q -c 1 -t 1 "$CASSANDRA_HOST" | grep PING | sed -e "s/).*//" | sed -e "s/.*(//"`
cassandra-web --hosts "$IP" --port '9042' --username "$CASSANDRA_USER" --password "$CASSANDRA_PASSWORD"