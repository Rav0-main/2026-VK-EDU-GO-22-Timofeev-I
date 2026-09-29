GO := $(shell command -v go 2> /dev/null)
RM := $(shell command -v rm 2> /dev/null)

.PHONY: release clear

release:
	$(GO) build -o calc.out .

clear:
	$(RM) -f *.out
