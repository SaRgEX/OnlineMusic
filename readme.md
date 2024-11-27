## Run
```sh
cd ./cmd
go run main.go
```
___
## ENV
The .env file has not been hidden for simplicity

Default values
```ini
# database
POSTGRES_USER="postgres"
POSTGRES_HOST="localhost"
POSTGRES_PORT="5432"
POSTGRES_DB="postgres"
POSTGRES_SSLMODE="disable"
POSTGRES_PASSWORD="postgres"

HTTP_SERVER_ADDRESS=":8080"

# Logger
LOG_LEVEL="debug"
LOG_FILE_PATH="/log/log.log"

API_URL="http://localhost:8000
```
___
## Endpoints
- **GET** `/songs?[name]&[performer]&[startDate]&[endDate]&[cursor]&[pageSize]`
- **GET** `/songs/id/lyrics?[group]&[song]&cursor=1&pageSize=1`
- **POST** `/songs` value see below
- **DELETE** `/songs/{id}`
- **PUT** `/songs/{id}`
___
## Add song's input
Input type JSON
```json
{
    "performer_id": 1,
    "song": "example"
}
```
___
## Technologies
**go packages**
- gin
- log/slog
- pgx

**utils**
- migrate
- swagger