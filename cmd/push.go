package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sh1yu/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push account info to remote server",
	Long:  `push account info to remote server`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {

		vaultConf, username, _ := runPreCheck(cmd)

		remoteServerAddr, err := cmd.Flags().GetString("remote")
		checkError(err)

		if remoteServerAddr == "" {
			remoteServerAddr = vaultConf.RemoteConf.ServerAddr
		}

		if remoteServerAddr == "" {
			fmt.Println("please provide remote server addr with -r or 'remote.server_addr' in config file.")
			os.Exit(1)
		}

		remoteToken := getRemoteCredential(vaultConf.Credentials, username)
		if remoteToken == "" {
			fmt.Println("please config 'credentials' for user " + username + " in config file before push.")
			os.Exit(1)
		}

		err = runPush(vaultConf.PersistConf.DataFile, username, remoteToken, remoteServerAddr)
		checkError(err)
	},
}

func init() {
	pushCmd.Flags().StringP("remote", "r", "", "given remote server addr.")
	rootCmd.AddCommand(pushCmd)
}

func runPush(dataFile, username, remoteToken, remoteServerAddr string) error {

	fmt.Printf("Push for username %v to remote %v ...\n", username, remoteServerAddr)

	records, err := persist.DumpRecord(dataFile, username)
	if err != nil {
		return err
	}

	recordsDataJson, err := json.Marshal(records)
	if err != nil {
		return err
	}
	recordsBase64 := base64.StdEncoding.EncodeToString(recordsDataJson)

	data := map[string]string{"username": username, "token": remoteToken, "records": recordsBase64}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(remoteServerAddr+"/up", "application/json", bytes.NewReader(dataJson))
	if err != nil {
		return err
	}

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(respContent))
	}

	return nil
}
