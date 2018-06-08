# NEW API SPEC

## License

1. GET: `/users/{userID}/licenses`
    * decription: APP call this API to get all userID's the licenses. Each license have an expire time. A license may enable one camera to cloud recording.
    * response:
    ```JSON
    {
        "license":[
            {
                "key":"abcd-456-jyugfhf-23425234",
                "expireTime":"2018-05-24",
                "camera":"ghsfasdfsdfa12345"
            },
            {
                "key":"abcd-799-gtyjkdf-23425234",
                "expireTime":"2018-05-24",
                "camera":""
            }
        ]
    }
    ```

1. POST: `/users/{userID}/licenses`
    * decription: APP call this API to add a new licenses to userID. Server will check the license is available.
    * payload:
    ```JSON
    {
        "license":[
            {"key":"abcd-123-abcdefg-23425234"}
        ]
    }
    ```
    * response:
    ```JSON
    {
        "license":[
            {
                "key":"abcd-456-jyugfhf-23425234",
                "expireTime":"2018-05-24",
                "camera":"ghsfasdfsdfa12345"
            },
            {
                "key":"abcd-799-gtyjkdf-23425234",
                "expireTime":"2018-05-24",
                "camera":""
            },
            {
                "key":"abcd-123-abcdefg-23425234",
                "expireTime":"2018-05-24",
                "camera":""
            }
        ]
    }
    ```
    * error
        1. 400  incorrect format
        1. 401  invaild licenses

1. PUT: `/users/{userID}/licenses`
    * decription: APP call this API to enable a camera cloud recording with a license. Server will check the camera is belong the userID.
    * payload:
    ```JSON
    {
        "license":[
            {
                "key":"abcd-123-abcdefg-23425234",
                "camera":"sdfsefsgxcdgerwef"
            }
        ]
    }
    ```
    * response:
    ```JSON
    {
        "license":[
            {
                "key":"abcd-456-jyugfhf-23425234",
                "expireTime":"2018-05-24",
                "camera":"ghsfasdfsdfa12345"
            },
            {
                "key":"abcd-799-gtyjkdf-23425234",
                "expireTime":"2018-05-24",
                "camera":""
            },
            {
                "key":"abcd-123-abcdefg-23425234",
                "expireTime":"2018-05-24",
                "camera":"sdfsefsgxcdgerwef"
            }
        ]
    }
    ```
    * error
        1. 400  incorrect format
        1. 401  invaild licenses
1. GET: `/devices/{devicesID}/security`
    * decription: APP call this API to get all security info for a camera to cloud record. Server will create new info from AWS every time APP get those info. Server will check the camera is belone the user by auth token.
    * response:
    ```JSON
    {
        "cert":"",
        "privateKey":"",
        "expireTime":""
    }
    ```
    * error
        1. 400  not add to any account or no such camera
        1. 401  not enable recording
