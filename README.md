# ATM Management System

## Description
This project is an ATM management system developed using the Go programming language. Database operations are performed using PostgreSQL and the GORM ORM (Object-Relational Mapping) library but  switch to etcd will be made. With Docker Compose, the PostgreSQL database and Go application are containerized for easy management and deployment.

## Features
- **User Authentication:** Secure user authentication using JWT (JSON Web Token).
- **Account Management:** Creating, updating, and deleting user accounts.
- **Balance Inquiry:** Checking the current balance of user accounts.
- **Deposit:** Depositing money into user accounts.
- **Withdrawal:** Withdrawing money from user accounts.
- **PIN Change:** Changing the PIN code of user accounts.
- **Logging:** Keeping a log of application activities.

## Technologies Used
- **Go:** Used for developing the application logic.
- **PostgreSQL:** Used as the database management system.
- **GORM:** Used to interact with PostgreSQL in the Go application.
- **Docker Compose:** Containerizes the application and database for easy management and deployment.

## Setup and Run

### Prerequisites
- Docker
- Docker Compose
- Dockerfile

### Running the Application

1. **Clone the repository:**
   ```sh
   git clone https://github.com/utkubayguven/newapiproject.git
   ```

2. **Create a `.env` file with the following content:**
   ```env
   DB_HOST=db
   DB_PORT=5432
   DB_USER=root
   DB_PASSWORD=root
   DB_NAME=test_db
   JWT_SECRET=your_jwt_secret
   ```

3. **Run the application using Docker Compose:**
   ```sh
   docker-compose up -d
   ```

This command will run the PostgreSQL database and Go application, making the application accessible at `http://localhost:8080`.

## API Endpoints

### User Routes
- **Register:** `POST /user/register`
- **Login:** `POST /user/login`

### Account Routes (Protected)
- **Balance Inquiry:** `GET /account/balance/:accountNumber`
- **Withdrawal:** `POST /account/withdrawal` (with JSON body parameter)
- **Deposit:** `POST /account/deposit` (with JSON body parameter)
- **PIN Change:** `POST /account/pin-change/:id`
- **Delete Account:** `DELETE /account/deleteacc/:accountNumber`

### User Routes (Protected)
- **Delete User:** `DELETE /user/:id`

### Swagger Documentation
Swagger documentation is available at `http://localhost:8080/swagger/index.html`.

## Configuration
The configuration is managed through a JSON file and environment variables. The API port and request limits can be configured in the `config.json` file.

**Full Changelog**: https://github.com/utkubayguven/newapiproject/commits/v1.0.0