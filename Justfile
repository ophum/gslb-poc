help:
    just --list

run:
    go run main.go
    
list:
    #!/bin/bash
    echo -e "ADDR\t\tPORT\tHEALTH"
    curl http://localhost:5000/servers -s  | jq -r '(.[] | [.Addr, .Port, .HealthOK]) | @tsv'

set ADDR PORT:
    curl -X POST -H "Content-Type: application/json" http://localhost:5000/servers -d '{"addr": "{{ ADDR }}", "port": {{ PORT }}}'