# Expenses API

An api to register expenses, incomings and take controll on how you move your money

## Tools

It is need to install the next tools to start the project

- [dotenv-cli](https://www.npmjs.com/package/dotenv-cli)
- [sqlacodegen](https://pypi.org/project/sqlacodegen/)

## Development

You can use this script to create entities.

```bash
sqlacodegen postgresql://development:development@localhost:3456/expenses_app-dev > gen_entities
```

Result is saved at `gen_entities`. You must copy this to `src/entities/__init__.py` ensuring you add needed imports
