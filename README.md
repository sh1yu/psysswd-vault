# psysswd-vault

A password vault for your password security.

It like a simply `1password`. I would like to see if it can be more function and easy to use in the future.

# Usage

use `--help` can see the usage:
```bash
$ go build
$ ./psysswd-vault --help
A password vault for your password security.

Usage:
  ./psysswd-vault [flags]
  ./psysswd-vault [command]

Available Commands:
  add         add a new account for given username
  export      export account info for given username
  find        find given account info
  help        Help about any command
  import      import account info for given username
  list        list account info for given username
  login       login vault and get a command shell
  register    register a new master account for storage password
  serve       serve start a server for remote sync
  sync        sync account info from remote addr
  version     print version

Flags:
  -c, --conf string       config file
  -h, --help              help for ./psysswd-vault
  -p, --password string   give your master password
  -s, --serve string      start a server for sync with given port
  -u, --username string   give your username

Use "./psysswd-vault [command] --help" for more information about a command.
```

- list all accounts (and use `-P` can show password directly):
```
./psysswd-vault list
please input your master password: *********
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
|        账号         |       用户名         |         密码         |           额外信息             |      更新时间       |
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
| test-account1       | test-user-name1     | ****************     | test1                         | 2020-09-28 17:50:13 |
| test-account2       | root                | *************        | test2                         | 2020-09-28 17:52:37 |
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
```

- search accounts and show password:
```bash
./psysswd-valut find test -P
please input your master password: *********
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
|        账号         |       用户名         |         密码         |           额外信息             |      更新时间       |
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
| test-account1       | test-user-name1     | 1234567890123456     | test1                         | 2020-09-28 17:50:13 |
| test-account2       | root                | 1234567890123        | test2                         | 2020-09-28 17:52:37 |
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
```
