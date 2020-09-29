package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/psy-core/psysswd-vault/config"
	"github.com/psy-core/psysswd-vault/persist"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

var serveCmd = &cobra.Command{
	Use:   "serve <port>",
	Short: "serve start a server for remote sync",
	Long:  `serve start a server for remote sync`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		vaultConf, err := config.InitConf(cmd.Flags().GetString("conf"))
		checkError(err)

		http.HandleFunc("/sync", syncHandlerWrapper(vaultConf))

		fmt.Println("server start at ", args[0], "...")
		err = http.ListenAndServe(":"+args[0], nil)
		checkError(err)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func syncHandlerWrapper(conf *config.VaultConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var data map[string]string
		err = json.Unmarshal(content, &data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		exist, valid, err := persist.CheckUser(conf.PersistConf.DataFile, data["username"], data["password"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		if !exist {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("user not registered: " + data["username"]))
			return
		}
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Permission Denied."))
			return
		}

		records, err := persist.DumpRecord(conf.PersistConf.DataFile, data["username"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		dataJson, err := json.Marshal(records)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write(dataJson)
	}
}
