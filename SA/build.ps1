$env:GOARCH="wasm";
$env:GOOS="js";
go build -o index-ga.wasm main.go
$env:GOARCH="";
$env:GOOS="";
