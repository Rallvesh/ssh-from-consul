.PHONY: build_amd
build_amd:
		GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/sfc cmd/ssh-from-consul/main.go

.PHONY: build_arm
build_arm:
		go build -o ./bin/darwin_arm64/sfc cmd/ssh-from-consul/main.go

.PHONY: run
run: build_arm
		./bin/darwin_arm64/sfc ls
