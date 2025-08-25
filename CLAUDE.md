# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a Go-based review service built with the Kratos framework following clean architecture principles:
- **API Layer**: Protobuf definitions in `api/review/v1/` with HTTP/GRPC endpoints
- **Business Logic**: `internal/biz/` contains use cases and business rules
- **Data Access**: `internal/data/` handles database operations with GORM
- **Service Layer**: `internal/service/` implements gRPC service handlers
- **Configuration**: YAML configs in `configs/` with protobuf-based configuration

## Key Technologies
- **Framework**: Kratos v2 with dependency injection via Wire
- **Database**: MySQL/SQLite with GORM ORM and code generation
- **ID Generation**: Snowflake algorithm for distributed ID generation
- **Service Discovery**: Consul registry integration
- **Protocol**: gRPC and HTTP/JSON with protobuf validation

## Development Commands

### Build and Run
```bash
# Build the application
make build

# Run with configs
./bin/review-service -conf ./configs

# Generate all code (protobuf, wire, etc.)
make all
```

### Code Generation
```bash
# Generate protobuf code
make api

# Generate wire dependency injection
cd cmd/review-service && wire

# Generate GORM models
go generate ./...
```

### Protobuf Management
```bash
# Update proto dependencies
make init

# Generate validation code
make validate

# Generate error definitions
make errors
```

## Service Structure
- **Main Entry**: `cmd/review-service/main.go` with Wire-based DI
- **Business Logic**: Review creation, auditing, replies in `internal/biz/review.go`
- **Data Models**: GORM-generated models in `internal/data/model/`
- **API Endpoints**: Defined in `api/review/v1/review.proto`

## Database Schema
Check `review.sql` for database schema. Uses GORM with generated query code in `internal/data/query/`.

## Configuration
- Server config: `configs/config.yaml` 
- Registry config: `configs/registry.yaml`
- Uses protobuf-based configuration in `internal/conf/`

## Key Features Implemented
- Review creation with order validation
- Review auditing system
- Store reply functionality
- Snowflake ID generation
- Consul service registration
- Input validation via protobuf rules