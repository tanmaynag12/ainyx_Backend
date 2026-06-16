# User API — Ainyx Backend Assessment

A REST API to manage users with their date of birth. Built with Go, GoFiber, PostgreSQL and SQLC.

## What this does

Five endpoints — create, get, update, delete, and list users. When fetching a user, age is calculated dynamically from their DOB (not stored in the DB). Also returns days until their next birthday.

The age calculation accounts for whether the birthday has actually passed this year or not. The naive `currentYear - birthYear` gives wrong results if the birthday hasn't come yet this year, so we handle that properly.

## Stack

- **Go + GoFiber** — web framework, handles HTTP routing
- **PostgreSQL** — database
- **SQLC** — generates type-safe Go code from raw SQL queries
- **Uber Zap** — structured JSON logging
- **go-playground/validator** — input validation via struct tags

## Project Structure

```
cmd/server/         → entry point, wires everything together
config/             → loads env variables
db/migrations/      → SQL to create tables
db/sqlc/            → SQLC generated database layer
internal/
  handler/          → HTTP layer, parses requests and sends responses
  service/          → business logic, age calculation lives here
  repository/       → only layer that talks to the database
  models/           → request and response structs
  middleware/       → request ID injection and duration logging
  logger/           → Zap logger setup
```

Each layer has one job and doesn't know about the others. If the database changes, only the repository layer needs updating.

## Setup

**Prerequisites:** Go 1.21+, PostgreSQL

**1. Clone the repo**

```bash
git clone https://github.com/tanmaynag12/ainyx_Backend
cd ainyx_Backend
```

**2. Create the database**

```bash
psql -U postgres -c "CREATE DATABASE ainyx_db;"
psql -U postgres -d ainyx_db -f db/migrations/001_create_users.sql
```

**3. Set up environment**

Copy the example env file and fill in your credentials:

```bash
cp .env.example .env
```

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=ainyx_db
APP_PORT=8080
```

**4. Install dependencies**

```bash
go mod tidy
```

**5. Run**

```bash
go run cmd/server/main.go
```

Server starts on `http://localhost:8080`

## API Endpoints

| Method | Endpoint          | Description                   |
| ------ | ----------------- | ----------------------------- |
| POST   | /api/v1/users     | Create a user                 |
| GET    | /api/v1/users/:id | Get user by ID (includes age) |
| PUT    | /api/v1/users/:id | Update a user                 |
| DELETE | /api/v1/users/:id | Delete a user                 |
| GET    | /api/v1/users     | List all users (includes age) |

### Create User

```
POST /api/v1/users
Content-Type: application/json

{
  "name": "Alice",
  "dob": "1990-05-10"
}
```

Response `201`:

```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get User by ID

```
GET /api/v1/users/1
```

Response `200`:

```json
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35,
  "next_birthday_in_days": 144
}
```

### Update User

```
PUT /api/v1/users/1
Content-Type: application/json

{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

Response `200`:

```json
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### Delete User

```
DELETE /api/v1/users/1
```

Response `204` — no body.

### List All Users

```
GET /api/v1/users
```

Response `200`:

```json
[
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10",
    "age": 35,
    "next_birthday_in_days": 144
  }
]
```

### Error responses

All errors follow this structure:

```json
{
  "code": "USER_NOT_FOUND",
  "message": "no user with id 42"
}
```

| Code                | When                           |
| ------------------- | ------------------------------ |
| `INVALID_BODY`      | Request body is not valid JSON |
| `VALIDATION_FAILED` | Missing or invalid fields      |
| `INVALID_ID`        | ID in URL is not a number      |
| `USER_NOT_FOUND`    | No user with that ID           |
| `CREATE_FAILED`     | DB error on create             |
| `UPDATE_FAILED`     | DB error on update             |
| `DELETE_FAILED`     | DB error on delete             |

## Running Tests

```bash
go test ./internal/service/...
```

Tests cover the age calculation — birthday already passed this year, birthday not yet this year, and birthday falling on today.

## A few design decisions worth noting

**Why `DATE` not `TIMESTAMP` for dob** — `TIMESTAMP` includes time and timezone. Timezone differences can shift the stored date by a day, which would break age calculation. `DATE` stores only the date, no drift possible.

**Why age is not stored in the database** — age changes every year on the user's birthday. Storing it means you'd have to update every user record every year. Calculating it fresh on fetch is simpler and always correct.

**Why `next_birthday_in_days`** — wasn't in the spec but made sense to add given we already have the DOB and are doing date math anyway. Costs nothing extra.

**Why SQLC over an ORM** — with SQLC you write raw SQL and it generates the Go code. You know exactly what query runs. With an ORM like GORM, it generates the SQL for you and it's not always obvious what's happening under the hood.

**Request tracing** — every request gets a unique ID in the `X-Request-ID` response header. The same ID appears in all log lines for that request, so you can trace exactly what happened for any given request.
