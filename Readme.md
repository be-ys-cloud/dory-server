# DORY - Server

Expose a simple API to manipulate Active Directory or OpenLDAP server.

* Password reinitialization
* Password changer
* Account Unlocking

**You must have LDAPS (port 636) active and open to use this project.**

## Configuration file

Must be name `configuration.json`. Content :

```json
{
  "ldap_server": {
    "admin": {
      "username": "username-that-can-manipulate-users-on-ad",
      "password": "password"
    },
    "base_dn": "base_dn",
    "filter_on": "(&(objectClass=person)(samaccountname=%s))",
    "address": "ad_address",
    "kind": "ad|openldap",
    "port": 636,
    "skip_tls_verify": true,
    "email_field": "mail"
  },
  "server": {
    "port": 8000,
    "base_path": "/",
    "database_path": ""
  },
  "totp": {
    "secret": "your_custom_key_here_which_is_at_least_25_characters_long",
    "custom_service_name": "TOTP display name (leave blank for auto)"
  },
  "features": {
    "disable_unlock": false,
    "disable_password_update": false,
    "disable_password_reinitialization": false,
    "disable_totp": false
  },
  "mail_server": {
    "address": "server_addr",
    "port": 25,
    "sender_address": "dory_noreply@localhost.local",
    "password": "Password (if any) to authenticate",
    "subject": "DORY",
    "skip_tls_verify": true,
    "sender_name": "DORY"
  },
  "front_address": "https://dory.local/"
}
```

* `ldap_server` : Handles the configuration of your LDAP server, which base values (bind DN, password, address, etc)
  * `kind` must be `openldap` or `ad` (which stands for Active Directory)
* `server` : Web server configuration
  * `database_path` : Location of the database file (only needed with TOTP enabled). Defaults to `./database.sql`
* `totp` : Enables TOTP feature : users can create a TOTP that can be used in replacement of email verification pipeline. This might be useful, especially if your LDAP server manages your mail server.
  *  `secret` : Must be a secret string, known only by server, which is at least 25 characters long. **Losing or changing this key will make all TOTP unusable !**
  * `custom_service_name` : Change the default value (which is `DORY - your_ldap_address`) to a custom value. Only useful for display.
* `features` : Allow users to disable some features of the tool. By default, all features are enabled (except `unlock` feature on OpenLDAP).

**Important note:** When using TOTP, this server **requires** a SQLite backend to store user-specific secrets.

## Generate doc

```shell
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g ./internal/swagger_expose.go -o ./api
```

## Run

* `docker build -t="dory:latest" .`
* `touch /path/to/your/database.sql && chmod 777 /path/to/your/database.sql`
* `docker run -v /path/to/your/configuration.json:/app/configuration.json -v /path/to/your/database.sql:/app/database.sql -p 8000:8000 dory:latest`

