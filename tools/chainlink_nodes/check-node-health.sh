#!/bin/bash

# Default chainlink ports
declare -a ports=("6711" "6722" "6733" "6744" "6755")

# Check each port
for port in "${ports[@]}"
do
    # Check once every 5 seconds, 10 times to see if the node is up
    for attempt in {1...10}
    do
        curl "http://localhost:$port" -q
        if [ $? -eq 0 ]
        then
            break
        fi
        sleep 5
    done
    if [ $? -ne 0 ]
    then
        echo "Unable to reach http://localhost:$port"
        exit 1
    fi
done
echo "All good"
exit 0