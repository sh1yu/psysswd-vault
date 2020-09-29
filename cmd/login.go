package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login vault and get a command shell",
	Long:  `login with master password, and get command shell`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		vaultConf, username, password := runPreCheck(cmd)

		servePort, err := cmd.Flags().GetString("serve")
		checkError(err)

		runLogin(vaultConf.PersistConf.DataFile, username, password, servePort)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(dataFile, username, password, servePort string) {

	ch := make(chan struct{})

	if servePort != "" {
		go func() {
			http.HandleFunc("/sync", syncHandlerWrapper(dataFile))

			fmt.Println("server start at ", servePort, "...")
			close(ch)
			err := http.ListenAndServe(":"+servePort, nil)
			checkError(err)
		}()
	}

	<-ch
	//give a shell
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("[%s] > ", username)

		//fixme 应该有办法监控处理上下左右箭头的按键事件
		line, err := stdinReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		checkError(err)

		token := strings.Fields(line)
		if len(token) == 0 {
			fmt.Println("usage: add | find | list | sync | exit ")
			continue
		}

		switch token[0] {
		case "add":
			_, args, flags, err := parseArgs(token, map[string]int{"-g": 0, "--genpass": 0})
			if err != nil {
				fmt.Println("error:", err, "usage: add <account-name> <account-user> [extra-message] [-g]")
				continue
			}
			if len(args) != 2 && len(args) != 3 {
				fmt.Println("usage: add <account-name> <account-user> [extra-message] [-g]")
				continue
			}

			_, ok1 := flags["-g"]
			_, ok2 := flags["--genpass"]
			err = runAdd(dataFile, username, password, ok1 || ok2, args)
			if err != nil {
				fmt.Println("error:", err, "usage: add <account-name> <account-user> [extra-message] [-g]")
				continue
			}
			fmt.Printf("add account %s success.\n", args[0])
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
			err = runFind(ok1 || ok2, dataFile, username, password, args[0])
			if err != nil {
				fmt.Println("error:", err, " usage: find <account-keyword> [-P]")
				continue
			}
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
			err = runFind(ok1 || ok2, dataFile, username, password, "")
			if err != nil {
				fmt.Println("error:", err, " usage: find <account-keyword> [-P]")
				continue
			}
		case "sync":
			if len(token) != 2 {
				fmt.Println("usage: sync <remote-addr>")
				continue
			}
			runSync(dataFile, username, password, token[1])
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("usage: add | find | list | sync | exit ")
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
