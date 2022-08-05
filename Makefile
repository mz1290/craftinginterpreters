all: build-golox build-clox
	cp clox/build/clox clox-1.0.0 && \
	cp golox/golox golox-1.0.0

shell-clox: build-clox
	cd clox && \
	./build/clox

shell-golox: build-golox
	cd golox && \
	./golox

build-clox:
	cd clox && \
	make

build-golox:
	cd golox && \
	make

test-golox: build-golox
	cd test/ && \
	interpreter=golox go test -v -count=1 ./...

test-clox: build-clox
	cd test/ && \
	interpreter=clox go test -v -count=1 ./...

clean-test:
	cd test/ && \
	go clean -testcache

clean-clox:
	cd clox/ && \
	make clean && \
	cd .. && rm -f clox-1.0.0

clean-golox:
	rm -f golox/golox golox-1.0.0

clean: clean-golox clean-clox