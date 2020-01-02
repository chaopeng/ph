#!/bin/bash

go test ./... && go install . && echo "ph installed."
