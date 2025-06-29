echo "List Products..."

curl "http://localhost:8080/api/v1/products?page=1&page_size=10&status=active&name=laptop"
