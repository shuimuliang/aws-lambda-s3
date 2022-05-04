#GOOS=linux go build -ldflags "-s -w" main.go
GOOS=linux go build main.go
zip function.zip main 
