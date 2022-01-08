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

	"github.com/balchua/demo-jetstream/pkg/dtrace"
	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/balchua/demo-jetstream/pkg/publisher"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	pubStreamName  string
	messageSubject string
	maxCount       int
	pubMessage     string
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish user transaction to NATS jetstream",
	Long:  `Publish user transaction message to NATS jetstream`,
	Run:   publish,
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&pubStreamName, "streamName", "s", "USER_TXN", "specify the stream name to publish message to")
	publishCmd.Flags().StringVarP(&messageSubject, "messageSubject", "b", "USER_TXN.>*", "specify the subject of the message to publish")
	publishCmd.Flags().StringVarP(&pubMessage, "message", "m", "", "the message to send to the stream")
	publishCmd.Flags().IntVarP(&maxCount, "maxCount", "c", 10, "the message to send to the stream")

}

func publish(cmd *cobra.Command, args []string) {

	d := dtrace.SetupTracer()
	defer d.Close()

	n, err := infra.NewNats(appConfig.P.SeedPath, appConfig.P.NatsUri)
	if err != nil {
		zap.S().Fatalf("%v", err)
	}

	ctx := context.Background()
	var span trace.Span

	pub := publisher.NewTransactionPublisher(n)
	for i := 0; i < maxCount; i++ {
		ctx, span = otel.Tracer("publisher").Start(ctx, "LoopPublish")
		//pub.Publish("{\"TransactionID\":1,\"UserID\":1,\"Status\":\"OK\",\"Amount\": 456.89}", "USER_TXN.maker")
		pub.SendMessage(ctx, pubMessage, messageSubject)
		span.End()
	}

}
