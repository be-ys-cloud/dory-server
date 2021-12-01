# DORY - Server

Expose a simple API to manipulate AD.
* Password reinitialization
* Password changer
* Account Unlocking

**You must have LDAPS (port 636) active and open to use this project. Domain administrators will not be able to use DORY unless you change the AdminSDGroup permissions.**

## Endpoints :

* `/ask_reinitialization`
  * You must provide in body (JSON format) : `username`
* `/reinitialize`
  * You must provide in body (JSON format) : `username`, `token`, `new_password`
* `/ask_unlock`
  * You must provide in body (JSON format) : `username`
* `/unlock`
  * You must provide in body (JSON format) : `username`, `token`
* `/change_password`
  * You must provide in body (JSON format) : `username`, `old_password`, `new_password`

**This is NOT a REST API, and is not intended to be one ! HTTP Codes are still used for comprehensiveness, but the comparison ends here.**

## Configuration file

Must be name `configuration.json`. Content : 

```json
{
  "active_directory": {
    "admin": {
      "username": "username-that-can-manipulate-users-on-ad",
      "password": "password"
    },
    "base_dn": "base_dn",
    "filter_on": "(&(objectClass=person)(samaccountname=%s))",
    "address": "ad_address",
    "port": 636,
    "skip_tls_verify": true,
    "email_field": "mail"
  },
  "server": {
    "port": 8000,
    "base_path": "/"
  },
  "mail_server": {
    "address": "server_addr",
    "port": 25,
    "sender_address": "dory_noreply@localhost.local",
    "password": "Password (if any) to authenticate",
    "subject": "DORY",
    "sender_name": "DORY"
  },
  "front_address": "https://dory.local/"
}
```

## Configuration - AD rights

The user that will be used for Dory **MUST** have the following Active Directory rights :
* Change password
* Reset Password
* Write "lastPwdSet" property
* Write "lockoutTime" property
* Read users in your AD (obviously)




# Run

## Method 1 : From our DockerHub image :

`docker run -d -v /path/to/your/configuration.json:/app/configuration.json -p 8000:8000 beys/dory-server:latest`

## Method 2 : From sources

1. Git clone : `git clone https://github.com/be-ys/dory-server.git` and go to the `dory-server` folder
2. `export GOPATH=$(pwd)`
3. Get deps :
   1. `go get -t github.com/go-ldap/ldap`
   2. `go get -t github.com/sirupsen/logrus`
   3. `go get -t golang.org/x/text/encoding/unicode` 
   4. `go get -t github.com/thanhpk/randstr`
   5. `go get -t github.com/gorilla/mux`
4. `go run ./src`

# License

Distributed under MIT license. You will find a copy of the license [here](LICENSE).