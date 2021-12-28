# Playing with NATS Jetstream

This project demonstrate the use of NATS Jetstream.  NATS Jetstream is the persistent storage engine built into core NATS server.

For more information about NATS Jetstream visit this [link](https://docs.nats.io/nats-concepts/jetstream)

This project will demonstrate the publishing of message from a publisher and a consumer based on different subjects. The consumer will then save the records into a simple PostgreSQL database.

## Rules

Publisher will publish the following a `UserTransaction` event.

```json

{
    "TransactionID":1,
    "UserID":1,
    "Status": "OK",
    "Amount": 456.89
}
```

The Subscriber will ingest this event and verify if the transaction is within the user threshold.

# Getting started

Clone this project

## Pre-requisites

* Install NATS using helm with Jetstream enabled.  
 
[Documentation](https://docs.nats.io/running-a-nats-service/introduction/running/nats-kubernetes/helm-charts)
[Charts](https://github.com/nats-io/k8s/tree/main/helm/charts/nats)

  
```shell
kubectl apply -f hack/jetstream-pvc.yaml
helm upgrade --install --namespace nats bnats nats/nats -f hack/values.yaml
```

Check out the [`values.yaml`](hack/values.yaml) in the hack directory.

The above command starts a NATS server with Jetstream default values.

### Manually create the Stream and Consumers

Follow the Jetstream walkthrough - https://docs.nats.io/nats-concepts/jetstream/js_walkthrough

_Connect to the bnats box_

```shell
kubectl exec -n nats -it deployment/bnats-box -- /bin/sh -l
```

_Create the stream_

```shell
nats stream add USER_TXN --subjects "USER_TXN.*" --ack --max-msgs=-1 --max-bytes=-1 --max-age=1y --storage file --retention limits --max-msg-size=-1 --discard=old --max-msgs-per-subject=-1 --dupe-window=1d --allow-rollup --no-deny-delete --no-deny-purge --replicas=1
```

_Create `MAKER` consumer_

```shell
nats consumer add USER_TXN GRP_MAKER --filter USER_TXN.maker --ack explicit --pull --deliver all --max-deliver=-1 --replay=instant --max-pending=100 --no-headers-only 
```

_Create `VALIDATOR` consumer_

```shell
nats consumer add USER_TXN GRP_VALIDATOR --filter USER_TXN.validator --ack explicit --pull --deliver all --max-deliver=-1 --replay=instant --max-pending=100 --no-headers-only 
```

_Publish messages to the stream for maker_

```shell
nats pub USER_TXN.maker --count=50 --sleep 1s "publication #{{Count}} @ {{TimeStamp}}"
```

_Publish messages to the stream for validator_

Must use the subject to publish

```shell
nats pub USER_TXN.validator --count=50 --sleep 1s "publication #{{Count}} @ {{TimeStamp}}"
```

_View stream contents_

```shell
nats stream view txn_stream
```

_Get messages from stream for maker_

```shell
nats consumer next USER_TXN GRP_MAKER --count 50
```
## Building

```shell
go build.
```

## Auto create the STREAM and CONSUMERS

This project will automatically create the `STREAMS` and `CONSUMERS`

```shell
demo-jetstream setup
```

This will setup the STREAM `USER_TXN` and the CONSUMER `GRP_MAKER`


## Generate Account and authorization using nkeys

Download the nk tools

```shell
go install github.com/nats-io/nkeys/nk
```

Generate the keys

```shell
nk -gen user -pubout
SUAJSNDKKS4SLYV4BWYIF3RHP72PCF7LAXI6SIUIWLZW72DEBGFY6CCSAI
UB6WFVVI6BKTAHT5XGS55BONYOE3TDF47ZD7F75YVPABRXJ7XHWZKX2W
```

In the output above 

`Seed` (private key) - `SUAJSNDKKS4SLYV4BWYIF3RHP72PCF7LAXI6SIUIWLZW72DEBGFY6CCSAI`
`User` (public key) - `UB6WFVVI6BKTAHT5XGS55BONYOE3TDF47ZD7F75YVPABRXJ7XHWZKX2W`


These generated nkeys are stored in the [`seed.txt`](hack/seed.txt), this is used in the code.

```shell
nk -gen user -pubout
SUAMKIAMDUJITCXXXTL2XMHTVT3OBSA3KWLIZQ3NFBA4FMD3SQ75GJEF6Y
UD736QEIGXPHB5CLR4UAPCOEXET6WIKDYWELPIFHJJDJRNKH3SDHZTLT
```

Keep the `SUAMKIAMDUJITCXXXTL2XMHTVT3OBSA3KWLIZQ3NFBA4FMD3SQ75GJEF6Y` into [`sys-seed.txt`](hack/sys-seed.txt) and add the user key to the values.yaml

## Application configuration
In this demo, the sample application configuration is in [config.yaml](hack/config.yaml)

```yaml
infra:
  natsUri: "localhost:32422"
  seedPath: "hack/sys-seed.txt"
publish:
  natsUri: "localhost:32422"
  seedPath: "hack/sys-seed.txt"
monitor:
  scheme: "http"
  host: "localhost"
  port: 32822
  account: "demo"
  consumerName: "USER_TXN.maker"
  streamName: "USER_TXN"
  pollSeconds: 1
```

## Publish message

This will publish 10 messages to the stream on subject `USER_TXN.maker`
```shell
./demo-jetstream publish --config "hack/config.yaml" --streamName "USER_TXN" --messageSubject "USER_TXN.maker" --maxCount "10" --message "{\"TransactionID\":1,\"UserID\":1,\"Status\":\"OK\",\"Amount\": 456.89}"
```

## Monitor message lag

Using the configuration [here](#application-configuration).

The monitoring will check the message lag of the Consumer `USER_TXN.maker` in the Stream `USER_TXN` using the account `demo`

```shell
./demo-jetstream monitor --config "hack/config.yaml"
...
2021-12-28T17:54:30.352+0800	INFO	monitoring/monitor.go:47	total lag is 10100
2021-12-28T17:54:31.354+0800	INFO	monitoring/monitor.go:47	total lag is 10100
2021-12-28T17:54:32.357+0800	INFO	monitoring/monitor.go:47	total lag is 10100
2021-12-28T17:54:33.360+0800	INFO	monitoring/monitor.go:47	total lag is 10100
2021-12-28T17:54:34.365+0800	INFO	monitoring/monitor.go:47	total lag is 10100
2021-12-28T17:54:35.368+0800	INFO	monitoring/monitor.go:47	total lag is 10100
2021-12-28T17:54:36.371+0800	INFO	monitoring/monitor.go:47	total lag is 10100
...
```

## Consume message

```shell

./demo-jetstream consume --config "hack/config.yaml" --consumerName "GRP_MAKER" --subscriberSubject "USER_TXN.maker"
```
## Create Postgres DB

TODO


