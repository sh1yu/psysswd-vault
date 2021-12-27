# psysswd-vault

A password vault for your password security.

It like a simply `1password`. I would like to see if it can be more function and easy to use in the future.

# Usage

use `--help` can see the usage:
```bash
$ go build -o pvlt
$ sudo mv pvlt /usr/local/bin
$ pvlt --help
A password vault for your password security.

Usage:
  pvlt [flags]
  pvlt [command]

Available Commands:
  add         add a new account for given username
  config      configure configurations
  export      export account info for given username
  find        find given account info
  help        Help about any command
  import      import account info for given username
  list        list account info for given username
  pull        pull account info from remote server
  push        push account info to remote server
  register    register a new master account for storage password
  remove      remove a account for given username
  serve       serve start a remote server
  version     print version

Flags:
  -c, --conf string       config file
  -h, --help              help for pvlt
  -p, --password string   give your master password
  -u, --username string   give your username

Use "pvlt [command] --help" for more information about a command.
```

- list all accounts (and use `-P` can show password directly):
```bash
$ pvlt list
please input your master password: ***********
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
|        账号         |       用户名         |         密码         |           额外信息             |      更新时间       |
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
| test-account1       | test-user-name1     | ****************     | test1                         | 2020-09-28 17:50:13 |
| test-account2       | root                | *************        | test2                         | 2020-09-28 17:52:37 |
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
```

- search accounts and show password:
```bash
$ pvlt find test -P
please input your master password: ***********
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
|        账号         |       用户名         |         密码         |           额外信息             |      更新时间       |
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
| test-account1       | test-user-name1     | 1234567890123456     | test1                         | 2020-09-28 17:50:13 |
| test-account2       | root                | 1234567890123        | test2                         | 2020-09-28 17:52:37 |
+---------------------+---------------------+----------------------+-------------------------------+---------------------+
```

you must do `pvlt register` first before you use. see `pvlt help register` for more help.

and you can config `user.defaultUserName` in `config.yaml` to avoid you use `-u` parameter for convenience.

# Remote sync

use `pvlt serve` then you get a remote server:
```bash
$ pvlt serve 8888                                                                                  2021-12-27 11:55:27
server start at  8888 ...
```

then you can use `push` or `pull` for your local data synced with this server:

```bash
$ pvlt push -r http://127.0.0.1:8888                                                               2021-12-27 11:57:15
please input your master password: ***********
Push for username user1 to remote http://127.0.0.1:8888 ...

$ pvlt pull -r http://127.0.0.1:8888                                                               2021-12-27 11:57:50
please input your master password: ***********
Pulling for username user1 from remote http://127.0.0.1:8888 ...
import complete. total: 24, insert: 0, update: 0, ignore:24, err: 0
```

you must configure `credentials` configurations before you use `push` or `pull` command.

you can config `remote.server_addr` in `config.yaml` to avoid use `-r` parameter for convenience.

And you can use `pvlt help <command>` for more help.

enjoy!
