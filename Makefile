# moonBASIC — common dev tasks (Unix / Git Bash / WSL).
# On Windows without Make, use:  powershell -File scripts/dev.ps1 <target>

.PHONY: build-compiler build-moonrun test check run-spin-cube help

help:
	@echo "Targets: build-compiler, build-moonrun, test, check, run-spin-cube"

build-compiler:
	go build -o moonbasic .

build-moonrun:
	go build -tags fullruntime -o moonrun ./cmd/moonrun

test:
	go test ./...

check:
	go run . --check examples/mario64/main_entities.mb

# Opens a window — requires CGO + full runtime (see docs/BUILDING.md).
run-spin-cube:
	CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
