# Calendar Backend

A lightweight, high-performance RESTful API for managing calendar events and subjects, built with Go and MongoDB. The project follows clean architecture principles, separating concerns into routing, controllers, middleware, repositories, and database layers.

## Tech Stack & Architecture

* **Language:** Go (Golang)
* **Database:** MongoDB
* **Environment Management:** Godotenv (`.env`)
* **Architecture Pattern:** Clean Architecture / Domain-Driven Design (DDD) layout:
    * **Controllers:** Handle incoming HTTP requests, input validation, and HTTP responses.
    * **Repositories:** Encapsulate database queries and direct data access logic.
    * **Models:** Define data structures and MongoDB BSON mappings.
    * **Routes & Middleware:** Manage API endpoint registration and cross-cutting concerns like authentication.

---

## Project Structure

```text
.
├── .env                  # Local environment configuration (ignored by git)
├── .gitignore            # Git ignore rules
├── go.mod                # Go module dependencies
├── go.sum                # Dependency checksums
├── main.go               # Application entry point
└── internal              # Private application and business logic
    ├── controller        # HTTP request handlers
    │   ├── event_controller.go
    │   └── subject_controller.go
    ├── database          # Database connection initializers
    │   └── mongodb.go
    ├── middleware        # Custom HTTP middlewares
    │   └── auth.go
    ├── model             # Data models / Struct definitions
    │   └── model.go
    ├── repository        # Database abstraction layer
    │   ├── event_repository.go
    │   └── subject_repository.go
    └── routes            # API routing definitions
        └── routes.go
```
## Key Features

* **Event Management:** Full CRUD operations for handling calendar events.
* **Subject Organization:** Categorize and link events to specific subjects or modules.
* **Middleware Authentication:** Secured endpoints protected by a dedicated auth middleware layer.
* **Decoupled Storage:** Abstracted database layer enabling decoupled development and easy testing.

---

## Setup

### 1. Prerequisites
* Go (1.20 or higher recommended)
* MongoDB instance (Local or Atlas)

### 2. Environment Setup
Create a `.env` file in the root directory and configure your variables:
```env
PORT=8080
MONGO_URI=mongodb://localhost:27017
DB_NAME=calendar_db
AUTH_SECRET=your_super_secret_jwt_key
```
### 3. Installation & Running

Clone the repository, navigate into the project directory and start the server:
```bash
git clone [https://github.com/boysonline/interneto-technologijos.git](https://github.com/boysonline/interneto-technologijos.git)
cd interneto-technologijos
go mod download
go run main.go
```
# API Documentation

All endpoints are grouped under the `/api` base path and are protected by API key middleware verification.

## Global Headers

Every request made to this API must include the following header for authentication:

| Header | Type | Description | Required |
| :--- | :--- | :--- | :--- |
| `X-API-Key` | String | Client validation token managed by middleware | Yes |

---

## Subjects API

Base Path: `/api/subjects`

### 1. Fetch All Subjects
* **URL:** `/`
* **Method:** `GET`
* **Description:** Retrieves a collection of all stored subjects.

### 2. Fetch Single Subject
* **URL:** `/:id`
* **Method:** `GET`
* **Description:** Retrieves details for a specific subject by its unique identifier.

### 3. Create Subject
* **URL:** `/`
* **Method:** `POST`
* **Description:** Provision a new subject resource.

### 4. Update Subject
* **URL:** `/:id`
* **Method:** `PATCH`
* **Description:** Performs partial modifications on an existing subject.

### 5. Delete Subject
* **URL:** `/:id`
* **Method:** `DELETE`
* **Description:** Permanently removes a subject resource.

---

## Events API

Base Path: `/api/events`

### 1. Fetch All Events
* **URL:** `/`
* **Method:** `GET`
* **Description:** Retrieves a collection of all calendar event structures.

### 2. Delete Event
* **URL:** `/:id`
* **Method:** `DELETE`
* **Description:** Permanently removes an event resource.

---

## Lectures API

Base Path: `/api/lectures`

### 1. Fetch Single Lecture
* **URL:** `/:id`
* **Method:** `GET`
* **Description:** Retrieves specific details for a calendar lecture entry.

### 2. Create Lectures
* **URL:** `/`
* **Method:** `POST`
* **Description:** Provision a new lecture entry or entries.

### 3. Update Lecture
* **URL:** `/:id`
* **Method:** `PATCH`
* **Description:** Performs partial modifications on a lecture schedule or content record.

### 4. Delete Lectures by Group
* **URL:** `/group/:group_id`
* **Method:** `DELETE`
* **Description:** Batch removes multiple lectures associated with a specific group identifier.

---

## Assignments API

Base Path: `/api/assignments`

### 1. Fetch Single Assignment
* **URL:** `/:id`
* **Method:** `GET`
* **Description:** Retrieves details and deadlines for a specific assignment.

### 2. Create Assignment
* **URL:** `/`
* **Method:** `POST`
* **Description:** Provisions a new assignment tracking instance.

### 3. Update Assignment
* **URL:** `/:id`
* **Method:** `PATCH`
* **Description:** Update parameters, details, or completion state for an assignment.
