# GitPages

Serve static files from your git repo.
Update with webhooks.
Customize with `.gitpages.toml` files just in your repo.

## Install

```bash
go get github.com/oott123/gitpages
```

## Config

### Server config

Create toml config file named `config.toml` located in working directory or
`config` directory inside working directory contains following content:

```toml
Endpoint = ":2289" # listen endpoint
StorageDir = "data" # git repos and worktrees are saved here

[[Servers]]
Host = "yelp.github.io" # match `Host` header
Remote = "https://github.com/Yelp/yelp.github.io.git" # git remote
Branch = "master" # which branch to serve
Dir = "/" # which dir inside git repo to serve
WebHookSecret = "gitpages" # update webook secret

[[Servers]]
Host = "*" # use `*` for wildcard matching
Remote = "https://github.com/oott123/gitpages-example.git"
WebHookSecret = "gitpages"
Branch = "master"
Dir = "/"
```

Checkout [godoc](https://godoc.org/github.com/oott123/gitpages/pkg/config#Config)
for more details.

### Access rules

Create toml config file named `.gitpagesfile` contains following contents:

```toml
AllowCORS = true
NotFoundErrorPage = "/bar/error.html"

[[Rules]]
Match = "^/foo/.*"
AllowCORS = false
AllowListDirectory = true
```

Check out [godoc](https://godoc.org/github.com/oott123/gitpages/pkg/fileserver#AccessConfig)
to see the full list of the options.

## License

AGPLv3
