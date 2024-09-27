# SLIB
for non-comertial use only
this library is just a help-lib that makes my work faster and more comfortable

---

# Getting Started
```bash
    go get https://github.com/sabbatD/slib
```
### Math package includes few math funcs.
### Stringutils package includes en alphabetical letters and digits, include pointers.

## Handle package
### ChangeConfig func obviously changes the configuration.
It recieves the folowing values
```go
rps      uint16
duration uint16
detain   bool
```
`RPS` - Requests per second (uint16)
`duration` - Duration of requests repetition  (uint16)
`detain` - Reflets method of a requests gerenation  (bool)

You must use ChangeConfig to configure programm.
```go
ChangeConfig(100, 5, true)
```
### Get() Post()
There is 2 functions of single request:
```go
Get(url)
Post(url, body)
```
`url` - URL of the server
`body` - Body of the request

### Attack()
And finnaly the Attack() function.
It recieves the `Method`, `URL`, and optionally `body` parameters and returns string, which describes how many requests was made per duration time.
Here is an example:

```go
Attack("GET", "https://localhost")
Attack("POST", "https://localhost", body)
```

Return:
```
1200 requests per 3.0003329s
```
---

### SABBAT