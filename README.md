# Must-gather REST wrapper

The wraper should provide HTTP API to allow trigger OpenShift must-gather (full or targeted) and provide output archive again via HTTP.

Under initial development, if this will work well, it could be moved under Forklift organization.

## Usage

This is an early development version, steps below might change.

### Get it and run

```
$ go get github.com/aufi/must-gather-rest-wrapper
$ cd ~/go/src/github.com/aufi/must-gather-rest-wrapper
$ go run pkg/must-gather-rest-wrapper.go # Note, run oc login first to allow wrapper use your KUBECONFIG or set environment variables to adjust config options
```

### API examples

Start must-gather execution

```
$ curl -X POST -H "Content-Type: application/json" -d '{"image": "quay.io/konveyor/forklift-must-gather", "timeout": "15m"}' http://localhost:8080/must-gather
```

Get must-gather execution (status field values: new, inprogress, completed, error)

```
$ curl  http://localhost:8080/must-gather/15
```

Download must-gather archive (available only if must-gather execution status == "completed")

```
$ curl -OJ http://localhost:8080/must-gather/15/data
```

List all must-gather executions

```
$ curl  http://localhost:8080/must-gather
```

Example of must-gather JSON object returned by API
```
{
  "id": 15,
  "created-at": "2021-06-30T15:19:17.514594773+02:00",
  "updated-at": "2021-06-30T15:23:05.415732774+02:00",
  "custom-name": "",
  "status": "completed",
  "image": "quay.io/konveyor/forklift-must-gather",
  "image-stream": "",
  "node-name": "",
  "command": "",
  "source-dir": "",
  "timeout": "30m",
  "server": "",
  "archive-size": 95334,
  "exec-output": "[must-gather      ] OUT Using must-gather plug-in image: quay.io/maufart/forklift-must-gather\n[must-gather      ]..."
}
```

### Available API parameters

Parameters below are passed from JSON API to must-gather execution.

API param name | Description | Example
--- | --- | ---
image | must-gather image name | ```quay.io/maufart/forklift-must-gather```
image-stream | must-gather image with tags | ```quay.io/maufart/forklift-must-gather:custom-tag```
node-name | node name to execute must-gather pod on | ```cluster-test4-fcjdb-master-1```
command | custom command to be executed | ```PLAN=plan1 /usr/bin/targeted```
timeout | timeout for must-gather execution | ```60m```
source-dir | must-gather pod directory to take data from | ```/var/local/something```
server | k8s API server which should be used | (taken from KUBECONFIG by default)
custom-name | custom ID to query the must-gather execution without remembering its ID | ```forklift-plan1```

All params are optional. Empty POST request will run must-gather with default options configured on server side (see below).

## Configuration

The wrapper default options can be adjusted by OS Environment variables.

Option | Default value | Description
--- | --- | ---
PORT | ```8080``` | Port where the wrapper listens on
DB_PATH | ```./gatherings.db``` | Local storage for must-gather executions records, can be ephemeral or just in memory ```file::memory:?cache=shared```
DEFAULT_IMAGE | ```quay.io/konveyor/forklift-must-gather``` | Image name to be used if it was not specified in API call
TIMEOUT | ```20m``` | Timeout for must-gather execution
ARCHIVE_FILENAME | ```must-gather.tar.gz``` | Archive filename to be searched in must-gather execution directory to be provided to user as the result archive
CLEANUP_MAX_AGE | ```-1``` | Maximum age of must-gather executions kept available in the wrapper, -1 disables the deletion, e.g. ```24h```

## Notes

Checkout [doc](doc/README.md) directory for more information.

### Initial development progress

- <del>prepare HTTP endpoints in gin-gonic</del>
- <del>prepare database/storage for gatherings metadata</del>
- <del>prepare image build scripts</del>
- <del>implement create / get / list gatherings (into db)</del>
- <del>implement raw single must-gather execution based on gathering db record</del>
- <del>implement all needed oc adm must-gather options support</del>
- <del>add ENV variables-driven default must-gather options values</del>
- <del>not needed now: handle async/not-blocking gathering with sane limits (e.g. max 10 simul.gatherings)</del>
- <del>prepare serving of gathered archive</del>
- <del>add periodical obsolete data cleanup</del>
- add ocp auth for gin-gonic if needed
- use passed admin token from operator to exec must-gather (from operator)
- <del>basic tests</del>
- <del>document API usage</del>
