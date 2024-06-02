#!/bin/bash

# Set the collection and environment file paths
COLLECTION="postman/collection.json"
ENVIRONMENT="postman/environment.json"

# Run Newman with the collection and environment
newman run $COLLECTION -e $ENVIRONMENT