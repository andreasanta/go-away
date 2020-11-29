# Introduction

"GO Away" is a simple service to check whether an IP address belogs to an abusive list fetched from [firehol/blocklist-ipsets](https://github.com/firehol/blocklist-ipsets).

# Quick Start

Pull the repository locally then run the following command to download the IP database:

`go run ./cmd/downloader/main.go`

Once the download is finished, you will find an SQLlite db in /data/ips.db.

Then you can run the HTTP server that responds to the malicious IP query:

`go run ./cmd/server/main.go`

> NOTE: give it a couple of minutes to load all the IP addresses in memory, they are about 2.6M at the time of writing. This should be the output once the service is ready to respond:

    2020/11/29 10:29:17 Loading database, will take roughly 1 minute
    2020/11/29 10:29:58 Loaded single IP DB rows 2590272
    2020/11/29 10:30:21 Loaded range records rows 27460 
    2020/11/29 10:30:21 Loaded single All DB rows 2617732 in 64052 ms
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
    - using env:   export GIN_MODE=release
    - using code:  gin.SetMode(gin.ReleaseMode)

    [GIN-debug] GET    /api/v1/ip/:ip            --> goaway/internal/http.GetIp (3 handlers)
    [GIN-debug] GET    /api/v1/reload            --> goaway/internal/http.GetReload (3 handlers)
    [GIN-debug] Listening and serving HTTP on :8080


# Querying the server

Once the service is running, you can query it by simply entering `http://localhost:8080/api/v1/ip/<your_ip_to_check>` in your browser. For instance:

    Request
    http://localhost:8080/api/v1/ip/1.10.16.0

    Response
    {
        "status": "ko",
        "description": "et_block.netset",
        "timestamp": 1606641062
    }

In this case the IP was found in the suspicious list and the name of the list has been returned, as well as the time of last update.

# Hot updating the list

The malicious IP list can be updated while the server is running by first executing the downloader again:

    go run ./cmd/downloader/main.go

Once finished, the `/reload` endpoint can be queried, this way:

    http://localhost:8080/api/v1/reload

Remember it will take rougly two minutes to reload the database from disk, after that the service will respond with a confirmation of the number of loaded records and the total time of loading, for instance:

    {
        "records": 2617732,
        "status": "ok",
        "time": 73880
    }

> NOTE: the reload process is totally thread safe since and it does not interrupt the service execution per se.

# Quick Doker containerization

Unfortunately I did not have much time to write a proper Dockerfile, but you can find a basic one in "Dockerfile".

A very basic instance of the service can be run with:

    docker build -t goaway:latest .
    docker run goaway:latest

# TODO

- Improve docker containerization
- Implement IPTrie data structure