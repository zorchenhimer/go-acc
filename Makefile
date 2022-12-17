
.PHONY: clean

acc: main.go
	go build -o $@

clean:
	-rm -f acc acc.exe
