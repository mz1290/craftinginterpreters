build-clox:
	cd clox && \
	make

build-golox:
	cd golox && \
	make

test-golox: build-golox
	cd test/ && \
	interpreter=golox go test -v -count=1 ./...

test-clox:
	cd test/ && \
	interpreter=clox go test -v -count=1 ./...

clean-test:
	cd test/ && \
	go clean -testcache

clean-clox:
	cd clox/ && \
	make clean

clean-golox:
	rm -f golox/golox

clean: clean-golox clean-clox