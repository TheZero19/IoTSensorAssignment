# IoT Sensor Assignment

A Go application that ingests IoT sensor readings, caches state in Redis, and periodically synchronizes aggregated data to PostgreSQL. The service exposes HTTP endpoints to register sensors and submit temperature readings. Authentication is enforced using:

- API Key + Sensor credentials for sensor registration
- Pre-Shared Key (PSK) for submitting readings

The project ships with a Docker Compose setup that starts the app, Redis, and PostgreSQL together.

---

## Features

- Sensor registration secured by an API key
- Reading ingestion secured by PSK (bcrypt hashed)
- Redis as the system of engagement for fast writes and reads
- Periodic background sync from Redis to PostgreSQL (upsert on conflict)
- Aggregated response: per-sensor averages and overall average

---

## Architecture

- HTTP Server (net/http)
  - GET `/` health probe
  - POST `/registerSensor` (Authorization: `API-KEY <API_KEY> <SensorID> <PSK>`) → hashes PSK and caches registration in Redis
  - POST `/inputPayloadFromSensor` (Authorization: `PSK <SensorID> <PSK>`) → verifies PSK, updates averages in Redis, returns aggregate averages
- Redis
  - Hash per sensor: key = `<SensorID>`, fields: `PSKHash`, `AverageTemperature`, `NumberOfReceivedReadings`
  - Set `dirty_sensors` tracking sensor IDs needing persistence
- PostgreSQL
  - Table `sensors` with unique `sensor_id`, hashed PSK, average temperature, and count of readings
- Background sync
  - Interval set by `SYNC_FACTOR` seconds
  - Reads each sensor in `dirty_sensors` from Redis and upserts into PostgreSQL; removes from set on success

---


## Project Structure

- `main.go` – server setup, routes, background sync start
- `environmentInit.go` – environment variable loading
- `Constants/` – global config, env var names, Redis keys, DB handles
- `Database/`
  - `dbInit.go` – PostgreSQL and Redis clients; table creation
  - `Synchronization/` – background sync loop from Redis → Postgres
  - `Models/` – in-memory models
- `Controllers/`
  - `Register/` – sensor registration handler
  - `SensorReading/` – reading ingestion and response building
  - `Utils/` – request validation, serialization, helper utilities
- `Auth/` – middleware interface and concrete implementations
  - `Concrete/ApiKeyAuthMiddleware.go` – API-KEY auth for registration
  - `Concrete/BcryptAuthMiddleware.go` – PSK auth for readings
- `Hash/` – bcrypt helpers
- `Dockerfile`, `docker-compose.yml`

---

## Quickstart (Docker Compose)

Prerequisites:
- Docker + Docker Compose

Steps:
1) Build and start all services:

```
docker compose up -d --build
```

2) Verify the app is running:

```
curl -i http://localhost:8080/
```

Expected response body: `IOT Assignment Server Up and Running`

---

## Configuration

Environment variables (provided via `docker-compose.yml` for the app service):
- `POSTGRES_USER` (default: `iotuser`)
- `POSTGRES_PASSWORD` (default: `secret`)
- `POSTGRES_DB` (default: `iotdb`)
- `POSTGRES_HOST` (default: `postgres` – the compose service name)
- `POSTGRES_PORT` (default: `5432`)
- `REDIS_HOST` (default: `redis` – the compose service name)
- `REDIS_PORT` (default: `6379`)
- `SENSOR_REGISTRATION_API_KEY` (default: `CEROSKY`)
- `SYNC_FACTOR` (default: `10`, seconds between background sync runs)

Notes:
- The registration API key is currently hard-coded in code as `CEROSKY`. For production, externalize this as an environment variable and avoid hard-coding secrets.

---

## Data Model

PostgreSQL table (auto-created on start if not exists):

```
CREATE TABLE IF NOT EXISTS sensors (
  ID SERIAL PRIMARY KEY,
  sensor_id TEXT NOT NULL UNIQUE,
  psk_hash TEXT NOT NULL,
  average_temperature FLOAT NOT NULL DEFAULT 0.0,
  num_of_received_readings INTEGER NOT NULL DEFAULT 0
);
```

Redis keys:
- Set `dirty_sensors`: tracks sensor IDs flagged for syncing
- Hash per sensor at key `<SensorID>` with fields:
  - `PSKHash` (string)
  - `AverageTemperature` (float stored as string)
  - `NumberOfReceivedReadings` (int stored as string)

---

## API

All endpoints are on port 8080.

1) Health
- Method: GET
- Path: `/`
- Response: `IOT Assignment Server Up and Running`

2) Register Sensor
- Method: POST
- Path: `/registerSensor`
- Auth: HTTP header `Authorization: API-KEY <API_KEY> <SensorID> <PSK>`
  - Example: `Authorization: API-KEY CEROSKY sensor-1 supersecret`
- Body: empty
- Response: `Json Received` (also writes to Redis and marks sensor dirty for syncing)

3) Submit Sensor Reading
- Method: POST
- Path: `/inputPayloadFromSensor`
- Auth: HTTP header `Authorization: PSK <SensorID> <PSK>`
  - Example: `Authorization: PSK sensor-1 supersecret`
- Body (JSON):

```
{
  "sensor_id": "sensor-1",
  "temperature": "22.5"
}
```

- Response (JSON):

```
{
  "overall_average": 22.5,
  "sensor_averages": {
    "sensor-1": 22.5
  }
}
```

Behavior:
- The service reads current `AverageTemperature` and `NumberOfReceivedReadings` for the sensor from Redis, computes a new average, writes it back, flags the sensor as dirty, and returns the overall average across sensors plus per-sensor averages.

Important: Overall and per-sensor averages are computed by scanning Redis keys matching the pattern `sensor*`. If your sensor IDs do not start with `sensor`, they will not be included in the overall aggregation unless you adjust the key pattern in code.

---

## Example Requests

Register a sensor:

```
curl -X POST \
  http://localhost:8080/registerSensor \
  -H "Authorization: API-KEY CEROSKY sensor-1 supersecret"
```

Submit a reading:

```
curl -X POST \
  http://localhost:8080/inputPayloadFromSensor \
  -H "Authorization: PSK sensor-1 supersecret" \
  -H "Content-Type: application/json" \
  -d '{"sensor_id":"sensor-1","temperature":"22.5"}'
```

---

## Running Locally Without Docker

Prerequisites:
- Go toolchain (matching the Dockerfile base; if the `golang:1.25-alpine` image tag is unavailable, use a recent Go version such as 1.22)
- Local Redis and PostgreSQL instances

Steps:
- Export the required environment variables (see Configuration) to point at your local DBs
- Build and run:

```
go build -o iot-app .
./iot-app
```

Note: The module imports use a local module path (e.g., `dependencies/...`). Ensure `go.mod` exists with `module dependencies` or update imports to your chosen module path.

---

## Background Sync

- Controlled by `SYNC_FACTOR` (seconds)
- On each tick:
  - Load sensor IDs from Redis set `dirty_sensors`
  - For each ID, read the Redis hash fields and upsert into PostgreSQL
  - On success, remove the sensor ID from `dirty_sensors`

---

## Security Considerations

- No TLS termination included; run behind a reverse proxy/load balancer that handles TLS
- No rate limiting or audit logging; consider adding for production
