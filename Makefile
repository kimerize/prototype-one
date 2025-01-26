.PHONY: example
demo:
	go run -gcflags=all="-N -l" ./cmd/kimerize example
