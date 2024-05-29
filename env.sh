#!/bin/bash

# Common
export NEO4J_URL=neo4j://localhost:7687
export NEO4J_PASSWORD=password
export AUTHENTICATION_TOKEN=7WYj5ciV7r0bsbYa3IrkC1zxiXM7VLOQhC61Iw/MRck=

# Transformer
export TRANSFORMER_PORT=5000
export IS_PRODUCTION=false

# Next-app and Producer
export TRANSFORMER_URL=http://localhost:5000

# Producer
export APP_PASSWORD=testapppassword
export KAFKA_URL=localhost:9092
export GIN_MODE=debug
export PRODUCER_PORT=:8080