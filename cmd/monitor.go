/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/balchua/demo-jetstream/pkg/monitoring"
	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "monitor the message lag",
	Long:  `Monitors and prints the message lag per second`,
	Run: func(cmd *cobra.Command, args []string) {
		c := http.Client{Timeout: time.Duration(1) * time.Second}
		uri := fmt.Sprintf("%s://%s:%d", appConfig.M.Scheme, appConfig.M.Host, appConfig.M.Port)
		m := monitoring.NewMonitor(appConfig.M, &c, uri)
		m.StartMonitor()
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
