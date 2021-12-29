# Manually setup Streams and Consumers

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