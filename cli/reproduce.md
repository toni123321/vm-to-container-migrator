# Reproduce

```sh
go mod init vm2cont
cobra-cli init --viper --author "Antonio Takev tonitakev.com" --license apache
go run main.go
go build
go install
```