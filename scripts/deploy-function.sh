#!/bin/bash
PROJECT_ID=a6c0754c-a2bd-41e4-b40a-85ebbc0efd1d
PROJECT_NAME=first-project
NAMESPACE_ID=fn-21d566fc-2901-4e9e-9434-d267db0d4d0d

# Connect to project's target namespace
doctl serverless connect $NAMESPACE_ID

doctl serverless deploy ../ --remote-build
