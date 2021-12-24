# Playing with MATS Jetstream

This project demonstrate the use of NATS Jetstream.  NATS Jetstream is the persistent storage engine built into core NATS server.

For more information about NATS Jetstream visit this [link](https://docs.nats.io/nats-concepts/jetstream)

This project will demonstrate the publishing of message from a publisher and a consumer based on different subjects. The consumer will then save the records into a simple PostgreSQL database.

## Rules

Publisher will publish the following a `UserTransaction` event.

```json

"userTransaction": {
    "userId": "123456789",
    "transactionMode": "buy",
    "amount": 123.456,
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
helm install --namespace nats bnats nats/nats --set nats.jetstream.enabled=true \
  --set nats.jetstream.fileStorage.enabled=true --set nats.jetstream.fileStorage.existingClaim=jetstream-pvc \
  --set nats.jetstream.fileStorage.claimStorageSize=1Gi
```
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

## Publish message

This will publish 10 messages to the stream on subject `USER_TXN.maker`
```shell
./demo-jetstream publish
```

## Create Postgres DB

TODO