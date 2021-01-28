# checkbuild
Sample app working in cli and server mode

### URL to check versions
[http://zabrosov.com:8888/](http://zabrosov.com:8888/)

## Requirements:
```bash
# docker -v
Docker version 20.10.2, build 2291f61

#kubectl version
Client Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.3", GitCommit:"1e11e4a2108024935ecfcb2912226cedeafd99df", GitTreeState:"clean", BuildDate:"2020-10-14T12:50:19Z", GoVersion:"go1.15.2", Compiler:"gc", Platform:"darwin/amd64"}

# helm version
version.BuildInfo{Version:"v3.5.0", GitCommit:"32c22239423b3b4ba6706d450bd044baffdcf9e6", GitTreeState:"dirty", GoVersion:"go1.15.6"}

#git version
git version 2.29.2
```

## Tree
```bash
.
├── Dockerfile              # use bulder-image for go build
├── LICENSE
├── README.md
├── build.sh
├── checkbuild.yaml          # sample config
├── cmd
│   ├── root.go              # CLI mode
│   └── server.go            # server mode
├── docker.youtrack.service  # sample for sytemd
├── go.mod
├── go.sum
├── main.go
└── pkg
    ├── cmp
    │   └── cmp.go           # compare build version
    └── controller
        └── index.go         # index.html page
```

### build
* git clone https://github.com/petr4/checkbuild.git
* cd ./checkbuild
* chmod +x build.sh && ./build.sh
* we use `.dockerignore` and `.gitignore`
```bash
chmod +x build.sh && ./build.sh
[+] Building 1.6s (16/16) FINISHED
 => [internal] load build definition from Dockerfile
 # skip......
 => => naming to docker.example.com/checkbuild:v1000
```
#### also check .dockerignore
```bash
cat .dockerignore
# Ignore everything
**

# Allow app
!*.go
!go.mod
!go.sum
```

### checkbuild as cli
```bash
docker run -it --rm --name ch docker.example.com/checkbuild:v1000 -h
WARN[0000] Using config file:
'checkbuild' apps check status, compare build version in URLs, Both contain a build number(git hash).
Make sure they are the same. If they are the same, test passes. If they are not, test fails.
This will require doing some light parsing.

Usage:
  checkbuild [flags]
  checkbuild [command]

Available Commands:
  help        Help about any command
  server      Run checkbuild as server

Flags:
  -c, --config string    config file (default is $HOME/checkbuild.yaml)
  -d, --debug            debug mode
  -h, --help             help for checkbuild
      --logfile string   log file (default is Stdout
  -u, --urls strings     Add urls, separated by ','; urls >=2 (default [https://qa.adobeprojectm.com/version,https://spark.adobe.com/version])
      --viper            use Viper for configuration (default true)

Use "checkbuild [command] --help" for more information about a command.

### run
docker run -it --rm --name ch docker.example.com/checkbuild:v1000 -u "https://qa.adobeprojectm.com/version,https://qa.adobeprojectm.com/version,https://spark.adobe.com/version" -d
WARN[0000] Using config file:
INFO[2021-01-28T01:20:31Z] Urls: [https://qa.adobeprojectm.com/version https://qa.adobeprojectm.com/version https://spark.adobe.com/version]
INFO[2021-01-28T01:20:31Z] Debug: true
INFO[2021-01-28T01:20:31Z] Logfile: true  # not implemented
INFO[2021-01-28T01:20:31Z] Resp: {0.0.1 1f006d8f70263f58a30c prod 200 https://spark.adobe.com/version}
INFO[2021-01-28T01:20:31Z] Resp: {0.0.1 a595fa0375dc7c291823 qa 200 https://qa.adobeprojectm.com/version}
INFO[2021-01-28T01:20:31Z] Resp: {0.0.1 a595fa0375dc7c291823 qa 200 https://qa.adobeprojectm.com/version}
https://spark.adobe.com/version:1f006d8f70263f58a30c: False
https://qa.adobeprojectm.com/version:a595fa0375dc7c291823: False
https://qa.adobeprojectm.com/version:a595fa0375dc7c291823: False
------------
Test: failed
```

### checkbuild as server
* must have configfile: checkbuild.yaml
* mount config to `$HOME`: /app/checkbuild.yaml

```bash
docker run -it --rm --name ch -p 8080:8080 -v $(PWD)/checkbuild.yaml:/app/checkbuild.yaml docker.example.com/checkbuild:v1000 server

WARN[0000] Using config file:
WARN[2021-01-28T01:26:39Z] Using config file: /app/checkbuild.yaml
WARN[2021-01-28T01:26:39Z] Server Init...:true
[GIN] 2021/01/28 - 01:27:27 | 200 |     1.273908s |      172.17.0.1 | GET      "/"
[GIN] 2021/01/28 - 01:28:20 | 404 |         8.3µs |      172.17.0.1 | HEAD     "/"
[GIN] 2021/01/28 - 01:29:17 | 200 |    1.3412889s |      172.17.0.1 | GET      "/"


# check answer
curl -Lv localhost:8080
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Thu, 28 Jan 2021 01:29:17 GMT
< Content-Length: 1483
< Content-Type: text/html; charset=utf-8
```

### Helm
* TODO configmount 
* Pass parameter 'server' to start command
* WIP in progress in case if you need it -> let to know me
* details could be here [deploy](https://github.com/petr4/app-sample#deploy-app-sample-container-with-deploysh)