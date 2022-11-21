TAG=`git describe --always --tags | cut -c 2-`

run:
	@ dotenv -e example.env -- go run cmd/api/*.go

mocking:
	@ mockery --all --with-expecter

test:
	@ dotenv -e test.env -- npx prisma db push
ifdef FILE
	@ dotenv -e test.env -- ginkgo $(FILE)
else
	@ dotenv -e test.env -- ginkgo -v ./...
endif

lint:
	@ golangci-lint run

build_image:
	@ docker build . -t expenses_api:latest
	@ docker build . -t "expenses_api:$(TAG)"

push_postgres:
ifdef ENV
	@ dotenv -e $(ENV).env -- npx prisma migrate dev --skip-generate --skip-seed
else
	@ dotenv -e example.env -- npx prisma migrate dev --skip-generate --skip-seed
endif

gen_swag:
	@ swag init --dir cmd/api --output cmd/api/docs --parseInternal --parseDependency --parseDepth 1

fmt:
	@ go fmt ./...
	@ swag fmt
