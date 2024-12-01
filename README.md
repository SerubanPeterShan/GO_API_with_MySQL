# GO API with MySQL - Time Service

A containerized Go API service that provides current time and time log functionality with MySQL database integration.

## Features

- Get current time in Toronto timezone
- View historical time request logs
- Containerized application using Docker
- MySQL database for persistent storage
- Timezone handling (America/Toronto)

## Prerequisites

- Docker
  - MySQL :latest
  - Golang : 1.23
  - Alpine :3.20
- Docker Compose
- Git

## Architecture

![Application Architecture](Architecture%20API.svg)

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/SerubanPeterShan/GO_API_with_MySQL.git
cd GO_API_with_MySQL
```

2. Start the containers:

```bash
docker-compose up -d
```

3. The API will be available at `http://localhost:80`

## API Endpoints

### GET /current-time

Returns the current time in Toronto timezone.

```json
{
    "current_time": "2024-11-30 22:05:50 EST (Toronto Time)"
}
```

### GET /request-logs

Returns a list of all time requests.
<table>
<tr>
<th>If no Records were made</th>
<th>with Records</th>
</tr>
<tr>
<td>
  
```json

{
    "message": "No time requests recorded yet"
}

```
  
</td>
<td>

```json

[
    "2024-11-30 22:10:11 EST (Toronto Time)",
    "2024-11-30 22:10:09 EST (Toronto Time)"
]

```

</td>
</tr>
</table>

## Database Setup

MySQL database is automatically initialized with the following configuration:

- Port: 3306
- Database: timedb
- User: root
- Password: rootpassword

The necessary tables are created automatically when the container starts.

### Database Initialization (init.sql)

The `init.sql` file is crucial for proper timezone handling in the MySQL database. It:

- Sets the global and session timezone to EST (UTC-05:00)
- Creates the time_log table with DATETIME type
- Uses proper character encoding (utf8mb4)
- Verifies timezone settings during initialization

This setup ensures consistent time storage and retrieval without timezone conversion issues, which is essential for accurate time logging in the Toronto timezone.

## How to read logs on DOCKER

```bash
docker ps
```

```bash
docker logs <container id>
```

## Error Messages and Troubleshooting

### Error opening database:

- **Cause**: This error occurs if there is an issue with the database connection string or the database server is not reachable.
- **Troubleshooting**:
  1. Verify that the environment variables `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME` are correctly set.
  2. Ensure that the MySQL server is running and accessible from the application.
  3. Check the network connectivity between the application and the database server.

### Error connecting to database:

- **Cause**: This error occurs if the application cannot establish a connection to the database after opening it.
- **Troubleshooting**:
  1. Ensure that the database credentials are correct.
  2. Verify that the database server is running and accepting connections.
  3. Check for any firewall rules or network issues that might be blocking the connection.

### Only GET method is allowed

- **Cause**: This error occurs if a request is made to the `/current-time` or `/request-logs` endpoints using a method other than GET.
- **Troubleshooting**:
  1. Ensure that the client making the request uses the GET method.
  2. Check the request method in the client code or API testing tool.

### Time zone conversion error

- **Cause**: This error occurs if there is an issue loading the "America/Toronto" timezone.
- **Troubleshooting**:
  1. Verify that the timezone "America/Toronto" is valid and supported by the system.
  2. Check for any issues with the Go `time` package or the system's timezone data.

### Database insert error

- **Cause**: This error occurs if there is an issue inserting the current time into the `time_log` table.
- **Troubleshooting**:
  1. Ensure that the `time_log` table exists in the database.
  2. Verify that the database user has the necessary permissions to insert data into the table.
  3. Check for any constraints or triggers on the `time_log` table that might be causing the insert to fail.

### Database query error

- **Cause**: This error occurs if there is an issue querying the `time_log` table.
- **Troubleshooting**:
  1. Ensure that the `time_log` table exists in the database.
  2. Verify that the database user has the necessary permissions to query data from the table.
  3. Check the query syntax and ensure it is compatible with the MySQL version being used.

### Database scan error

- **Cause**: This error occurs if there is an issue scanning the results of the database query.
- **Troubleshooting**:
  1. Ensure that the query returns the expected columns and data types.
  2. Check for any issues with the database driver or the Go `database/sql` package.

### Time parsing error

- **Cause**: This error occurs if there is an issue parsing the timestamp from the database.
- **Troubleshooting**:
  1. Verify that the timestamps in the `time_log` table are in the expected format (`2006-01-02 15:04:05`).
  2. Check for any discrepancies in the timestamp format or data corruption in the table.