# Device Center Monitoring

A Go-based microservices system for monitoring and managing device sensor data. The system consists of three main services: Master Service, Client Service, and Device Simulator.

## Architecture

The system consists of three services:

- **Master Service**: Main API service that handles sensor data ingestion, device management, and processes sensor readings via Redis queue
- **Client Service**: Receives and processes sensor data from the master service
- **Device Simulator**: Simulates IoT devices sending sensor data to the master service

## Prerequisites

- Go 1.25.2 or higher
- Docker and Docker Compose (for running infrastructure services)
- Make (for using Makefile commands)

## Infrastructure Setup

The project requires PostgreSQL, MongoDB, and Redis. You can start them using Docker Compose:

```bash
docker compose up -d
```

This will start:
- **PostgreSQL** on port `55432`
  - Username: `postgres`
  - Password: `postgres`
  - Database: `postgres`
- **MongoDB** on port `27017`
  - Username: `admin`
  - Password: `password`
  - Database: `myapp`
- **Redis** on port `6378`

To stop the infrastructure:

```bash
docker compose down
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/zevinza/device-center-monitoring.git
cd device-center-monitoring
```

2. Install Go dependencies:
```bash
go mod download
```

3. (Optional) Create a `.env` file to override default configuration:
```bash
cp .env.example .env
```

Or manually create a `.env` file with your configuration. See the configuration section below for available options.

## Running the Services

### Using Makefile Commands

The project includes a Makefile with convenient commands to run each service:

#### Run Master Service
```bash
make run-master
```
Starts the master service on port `8000` (default). The service includes:
- REST API endpoints for device and sensor management
- Sensor data ingestion endpoint
- Redis queue consumer for processing sensor readings
- Swagger documentation available at `/swagger/index.html`

#### Run Client Service
```bash
make run-client
```
Starts the client service on port `8002` (default). This service receives sensor data from the master service.

#### Run Device Simulator
```bash
make run-device
```
Starts the device simulator that sends sensor data to the master service every 3 seconds.

#### Generate Swagger Documentation
```bash
make swag-master
```
Generates Swagger/OpenAPI documentation for the master service. The documentation will be available in the `docs/master` directory.

### Running Services Manually

You can also run services directly using Go:

```bash
# Master Service
go run ./app/master-service

# Client Service
go run ./app/client-service

# Device Simulator
go run ./app/device-simulator
```

## Makefile Commands Summary

| Command | Description |
|---------|-------------|
| `make run-master` | Run the master service |
| `make run-client` | Run the client service |
| `make run-device` | Run the device simulator |
| `make swag-master` | Generate Swagger documentation for master service |

## Configuration

Each service has its own configuration file in `app/<service-name>/config/environment.go`. You can override these settings using environment variables or a `.env` file.

### Master Service Configuration
- Port: `8000` (default)
- PostgreSQL: `localhost:55432`
  - Database: `postgres`
  - Username: `postgres`
  - Password: `postgres`
- MongoDB: `localhost:27017`
- Redis: `localhost:6378`
- API Secret Key: `Fr46VTqmt3j7AjT0hDa` (default)

### Client Service Configuration
- Port: `8002` (default)

### Device Simulator Configuration
- Server Host: `localhost`
- Server Port: `8000`
- Server Endpoint: `/api/v1/master`
- Sensor ID: Configured in environment config

## API Endpoints

### Master Service

The master service provides the following main endpoints:
- `GET /api/v1/master/devices` - List devices
- `POST /api/v1/master/devices` - Create device
- `GET /api/v1/master/sensors` - List sensors
- `POST /api/v1/master/sensors` - Create sensor
- `POST /api/v1/master/sensors/ingest` - Ingest sensor reading
- `GET /swagger/index.html` - Swagger documentation

### Client Service

- `GET /` - Health check
- `POST /receive` - Receive sensor data from master service

## Development

### Project Structure

```
device-center-monitoring/
├── app/
│   ├── master-service/      # Main API service
│   ├── client-service/      # Client service
│   └── device-simulator/    # Device simulator
├── services/                # Shared services (database, cache, queue)
├── middleware/              # HTTP middleware
├── utils/                   # Utility functions
├── entity/                  # Base entities
├── migrations/              # Database migrations
├── docs/                    # API documentation
├── docker-compose.yml       # Infrastructure setup
└── Makefile                 # Build and run commands
```

## Notes

- The master service uses PostgreSQL and MongoDB for data persistence, and Redis for queue management
- Sensor readings are processed asynchronously via Redis queue
- The device simulator sends sensor data every 3 seconds by default
- All services use Fiber web framework
- Environment variables can be set via `.env` file or system environment

## License

See LICENSE file for details.
