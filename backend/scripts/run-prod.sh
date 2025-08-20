#!/bin/bash
export APP_ENV=production
source .env.production
go run ./cmd/api/