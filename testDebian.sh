go run assets/asset_generate.go
go install
$HOME/go/bin/OpenPipelineIO -http :8080 -reviewrender -scanplaterender
