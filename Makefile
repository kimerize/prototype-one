.PHONY: example
example:
	go run -gcflags=all="-N -l" ./cmd/kimerize example
