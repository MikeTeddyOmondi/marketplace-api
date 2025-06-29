echo "Update Product..."

curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 11000,
    "description": Super High Performance PC"
  }'