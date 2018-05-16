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
        * GIVE: cameraID = "test", filename = "test.txt", userID="user" in header
        * WHEN: send a request to a test http handler
        * THEN: get response with a file in payload, file is the same on `S3:bucket/{userID}/{cameraID}/{filename}`
    1. get error when no filename
        * GIVE: cameraID = "test", filename = ""
        * WHEN: send a request to a test http handler
        * THEN: get a 400 response
    1. A API will return a route with many API handler
    1. A handle suport url `/recstorage/{cameraID}?p={filename}` with POST method.
1. storage
    1. a writer interface wrapper getObject SDK API.