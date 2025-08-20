#!/bin/bash
export APP_ENV=staging
source .env.staging
go run ./cmd/api/