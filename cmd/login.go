package cmd

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "strings"

    "github.com/psy-core/psysswd-vault/config"
    "github.com/psy-core/psysswd-vault/persist"

    "github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "login vault and get a command shell",
    Long:  `login with master password, and get command shell`,
    Args:  cobra.NoArgs,
    Run:   RunLogin,
}

func init() {
    rootCmd.AddCommand(loginCmd)
}

func RunLogin(cmd *cobra.Command, args []string) {
    vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
    checkError(err)
    username, password, err := readUsernameAndPassword(cmd, vaultConf)
    checkError(err)

    exist, valid, err := persist.CheckUser(vaultConf.PersistConf.DataFile, username, password)
    checkError(err)
    if !exist {
        fmt.Println("user not registered: ", username)
        os.Exit(1)
    }
    if !valid {
        fmt.Println("Permission Denied.")
        os.Exit(1)
    }

    runLogin(vaultConf.PersistConf.DataFile, username, password)
}

func runLogin(dataFile, username, password string) {
    //give a shell
    stdinReader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("[%s] > ", username)

        //fixme 应该有办法监控处理上下左右箭头的按键事件
        line, err := stdinReader.ReadString('\n')
        checkError(err)

        token := strings.Fields(line)
        if len(token) == 0 {
            fmt.Println("usage: add | find | list | exit ")
            continue
        }

        switch token[0] {
        case "add":
            if len(token) != 3 && len(token) != 4 {
                fmt.Println("usage: add <account-name> <account-user> [extra-message]")
                continue
            }
            runAdd(dataFile, username, password, token[1:])
        case "find":
            _, args, flags, err := parseArgs(token, map[string]int{"-P": 0, "--plain": 0})
            if err != nil {
                fmt.Println("error:", err, " usage: find <account-keyword> [-P]")
                continue
            }
            if len(args) != 1 {
                fmt.Println("usage: find <account-keyword> [-P]")
                continue
            }

            _, ok1 := flags["-P"]
            _, ok2 := flags["--plain"]
            runFind(ok1 || ok2, dataFile, username, password, args[0])
        case "list":
            _, args, flags, err := parseArgs(token, map[string]int{"-P": 0, "--plain": 0})
            if err != nil {
                fmt.Println("error:", err, " usage: list [-P]")
                continue
            }
            if len(args) != 0 {
                fmt.Println("usage: list [-P]")
                continue
            }
            _, ok1 := flags["-P"]
            _, ok2 := flags["--plain"]
            runFind(ok1 || ok2, dataFile, username, password, "")
        case "exit":
            os.Exit(0)
        default:
            fmt.Println("usage: add | find | list | exit ")
        }
    }
}

func parseArgs(token []string, flagDesc map[string]int) (string, []string, map[string][]string, error) {
    if len(token) == 0 {
        return "", nil, nil, errors.New("empty token")
    }
    cmd := token[0]
    args := make([]string, 0)
    flags := make(map[string][]string)

    for offset := 1; offset < len(token); {
        t := token[offset]
        if flagParamNum, ok := flagDesc[t]; ok {
            if _, ok2 := flags[t]; !ok2 {
                flags[t] = make([]string, 0)
            }
            for j := flagParamNum; j > 0; j-- {
                offset++
                if offset >= len(token) {
                    return "", nil, nil, errors.New("invalid parameters for flag:" + t)
                }
                flags[t] = append(flags[t], token[offset])
            }
        } else {
            args = append(args, t)
        }
        offset++
    }

    return cmd, args, flags, nil
}
