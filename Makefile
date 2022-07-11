test-golox:
	cd test/ && \
	interpreter=golox go test -v -count=1 ./...

test-clox:
	cd test/ && \
	interpreter=clox go test -v -count=1 ./...

test-clean:
	cd test/ && \
	go clean -testcache