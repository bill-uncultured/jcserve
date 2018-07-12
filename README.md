# jcserve
jcserve will give you a SHA512 hash of a password, encoded in base64

## Running jcserve
In the httpserver directory, run
```
go build
./httpserver
```

Then, with your http client, POST to
```
http://localhost:8080/hash
```
with the data 
```
password=YOUR_PASSWORD
```
The server will log the number of active and total connections.


To shut down the server, user your http client to access
```
http://localhost:8080/shutdown
```
