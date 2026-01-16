# Go DDD Example

This project is a practical example of a Go application structured using **Domain-Driven Design (DDD)** principles. It demonstrates a clean architecture layout, integrating modern Go libraries and best practices for building scalable Backend APIs.

## 🚀 Tech Stack

The project utilizes a robust set of tools and libraries:

-   **Language:** [Go](https://golang.org/) (1.24+)
-   **Web Framework:** [Gin](https://github.com/gin-gonic/gin)
-   **ORM:** [GORM](https://gorm.io/gorm) (Supports MySQL/PostgreSQL)
-   **Dependency Injection:** [Wire](https://github.com/google/wire)
-   **Authorization:** [Casbin](https://github.com/casbin/casbin)
-   **Caching:** [Redis](https://github.com/redis/go-redis)
-   **Validation:** [Go Playground Validator](https://github.com/go-playground/validator)
-   **Migration:** [Gormigrate](https://github.com/go-gormigrate/gormigrate)
-   **Logging:** [Logrus](https://github.com/sirupsen/logrus)
-   **Configuration:** [Viper](https://github.com/spf13/viper)
-   **CLI:** [Urfave CLI](https://github.com/urfave/cli)
-   **Tracing:** OpenTelemetry (as seen in imports)

## 📂 Project Structure

The project follows a layered architecture to separate concerns:

```
├── config/         # Application configuration
├── db/             # Database connection setup (MySQL, Redis)
├── migration/      # Database migration scripts
├── model/          # Domain entities and data models
├── repository/     # Data access layer (DB operations)
├── usecase/        # Application business logic (Service layer)
├── router/         # HTTP handlers and route definitions
├── singleton/      # Singleton instances (Session, etc.)
├── util/           # Utility functions (Masking, Logging, etc.)
└── main.go         # Application entry point
```

## 🛠️ Prerequisites

Ensure you have the following installed:

-   **Go** 1.24 or higher
-   **MySQL** (or compatible database)
-   **Redis**

## ⚙️ Configuration

1.  Clone the repository:
    ```bash
    git clone https://github.com/a5932016/go-ddd-example.git
    cd go-ddd-example
    ```

2.  Initialize the configuration file:
    Copy the example environment file to a new `.env` file.
    ```bash
    cp .env.example .env
    ```

3.  Update the `.env` file with your local credentials:
    ```ini
    # Server Configuration
    CORE_FE_API_MODE=debug
    CORE_FE_API_PORT=8010

    # Database Settings
    MYSQL_HOST=127.0.0.1
    MYSQL_PORT=3306
    MYSQL_USER=your_user
    MYSQL_PASSWORD=your_password
    MYSQL_DB_NAME=your_db_name

    # Redis Settings
    REDIS_HOST=127.0.0.1
    REDIS_PORT=6379
    REDIS_PASSWORD=your_redis_password
    ```

## 🏃 Getting Started

1.  **Install Dependencies:**
    ```bash
    go mod download
    ```

2.  **Run the Application:**
    ```bash
    go run main.go
    ```

    The server will start on the port specified in `.env` (default is `8010`).

3.  **Generate Dependency Injection (Optional):**
    If you modify the dependency graph, regenerate `wire_gen.go`:
    ```bash
    wire
    ```

## 📝 API Overview

The application exposes RESTful APIs using Gin. Key features typically include:

-   **Authentication:** Login, Logout (Session-based).
-   **User Management:** CRUD operations for users.
-   **Authorization:** Role-based access control using Casbin.

_(Check `router/` for detailed route definitions)_

## 📄 License

This project is licensed under the [MIT License](LICENSE).
