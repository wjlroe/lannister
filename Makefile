include $(GOROOT)/src/Make.inc

TARG=lannister
GOFILES=lannister.go

include $(GOROOT)/src/Make.cmd

docs:
	@pandoc -s -w man -o lannister.1 README.md
	@godoc -html > docs/lannister.html