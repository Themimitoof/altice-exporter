# Altice Exporter

This project is a Prometheus Exporter for consumer hardware created by Altice Labs (previously
Portugal Telecom Inovacão) and deployed by ISPs of Altice's group like SFR and Altice Portugal (MEO).

The main purpose of this Prometheus Exporter is to collect metrics from the xPON transceiver of the
Fibergateway routers or the ONT provided by the ISP and setup a bunch of alerts on Alertmanager to
be informed when an important change in the DBm measured by the xPON transceiver happen or simply
when the xPON transceiver have loss the signal from the PON.

Of course, in order to be alerted when the xPON transceiver have loss the signal from the PON, a
backup path (xDSL, Cellular, LoRaWAN, etc.) should be available.

If you are a SFR/Red-By-SFR customer and still have a NB6, you probably don't need this Prometheus
Exporter. The NB6 exposes a "REST API" that connects to the ONT and return all the information in an
XML format. A small Python script could be easily created to push the metrics to Prometheus via a
Pushgateway. The documentation for REST API of the NB6 is available in the attachment of
[this message on the forum *lafibre.info*](https://lafibre.info/sfr-les-news/spec-api-rest-box-de-sfr/msg772775/#msg772775).

This Prometheus Exporter becomes more useful if you bypassed your NB6 in order to use your own
router. In that case you only need to add an IP to your WAN interface in order to communicate with
the ONT.

If you are a customer with a FiberGateway (all MEO customers with a FTTH and now SFR customers with
the new SFR BOX 8), the bypass is now more complicated since the xPON transceiver is now integrated
inside the router. This Prometheus Exporter can connect to your FiberGateway to retrieve the
information of the xPON transceiver.


## Device compatibility

| Reference  | Type         | Compatibility | Default IP Address              |
|------------|:-------------|:--------------|:--------------------------------|
| GR241AG    | FiberGateway |      ✅       | MEO: 192.168.1.254            |
| GS0100GH   | ONT          |      ✅       | SFR/Red-By-SFR: 192.168.4.254 |

*Note for SFR customers:* If you have a mini-ONT (from Alcatel Lucent), you probably have 99% of
chance that this Prometheus Exporter works. The *GS0100GH* ONT from Altice Labs mimics the Alcatel
Lucent CLI (probably the ONT is simply a disguised Alcatel Lucent hardware 🤔).

*Note 2 for SFR customers:* In order to avoid any problems with Altice, the username and the
password of the ONT used by SFR is available somewhere in
[this message of this *lafibre.info* topic](https://lafibre.info/remplacer-sfr/captures-de-linterface-du-mini-ont-sfr/msg256505/#msg256505).

*Note for all customers with a FiberGateway:* In theory all FiberGateways should work with this
Prometheus Exporter but have not been tested. The FiberGateways that should not work with this
Prometheus Exporter are customers with a coaxial termination (FTTLa).


## Metrics

| Name            | Type  | Description                                                                                    |
|-----------------|-------|------------------------------------------------------------------------------------------------|
| `link_status`   | Gauge | Tells the status of the PON. Value is `1` if a signal came from the PON, `0` if not.           |
| `rxrssi_metric` | Gauge | Gives the RX RSSI value in `dBm` from the PON receiver. In reception, the value is negative.   |
| `txrssi_metric` | Gauge | Gives the TX RSSI value in `dBm` from the PON receiver. In transmission, the value is positive.|


## Usage

```bash
$ ./altice-exporter --help
usage: altice-exporter --device.type=DEVICE.TYPE --device.hostname=DEVICE.HOSTNAME
--device.username=DEVICE.USERNAME --device.password=DEVICE.PASSWORD [<flags>]

Flags:
  --help                        Show context-sensitive help (also try --help-long and --help-man).
  --device.type=DEVICE.TYPE     Device type to collect
  --device.hostname=DEVICE.HOSTNAME
                                IP Address:port of the device to collect
  --device.username=DEVICE.USERNAME
                                Username of the device to collect
  --device.password=DEVICE.PASSWORD
                                Password of the device to collect
  --web.listen-address=":9876"  Update the bind address:port for the exporter
  --web.route-path="/metrics"   Update the route where the metrics will be exposed
```

The values available for `--device.type=` are the same as the references in the
`Device compatibility` table. Until the compatibility list is updated, for all FiberGateways, use
the value `GR241AG`. For all separate ONTs (and the SFR mini-ONT), use the value `GS0100GH`.

If you have a FiberGateway, here the command to use:

```bash
./altice-exporter \
    --device.type GR241AG \
    --device.hostname 192.168.1.254:22 \
    --device.username meo \
    --device.password meo
```

For a ONT, here the command to use:

```bash
./altice-exporter \
    --device.type GS0100GH \
    --device.hostname 192.168.4.254:22 \
    --device.username xxx \
    --device.password xxx
```

Here's an example of SystemD unit file:

```ini
[Unit]
Description=Altice Labs Prometheus Exporter
After=network.target

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
User=prometheus
ExecStart=/usr/local/bin/altice-exporter --device.type GS0100GH --device.hostname 192.168.4.254:22 --device.username xxx --device.password xxx
```

## License

This project is licensed under [*MIT License*](LICENSE.md).
