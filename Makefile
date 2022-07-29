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
ifdef ENV
	@ dotenv $($(ENV).env) -- npx prisma db push
else
	@ dotenv -- npx prisma db push
endif
