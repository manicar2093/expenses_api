# Expenses API

An api to register expenses, incomings and take controll on how you move your money

## Tools

It is need to install the next tools to start the project

- [dotenv-cli](https://www.npmjs.com/package/dotenv-cli)

## Development

To start the docker container you can run:

```bash
docker run -p 8000:8000 -e DATABASE_URL="postgresql://development:development@<ip:port>/expenses_app-dev?sslmode=disable" --name expenses-api-containter expenses_api:latest
```

For example:

```bash
docker run -p 8000:8000 -e DATABASE_URL="postgresql://development:development@192.168.100.48:3456/expenses_app-dev?sslmode=disable" --name expenses-api-containter expenses_api:latest
```
