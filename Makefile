run:
	go run main.go

test:
	go test -v -cover ./...

benchmark:
	go test -bench=. -benchtime=500ms ./...
