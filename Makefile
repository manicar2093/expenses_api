files_tests = `find ./specs -name "*.py"`
files = `find ./src ./specs -name "*.py"`

run:
	@ dotenv -- go run cmd/api/*.go

mocking:
	@ mockery --all --with-expecter

test:
	@ ginkgo ./...
