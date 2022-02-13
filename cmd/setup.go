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
	"time"

	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up the necessary Jetstream stream and consumer",
	Long:  `This command sets up the Stream and Consumers needed by the application, uses direct Jetstream manager`,
	Run:   setupInfra,
}

var (
	streamName      string
	consumerName    string
	streamSubjects  string
	consumerSubject string
	streamReplicas  int
)

const SECONDS_IN_A_YEAR = 24 * 365 * time.Hour

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringVarP(&streamName, "streamName", "s", "USER_TXN", "specify the stream name to use")
	setupCmd.Flags().StringVarP(&streamSubjects, "streamSubjects", "b", "USER_TXN.>*", "specify the subject the stream is related to, example USER_TXN.>")
	setupCmd.Flags().StringVarP(&consumerName, "consumerName", "c", "GRP_MAKER", "The durable name of the consumer.")
	setupCmd.Flags().StringVarP(&consumerSubject, "consumerSubject", "j", "USER_TXN.maker", "The subject where the consumer is subscribed to.")
	setupCmd.Flags().IntVarP(&streamReplicas, "streamReplicas", "r", 1, "The number of replicas of the stream.")
}

func setupInfra(cmd *cobra.Command, args []string) {

	jsi, err := infra.NewJetStream(appConfig.I.SeedPath, appConfig.I.NatsUri)

	if err != nil {
		zap.S().Fatalf("setup failure %v", err)
	}

	i := infra.NewInfraSetup(jsi, streamName, streamSubjects, consumerName, consumerSubject, SECONDS_IN_A_YEAR, streamReplicas)
	if err := i.Setup(); err != nil {
		zap.S().Fatalf("%v", err)
	}
	jsi.Close()

}
