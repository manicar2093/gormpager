test:
	@ ginkgo -v ./...

mocking:
	@ rm -r mocks
	@ mockery

version:
	@ cz version -p
