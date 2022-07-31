run:
	@ dotenv -e example.env -- go run cmd/api/*.go

mocking:
	@ mockery --all --with-expecter

test:
	@ make push_mongo ENV=test
ifdef FILE
	@ dotenv -e test.env -- ginkgo $(FILE)
else
	@ dotenv -e test.env -- ginkgo -v ./...
endif

lint:
	@ golangci-lint run

build_image:
	@ docker build -t expenses_api:latest .

push_mongo:
ifdef ENV
	@ dotenv -e $(ENV).env -- npx prisma db push
else
	@ dotenv -e example.env -- npx prisma db push
endif
