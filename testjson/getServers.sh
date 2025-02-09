#!/bin/bash


# Replace 'your_username' and 'your_password' with actual values
curl -c cookies.txt -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"username": "u1", "password": "1"}'

# Replace 'actual_user_id' with the actual user ID you want to access
curl -b cookies.txt -X GET http://localhost:8080/api/user/1/servers \
     -H "Content-Type: application/json"

