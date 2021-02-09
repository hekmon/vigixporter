# Vigixporter

"[Vigicrue](https://www.vigicrues.gouv.fr/)" Victoria Metrics exporter/pusher. Data is probbed from [hubeau](https://hubeau.eaufrance.fr/page/api-hydrometrie).

## How it works

To respect the timestamp of the original data source, vigixporter is not scrapable: it pushes the converted data with the original timestamp.

When a new station is added, the maximum backlog is fetched from the hubeau API (actually 30 days).

### State & Cache

The exporter maintains a state and caches that are dumped to disk in order to be restored at next start. The file (`vigixporter_state.json`) is written on the current working directory. For the deb package, the working directory is the home of the low privileged user created by the package and used by the systemd service; you can check this with `systemctl cat vigixporter.service`.

For each station the last data point timestamp is saved. On polling, vigixporter selects the oldest timestamp seen for all the tracked stations in order to request the minimum information from hubeau. Every point already known is skipped during processing to avoid duplicates. Once retreived, data is then converted and transmitted to the Victoria Metrics pusher component which maintains its own cache. If the victoria Metrics remote is offline, the converted datapoints are kept in cache until successfully pushed.

## Configuration

### environment variables

* `VIGIXPORTER_STATIONS` - The list of stations IDs to scrap from hubeau. An [online tool](https://hubeau.eaufrance.fr/sites/default/files/api/demo/hydro/index.htm) can be used to find IDs. Ex: `VIGIXPORTER_STATIONS=F700000103,F490000104,F664000404`.
* `VIGIXPORTER_VMURL` - The URL of the victoria metrics [import endpoint](https://github.com/VictoriaMetrics/VictoriaMetrics#how-to-import-data-in-json-line-format).
* `VIGIXPORTER_VMUSER` & `VIGIXPORTER_VMPASS` - Optional, HTTP basic auth for Victoria Metrics. Recommended for public URL (with `https` !).

### systemd (deb package)

The debian package has a systemd unit ready to be used: `vigixporter.service`. It reads its environment variables from the `/etc/default/vigixporter` file.

Once configured, start (and activated at startup) with: `systemctl enable --now vigixporter.service`. Then check the logs with `systemctl status vigixporter.service` or `journalctl -f -u vigixporter.service`.

## Pushed metrics

2 metrics are created with the same tag set. Exemple:

* `vigixporter_water_level{latitude="48.844690", longitude="2.365511", site_code="F7000001", station_code="F700000103"}` value is in millimeters
* `vigixporter_water_flow{latitude="48.844690", longitude="2.365511", site_code="F7000001", station_code="F700000103"}` value is in liters per second (divide by 1000 to obtain m3/s)
