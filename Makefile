
.PHONY: clean all

all: acc_freebsd acc acc.exe

acc_freebsd: *.go
	GOOS=freebsd GOARCH=amd64 go build -o $@

acc: *.go
	GOOS=linux   GOARCH=amd64 go build -o $@

acc.exe: *.go
	GOOS=windows GOARCH=amd64 go build -o $@

clean:
	-rm -f acc acc.exe
