# GACHA MACHINE TODO LIST

1. machine
    1. run a server and leason a port. `"GET" "/"` will resplnse "hello world"
        * GIVE: 80 port
        * WHEN: start a http server, send `"GET" "/"` request
        * THEN: get response "hello world"
        * GIVE: 8080 port
        * WHEN: start a http server, send `"GET" ":8080/"` request
        * THEN: get response "hello world"
    1. server support add handlers by route
        * GIVE: 80 port, a handler for `"GET" "/hello"` request
        * WHEN: start a server, add hander to router, send `"GET" "/hello"` request
        * THEN: get response "hello world"
1. gacha
    1. A handle suport url `/recstorage/{cameraID}?p={filename}` with GET method. Resposne a media file
        * GIVE: cameraID = "test", filename = "test.txt", userID="user"
        * WHEN: send a get file request.
        * THEN: get response with a file in payload, file is the same on `S3:bucket/{userID}/{cameraID}/{filename}`
    1. A handle suport url `/recstorage/{cameraID}?p={filename}` with GET method and credential. Resposne a media file
        * GIVE: cameraID = "test", filename = "test.txt", userID="user", credential in header
        * WHEN: send a get file request.
        * THEN: get response with a file in payload, file is the same on `S3:bucket/{userID}/{cameraID}/{filename}`
    1. get error when no filename
        * GIVE: cameraID = "test", filename = ""
        * WHEN: send a request to a test http handler
        * THEN: get a 400 response
    1. A handle suport url `/recstorage/{cameraID}?p={filename}` with POST method. Put a media file to storage
        * GIVE: cameraID = "test" and filename = "test.txt" in URL, userID="user" in header
        * WHEN: send a put file request
        * THEN: get ok response. storage will receive a new file on `S3:bucket/{userID}/{cameraID}/{filename}`
    1. A handle suport url `/recstorage/{cameraID}/date?s={startDay}&e={endDay}` with GET method. to get the list of db files. assume prifx of filename is date
        * GIVE: cameraID = "test", start day = "2018-05-08", end day ="2018-05-08" in URL, userID="user" in header
        * WHEN: send query db list request
        * THEN: get ok response and db list from `S3:bucket/{userID}/{cameraID}/{dates}*` in json fromat
    1. A API will return a route with many API handler
    1. read config file to set regain and bucket
1. storage
    1. API to get a new storage
        * GIVE: none
        * GIVE: call new storage API
        * THEN: get a storage structure with download()/upload() func
    1. a s3 downloader API wrapper getObject() from AWS SDK
        * GIVE: a file info(Bucket="test", filename="test/test.txt") and a writer
        * WHEN: call download API
        * THEN: write got file data="test" to writer
    1. API to get a new s3 client
        * GIVE: Region, AccessKeyID, SecretKey, SessionToken
        * GIVE: call connect API
        * THEN: get a client structure with s3 API func
    1. a s3 uploader API wrapper putObject() from SDK API.
        * GIVE: a file info(Bucket="test", filename="test/test.txt") and a reader
        * WHEN: call download API
        * THEN: put file data="test" from reader to s3
