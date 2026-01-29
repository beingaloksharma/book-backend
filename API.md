# Book Store API Documentation

## Authentication
### Signup
- **POST** `/auth/signup`
- Body:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "USER" 
}
```
(role can be `ADMIN` or `USER`)

### Login
- **POST** `/auth/login`
- Body:
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```
- Response: `{"token": "JWT_TOKEN"}`

## Books (Public)
### List Books
- **GET** `/api/books`
- Headers: `Authorization: Bearer <TOKEN>`

### Get Book
- **GET** `/api/books/:id`
- Headers: `Authorization: Bearer <TOKEN>`

## Admin Operations
### Create Book
- **POST** `/api/admin/books`
- Headers: `Authorization: Bearer <ADMIN_TOKEN>`
- Body:
```json
{
  "title": "Go Programming",
  "author": "Google",
  "description": "Learn Go",
  "price": 29.99,
  "stock": 100
}
```

### Update Book
- **PUT** `/api/admin/books/:id`
- Headers: `Authorization: Bearer <ADMIN_TOKEN>`

### Delete Book
- **DELETE** `/api/admin/books/:id`
- Headers: `Authorization: Bearer <ADMIN_TOKEN>`

## User Operations
### Profile
- **GET** `/api/profile`
- Headers: `Authorization: Bearer <TOKEN>`

### Address Management
- **POST** `/api/addresses`
- Body:
```json
{
  "street": "123 Main St",
  "city": "Metropolis",
  "state": "NY",
  "zip_code": "10001",
  "country": "USA"
}
```
- **GET** `/api/addresses`

### Cart
- **POST** `/api/cart`
- Body:
```json
{
  "book_id": 1,
  "quantity": 1
}
```
- **GET** `/api/cart`

### Order
- **POST** `/api/orders`
- Body:
```json
{
  "address_id": 1
}
```
- **GET** `/api/orders`
