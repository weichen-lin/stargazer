#!/bin/bash

export NEO4J_URL=neo4j://localhost:7687
export NEO4J_PASSWORD=password
export AUTHENTICATION_TOKEN=localtesttotken

export TRANSFORMER_PORT=5000
export IS_PRODUCTION=false

export TRANSFORMER_URL=http://localhost:5000

export APP_PASSWORD=testapppassword
export KAFKA_URL=localhost:9092
export GIN_MODE=debug