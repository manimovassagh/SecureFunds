# Bank API

This project is a simple bank API built using Go, Echo framework, and PostgreSQL. The API supports user authentication, account management, and transaction operations, including deposits, withdrawals, transfers, and fetching account history with pagination.

## Features

- User Authentication (Signup, Login)
- Account Management (Deposit, Withdraw, Transfer)
- Transaction History with Pagination
- Secure Access to Accounts

## Tech Stack

- Go
- Echo framework
- PostgreSQL
- GORM
- JWT for authentication

## Setup

### Prerequisites

- Go (version 1.16+)
- Docker
- Docker Compose

### Running the Application

1. **Clone the repository:**

    ```sh
    git clone https://github.com/yourusername/bank-api.git
    cd bank-api
    ```

2. **Create a `.env` file with the following environment variables:**

    ```env
    DB_HOST=db
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=bank
    JWT_SECRET=your_jwt_secret
    ```

3. **Run Docker Compose to start PostgreSQL:**

    ```sh
    docker-compose up -d
    ```

4. **Run the application:**

    ```sh
    go run main.go
    ```

## API Endpoints

### Authentication

- **Signup:**
    - URL: `/signup`
    - Method: `POST`
    - Body:
      ```json
      {
          "username": "testuser",
          "password": "testpassword"
      }
      ```

- **Login:**
    - URL: `/login`
    - Method: `POST`
    - Body:
      ```json
      {
          "username": "testuser",
          "password": "testpassword"
      }
      ```

### Account Operations

- **Deposit:**
    - URL: `/deposit`
    - Method: `POST`
    - Headers: `Authorization: Bearer <jwt_token>`
    - Body:
      ```json
      {
          "account_id": 1,
          "amount": 100.0
      }
      ```

- **Withdraw:**
    - URL: `/withdraw`
    - Method: `POST`
    - Headers: `Authorization: Bearer <jwt_token>`
    - Body:
      ```json
      {
          "account_id": 1,
          "amount": 50.0
      }
      ```

- **Transfer:**
    - URL: `/transfer`
    - Method: `POST`
    - Headers: `Authorization: Bearer <jwt_token>`
    - Body:
      ```json
      {
          "from_account_id": 1,
          "to_account_id": 2,
          "amount": 25.0
      }
      ```

### Transaction History

- **Get Account History:**
    - URL: `/account/:account_id/history`
    - Method: `GET`
    - Headers: `Authorization: Bearer <jwt_token>`
    - Query Parameters:
        - `page`: Page number (default: 1)
        - `limit`: Number of transactions per page (default: 10)

- **Get Transfer History:**
    - URL: `/account/:account_id/transfer-history`
    - Method: `GET`
    - Headers: `Authorization: Bearer <jwt_token>`
    - Query Parameters:
        - `page`: Page number (default: 1)
        - `limit`: Number of transactions per page (default: 10)

## Postman Collection

You can import the Postman collection provided in the `postman_collection.json` file to test the API endpoints.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details..