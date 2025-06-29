echo "Create Product..."

curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "code": "P001",
    "name": "HP 15 Laptop",
    "description": "High-performance laptop",
    "price": 12000,
    "user_id": 1
  }'

