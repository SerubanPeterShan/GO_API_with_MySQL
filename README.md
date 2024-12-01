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

## Custom Start

To customize the start of the application, you can modify the `docker-compose.yml` file or the environment variables as needed.

1. Modify the `docker-compose.yml` file to change the database credentials or other configurations:

```yaml
version: '3.8'
services:
  database:
    image: mysql:latest
    container_name: mysql
    environment:
      # Database environment variables
      # Set the root password for MySQL
      MYSQL_ROOT_PASSWORD: <new-rootpass-word>
      # Set the name of the database to be created
      MYSQL_DATABASE: <new-timedb>
      # Set the MySQL user
      MYSQL_USER: <new-user>
      # Set the password for the MySQL user
      MYSQL_PASSWORD: <new-password>
      TZ: America/Toronto
    ports:
      - "3306:3306" #3306:<can change the this port number>
    volumes:
      - db_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "newuser", "-p$$MYSQL_PASSWORD"]
      interval: 8s
      timeout: 5s
      retries: 5

  app:
    build: .
    container_name: GO_TIME_API
    ports:
      - "80:80"
    depends_on:
      database:
        condition: service_healthy
    environment:
      DB_HOST: database
      DB_PORT: 3306
      # Set the MySQL user
      DB_USER: <newuser>
      # Set the password for the MySQL user
      DB_PASSWORD: <new-password>
      # Set the name of the database
      DB_NAME: <new-db-name>
      TZ: America/Toronto

volumes:
  db_data:

networks:
  default:
    driver: bridge
```

2. Start the containers with the modified configuration:

```bash
docker-compose up -d
```

3. The API will be available at `http://localhost:80` with the new configurations.

## Manual Configuration Without Docker Compose

### Manual Configuration Using Docker

If you prefer to manually configure and run the application using Docker without Docker Compose, follow these steps:

### MySQL Setup

1. **Pull the MySQL Docker Image**:

```bash
docker pull mysql:latest
```

2. **Run the MySQL Container**:

```bash
docker run --name mysql -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_DATABASE=timedb -e MYSQL_USER=newuser -e MYSQL_PASSWORD=new-password -e TZ=America/Toronto -p 3306:3306 -v db_data:/var/lib/mysql -d mysql:latest
```

3. **Initialize the Database**:

  - Copy the `init.sql` script to the MySQL container and run it to set up the necessary tables and timezone settings:

```bash
docker cp init.sql mysql:/init.sql
docker exec -it mysql mysql -u newuser -pnew-password timedb < /init.sql
```

### Go Application Setup

1. **Build the Go Application Docker Image**:

```bash
docker build -t go_time_api .
```

2. **Run the Go Application Container**:

```bash
docker run --name go_time_api -e DB_HOST=mysql -e DB_PORT=3306 -e DB_USER=newuser -e DB_PASSWORD=new-password -e DB_NAME=timedb -e TZ=America/Toronto --link mysql:mysql -p 80:80 -d go_time_api
```

3. **Access the API**:

  - The API will be available at `http://localhost:80`.

By following these steps, you can manually configure and run the Go API with MySQL using Docker without Docker Compose.

### Manual Configuration without docker

If you prefer to manually configure and run the application without using Docker Compose, follow these steps:

#### MySQL Setup

1. **Install MySQL**:

- Download and install MySQL from the official website: [MySQL Downloads](https://dev.mysql.com/downloads/).

2. **Configure MySQL**:

- Start the MySQL server.
- Create a new database and user:

  ```sql
  CREATE DATABASE timedb;
  CREATE USER 'newuser'@'localhost' IDENTIFIED BY 'new-password';
  GRANT ALL PRIVILEGES ON timedb.* TO 'newuser'@'localhost';
  FLUSH PRIVILEGES;
  ```

3. **Initialize the Database**:

- Run the `init.sql` script to set up the necessary tables and timezone settings:

  ```bash
  mysql -u newuser -p new-password timedb < init.sql
  ```

### Go Application Setup

1. **Install Go**:

- Download and install Go from the official website: [Go Downloads](https://golang.org/dl/).

2. **Clone the Repository**:

  ```bash
  git clone https://github.com/SerubanPeterShan/GO_API_with_MySQL.git
  cd GO_API_with_MySQL
  ```

3. **Configure Environment Variables**:

- Set the necessary environment variables for the application:

#### Bash

  ```bash
  export DB_HOST=localhost
  export DB_PORT=3306
  export DB_USER=newuser
  export DB_PASSWORD=new-password
  export DB_NAME=timedb
  export TZ=America/Toronto
  ```

#### PowerShell

  ```powershell
  $env:DB_HOST="localhost"
  $env:DB_PORT="3306"
  $env:DB_USER="newuser"
  $env:DB_PASSWORD="new-password"
  $env:DB_NAME="timedb"
  $env:TZ="America/Toronto"
  ```


4. **Build and Run the Application**:

  ```bash
  go build -o go_time_api .
  ./go_time_api
  ```

5. **Access the API**:

- The API will be available at `http://localhost:80`.

By following these steps, you can manually configure and run the Go API with MySQL without using Docker Compose.

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