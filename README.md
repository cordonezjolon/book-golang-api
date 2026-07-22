# Project README

## Configuration

This project requires a PostgreSQL database. Set the following environment variable in your `.env` file:

```env
DATABASE_URL=host=postgres port=5432 user= password= dbname=appdb sslmode=disable
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable` |
