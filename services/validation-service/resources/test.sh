#!/bin/bash

echo "posting invalid schema"
curl http://localhost:11000/schema/config-schema -X POST -d @invalid.json -v
echo ""
echo "posting valid schema should return 201"
curl http://localhost:11000/schema/config-schema -X POST -d @config-schema.json -v
echo ""
echo "get schema"
curl http://localhost:11000/schema/config-schema -X GET -v
echo ""
echo "posting invalid json doc"
curl http://localhost:11000/validate/config-schema -X POST -d @invalid.json -v
echo ""
echo "posting invalid config json -> validation fails"
curl http://localhost:11000/validate/config-schema -X POST -d @config-invalid.json -v
echo ""
echo "posting valid config json -> validation returns 200"
curl http://localhost:11000/validate/config-schema -X POST -d @config.json -v
