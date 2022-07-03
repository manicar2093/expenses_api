files_tests = `find ./specs -name "*.py"`
files = `find ./src ./specs -name "*.py"`

run:
	@ dotenv -- poetry run uvicorn src.cmd.api.expenses:app

mocking:
	@ mockery --all --with-expecter

test:
	@ ginkgo ./...
