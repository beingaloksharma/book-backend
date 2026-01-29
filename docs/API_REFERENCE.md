# API Reference Guide

This document provides a comprehensive guide to the Book Store REST API. The API is designed around RESTful principles and uses JSON for data exchange.

## Base URL
`http://localhost:8080`

## Authentication
The API uses JWT (JSON Web Tokens) for authentication.
1. Obtain a token via `/auth/login`.
2. Include the token in the `Authorization` header for protected endpoints:
   `Authorization: Bearer <your_token>`

---

## 🔒 Authentication

### Register a New User
Create a new account with a specific role.

- **Endpoint**: `POST /auth/signup`
- **Access**: Public
- **Request Body**:
  ```json
  {
    "name": "Alok Sharma",
    "email": "alok@example.com",
    "password": "securepassword",
    "role": "USER" 
  }
  ```
  *Note: `role` can be "USER" or "ADMIN".*
- **Response** (201 Created):
  ```json
  {
    "message": "User created successfully"
  }
  ```

### Login
Authenticate and receive a JWT token.

- **Endpoint**: `POST /auth/login`
- **Access**: Public
- **Request Body**:
  ```json
  {
    "email": "alok@example.com",
    "password": "securepassword"
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

---

## 📚 Books

### List All Books
Retrieve a list of all available books.

- **Endpoint**: `GET /api/books`
- **Access**: Authenticated (User/Admin)
- **Response** (200 OK):
  ```json
  [
    {
      "ID": 1,
      "title": "The Go Programming Language",
      "author": "Alan A. A. Donovan",
      "price": 35.99,
      "stock": 50
    }
  ]
  ```

### Get Book Details
Retrieve details for a specific book.

- **Endpoint**: `GET /api/books/{id}`
- **Access**: Authenticated
- **Response** (200 OK):
  ```json
  {
    "ID": 1,
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan",
    "description": "The authoritative resource for Go.",
    "price": 35.99,
    "stock": 50
  }
  ```

---

## 🛡️ Admin Operations

### Add a Book
Add a new book to the inventory.

- **Endpoint**: `POST /api/admin/books`
- **Access**: Admin Only
- **Request Body**:
  ```json
  {
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "description": "A Handbook of Agile Software Craftsmanship",
    "price": 29.99,
    "stock": 100
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "message": "Book created successfully"
  }
  ```

### Update a Book
Modify existing book details.

- **Endpoint**: `PUT /api/admin/books/{id}`
- **Access**: Admin Only
- **Request Body**:
  ```json
  {
    "title": "Clean Code",
    "author": "Robert C. Martin",
    "description": "Updated Description",
    "price": 32.99,
    "stock": 90
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "message": "Book updated successfully"
  }
  ```

### Delete a Book
Remove a book from the inventory.

- **Endpoint**: `DELETE /api/admin/books/{id}`
- **Access**: Admin Only
- **Response** (200 OK):
  ```json
  {
    "message": "Book deleted successfully"
  }
  ```

---

## 👤 User Profile & Address

### Get Profile
View current user's profile information.

- **Endpoint**: `GET /api/profile`
- **Access**: Authenticated
- **Response** (200 OK):
  ```json
  {
    "ID": 2,
    "name": "Alok Sharma",
    "email": "alok@example.com",
    "role": "USER"
  }
  ```

### Add Address
Save a new shipping address.

- **Endpoint**: `POST /api/addresses`
- **Access**: Authenticated
- **Request Body**:
  ```json
  {
    "street": "123 Tech Lane",
    "city": "Bengaluru",
    "state": "KA",
    "zip_code": "560001",
    "country": "India"
  }
  ```
- **Response** (201 Created):
  ```json
  {
    "message": "Address added successfully"
  }
  ```

### List Addresses
View all saved addresses.

- **Endpoint**: `GET /api/addresses`
- **Access**: Authenticated
- **Response** (200 OK):
  ```json
  [
    {
      "ID": 1,
      "street": "123 Tech Lane",
      "city": "Bengaluru",
      "...": "..."
    }
  ]
  ```

---

## 🛒 Cart & Orders

### Add to Cart
Add a book to the shopping cart.

- **Endpoint**: `POST /api/cart`
- **Access**: Authenticated
- **Request Body**:
  ```json
  {
    "book_id": 1,
    "quantity": 2
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "message": "Item added to cart"
  }
  ```

### View Cart
Review items in the current cart.

- **Endpoint**: `GET /api/cart`
- **Access**: Authenticated
- **Response** (200 OK):
  ```json
  {
    "ID": 5,
    "items": [
      {
        "book_id": 1,
        "quantity": 2,
        "book": { "title": "Clean Code", "price": 29.99 }
      }
    ]
  }
  ```

### Place Order
Checkout and place an order using a saved address.

- **Endpoint**: `POST /api/orders`
- **Access**: Authenticated
- **Request Body**:
  ```json
  {
    "address_id": 1
  }
  ```
- **Response** (200 OK):
  ```json
  {
    "message": "Order placed successfully"
  }
  ```

### List Orders
View order history.

- **Endpoint**: `GET /api/orders`
- **Access**: Authenticated
- **Response** (200 OK):
  ```json
  [
    {
      "ID": 101,
      "amount": 59.98,
      "status": "PENDING",
      "items": [...]
    }
  ]
  ```
