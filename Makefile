dev: 
	nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go
