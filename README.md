# ghsearch
> search repository or code in github

## Require
You need a github api token where you can get from [here](https://github.com/settings/tokens)

## Install
`go install -ldflags '-X main.GITHUB_TOKEN=${YOUR_API_TOKEN}' github.com/byebyebruce/ghsearch/cmd/ghsearch@latest`  
or  
`go install github.com/byebyebruce/ghsearch/cmd/ghsearch@latest`  

## Usage
- search repository: `ghsearch microservice grpc`
- search code: `ghsearch --lang=rust --code example grpc`
- help: `ghsearch -h`
- if you didn't build github api token into bin you should use: `GITHUB_TOKEN=xxx ghsearch microservice grpc`