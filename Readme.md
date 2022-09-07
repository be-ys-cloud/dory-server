# DORY - Server

Expose a simple API to manipulate AD.
* Password reinitialization
* Password changer
* Account Unlocking

**You must have LDAPS (port 636) active and open to use this project.**

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

## Generate doc

```shell
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./internal/swagger_expose.go -o ./api
```

## Run

* `docker build -t="dory:latest" .`
* `docker run -v /path/to/your/configuration.json:/app/configuration.json -p 8000:8000 dory:latest`

