dev: 
	nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run cmd/todos/main.go
