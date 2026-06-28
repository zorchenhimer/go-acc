
.PHONY: clean all

all: acc_freebsd acc

acc_freebsd: *.go
	GOOS=freebsd GOARCH=amd64 go build -o $@

acc: *.go
	go build -o $@

clean:
	-rm -f acc acc.exe
