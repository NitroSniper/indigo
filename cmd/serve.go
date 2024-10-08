/*
Copyright © 2024 Nitro Sniper <nitro@ortin.dev>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"time"

	"strconv"

	"github.com/NitroSniper/indigo/server"
	"github.com/NitroSniper/indigo/server/flavors"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"os"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), func(cmd *cobra.Command, args []string) error {
		path := args[0]
		stat, err := os.Stat(path)
		if os.IsNotExist(err) {
			return err
		}
		if stat.IsDir() {
			return fmt.Errorf("serve %s: is a directory, currently not supported", path)
		} else {
			return nil
		}
	}),
	RunE: func(cmd *cobra.Command, args []string) error {
		duration, err := time.ParseDuration(interval)
		if err != nil {
			return err
		}
		server.NewMarkdownServer("./example.md", duration, flavor, ":"+strconv.Itoa(port)).HostServer()
		return nil
	},
}

// flags
var (
	port     int
	interval string

	flavorsIds = map[flavors.Enum][]string{
		flavors.GitHub: {"github"},
		flavors.Pico:   {"pico"},
	}
	flavor flavors.Enum
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&port, "port", "p", 8000, "Port number to host the server on")
	serveCmd.Flags().StringVarP(&interval, "interval", "i", "1s", "Poll the file for changes at specified interval (e.g., 1s, 500ms, 2s)")
	serveCmd.Flags().VarP(
		enumflag.New(&flavor, "string", flavorsIds, enumflag.EnumCaseInsensitive),
		"flavor", "f",
		"CSS theme of markdown file; can be 'github' or 'pico'",
	)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
