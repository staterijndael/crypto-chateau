# test timeout
TEST_TIMEOUT?=5m

# Запустить тесты
.PHONY: test
test:
	go test -parallel=10 -timeout=$(TEST_TIMEOUT) $(PWD)/...

# create coverage web page
.PHONY: cover
cover:
	go test -timeout=$(TEST_TIMEOUT) -v -coverprofile=coverage.out ./...  && go tool cover -html=coverage.out
