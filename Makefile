files_tests = `find ./specs -name "*.py"`
files = `find ./src ./specs -name "*.py"`

run:
	@ dotenv -- go run cmd/api/*.go

mocking:
	@ mockery --all --with-expecter

test:
	@ dotenv -e test.env -- ginkgo ./...

single_test:
	@ dotenv -e test.env -- ginkgo $(FILE)

lint:
	@ golangci-lint run

build_image:
	@ docker build -t expenses_api:latest .

push_mongo:
	@ dotenv -- npx prisma db push
