# NEW API SPEC

## License

1. GET: `/users/{userID}/licenses`
    * response:
    ```JSON
    {
        "license":[
            {
                "key":"abcd-123-abcdefg-23425234",
                "expiretime":"2018-05-24",
                "camera":"ghsfasdfsdfa12345"
            },
            {
                "key":"abcd-456-jyugfhf-23425234",
                "expiretime":"2018-05-24",
                "camera":"dafssdfsfsdfa1234"
            },
            {
                "key":"abcd-799-gtyjkdf-23425234",
                "expiretime":"2018-05-24",
                "camera":"sdfsefsgxcdgerwef"
            }
        ]
    }
    ```

1. POST: `/users/{userID}/licenses`
    * payload:
    ```JSON
    {
        "license":[
            {"key":"abcd-123-abcdefg-23425234"}
        ]
    }
    ```
    * error
        1. 400  incorrect format
        1. 401  invaild licenses

1. PUT: `/users/{userID}/licenses`
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
    * error
        1. 400  incorrect format
        1. 401  invaild licenses