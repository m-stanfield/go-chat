#!/bin/bash

# Get the number from the first argument, default to "1" if not provided
num=${1:-1}

# Construct the username and password using the number
username="u$num"
password="$num"

echo "Logging in with username: $username and password: $password"

# Login and store cookies
curl -c cookies.txt -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d "{\"username\": \"$username\", \"password\": \"$password\"}"

# Generate server name with timestamp
server_name="server-$(date '+%Y-%m-%d_%H-%M-%S')"

echo "Creating server with name: $server_name"

# Create a new server with the generated name
curl -b cookies.txt -X POST http://localhost:8080/api/newserver \
     -H "Content-Type: application/json" \
     -d "{\"servername\": \"$server_name\"}"

