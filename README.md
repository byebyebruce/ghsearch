# ghsearch
> search repository or code in github

## Require
You need a github [token](https://github.com/settings/tokens)

## Install
`go install github.com/byebyebruce/ghsearch/cmd/ghsearch@latest`  
or build your token into binary  
`go install -ldflags '-X main.GITHUB_TOKEN=${YOUR_TOKEN}' github.com/byebyebruce/ghsearch/cmd/ghsearch@latest`

## Usage
- search repository:`ghsearch microservice grpc`
- search code `ghsearch --lang=rust --code example grpc`
- help `ghsearch -h`
