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
	"context"
	"os"
	"os/signal"

	"github.com/balchua/demo-jetstream/pkg/consumer"
	"github.com/balchua/demo-jetstream/pkg/dtrace"
	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consumes a message from the Consumer",
	Long:  `Consumes a message from a given stream/consumer`,
	Run:   consume,
}

var (
	subscriberSubject     string
	subscribeConsumerName string
)

func init() {
	rootCmd.AddCommand(consumeCmd)
	consumeCmd.Flags().StringVarP(&subscribeConsumerName, "consumerName", "n", "GRP_MAKER", "The durable name of the consumer.")
	consumeCmd.Flags().StringVarP(&subscriberSubject, "subscriberSubject", "s", "USER_TXN.maker", "The subject where the subscriber is subscribed to.")

}

func consume(cmd *cobra.Command, args []string) {

	d := dtrace.SetupTracer()
	defer d.Close()
	n, err := infra.NewNats(appConfig.S.SeedPath, appConfig.S.NatsUri)
	if err != nil {
		zap.S().Fatalf("%v", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	con := consumer.NewConsumer(n)

	worker := make(chan bool)
	go con.Listen(ctx, worker, subscriberSubject, subscribeConsumerName, appConfig.S.SleepTimeInMillis)
	ch := make(chan os.Signal)
	signal.Notify(ch, unix.SIGPWR)
	signal.Notify(ch, unix.SIGINT)
	signal.Notify(ch, unix.SIGQUIT)
	signal.Notify(ch, unix.SIGTERM)
	<-ch
	zap.S().Info("Waiting for worker to complete")
	cancelFunc()
	<-worker
	zap.S().Info("ending process")
}
