#!/bin/sh

set -ex

sleep $TRAFFIC_DRIVER_DELAY 
./hey -c $CONCURRENT_REQUESTS -q $REQUESTS_PER_SECOND -z $DURATION http://localhost:8000$SERVICE_ENDPOINT