files_tests = `find ./specs -name "*.py"`
files = `find ./src ./specs -name "*.py"`

run:
	@ uvicorn src.cmd.api.expenses:app

test:
	@poetry run mamba $(files_tests) --format documentation --enable-coverage

fmt: ## Format all project files
	@add-trailing-comma $(files)
	@pyformat -i $(files)
	@isort src specs
