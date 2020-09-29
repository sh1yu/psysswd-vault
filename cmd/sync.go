package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"time"
)

var syncCmd = &cobra.Command{
	Use:   "sync <remote-addr>",
	Short: "sync account info from remote addr",
	Long:  `sync account info from remote addr`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultConf, username, password := runPreCheck(cmd)
		runSync(vaultConf.PersistConf.DataFile, username, password, args[0])
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSync(dataFile, username, password, remoteAddr string) {

	data := map[string]string{"username": username, "password": password}
	dataJson, err := json.Marshal(data)
	checkError(err)

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Post(remoteAddr+"/sync", "application/json", bytes.NewReader(dataJson))
	checkError(err)

	respContent, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("sync failed for status code:", resp.StatusCode, string(respContent))
	}

	var records []*persist.AccountRecord
	err = json.Unmarshal(respContent, &records)
	checkError(err)

	err = persist.ImportRecord(dataFile, records)
	checkError(err)
}
