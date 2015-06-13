all:
	goimports -w .
	go build -o ldjtab cmd/ldjtab/main.go

clean:
	rm -f ldjtab
