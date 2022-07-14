files_tests = `find ./specs -name "*.py"`
files = `find ./src ./specs -name "*.py"`

run:
	@ dotenv -- go run cmd/api/*.go

mocking:
	@ mockery --all --with-expecter

test:
	@ dotenv -e test.env -- ginkgo ./...

lint:
	@ golangci-lint run

build_image:
	@ docker build -t expenses_api:latest .
