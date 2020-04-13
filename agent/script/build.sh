#!/bin/bash

go build -o /g/agent/agent-$1.exe -ldflags "-X main.buildVersion=$1"