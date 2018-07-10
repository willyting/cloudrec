# GACHA MACHINE

## INTURDCTION

* a cloud recording server

## BUILD

1. install GOLANG > v1.6
1. create folder `C:\workspace\gopath` or `/home/user/gopath`
1. set env. `set GOPATH=C:\workspace\gopath` or `export GOPATH=/home/user/gopath`
1. clone this repository under `%GOPATH%\src\` or `$GOPATH/src/`
1. (!!!!!!) change the folder name `cloudrec` to `gachamachine`
1. change work path to GOPATH
1. run `go build -v -o recserver.exe gachamachine

## API

1. GET/POST `/recstorage/{cameraID}?p={filename}`
  * HEADERS:
    1. "X-identityID"
    1. "X-accessKeyID"
    1. "X-secretKey"
    1. "X-sessionToken"
  * FUNC: 
    * download/upload file to S3 with key {identityID}/{cameraID}/{filename}
    * download/upload file to local filesystem with filename ./storage/{identityID}.{cameraID}.{filename}
