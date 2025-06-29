echo "Create User..."

curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jon Doe",
    "email": "jon@doe.com"
  }'

