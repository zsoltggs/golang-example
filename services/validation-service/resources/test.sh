#!/bin/bash

curl http://localhost:11000/schema/config-schema -X POST -d @config-schema.json -v
echo ""
echo ""
curl http://localhost:11000/schema/config-schema -X GET -v
echo ""
echo ""
curl http://localhost:11000/validate/config-schema -X POST -d @config.json -v
