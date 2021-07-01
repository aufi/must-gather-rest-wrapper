# Notes on must-gather REST wrapper development and setup

## Wrapper app

Should run in a container and provide HTTP API for Forklift UI to allow trigger must-gather (full or targeted) and provide output archive again via HTTP. The must-gather operation can take 10s seconds up to minutes, so it should be async. Consumer (Forklift UI) should periodically ask for result based on ID (received after must-gather trigger was accepted).

### Usage workflow

- UI creates request for triggering must-gather (POST to /must-gather endpoint with optional params for targeted gathering)
- Wrapper returns identifier for the must-gather result
- UI is polling for must-gather result (GET /must-gather/:identifier, e.g. every 5s, until status field is not completed or error)
- Once must-gather finished, status and link to must-gather archive is returned
- UI provides link to the archive to user and/or initiates download

Must-gather request status

- new
- inprogress
- completed
- error

## Open questions

Sqlite backend or just in memory? (after container restart there will be no archives anyway so no need for persistent db)
-> In memory (or use sqlite just with path flag inmemory) and communicate via chans

Authentication Wrapper <-> User from UI: not sure ?

Authentication Wrapper <-> OCP: need admin service account/token to execute must-gather (to be provided by operator)

## Command exec options

### oc adm must-gather

Supported options via API:

- <del>dest-dir - controller by the wrapper to be able find the archive</del>
- node-name
- image, image-stream
- source-dir - keep default for now MTV/MTC
- command
- timeout

```
$ oc adm must-gather --help
Launch a pod to gather debugging information

 This command will launch a pod in a temporary namespace on your cluster that gathers debugging information and then
downloads the gathered information.

 Experimental: This command is under active development and may change without notice.

Usage:
  oc adm must-gather [flags]

Examples:
  # gather information using the default plug-in image and command, writing into ./must-gather.local.<rand>
  oc adm must-gather
  
  # gather information with a specific local folder to copy to
  oc adm must-gather --dest-dir=/local/directory
  
  # gather audit information
  oc adm must-gather -- /usr/bin/gather_audit_logs
  
  # gather information using multiple plug-in images
  oc adm must-gather --image=quay.io/kubevirt/must-gather --image=quay.io/openshift/origin-must-gather
  
  # gather information using a specific image stream plug-in
  oc adm must-gather --image-stream=openshift/must-gather:latest
  
  # gather information using a specific image, command, and pod-dir
  oc adm must-gather --image=my/image:tag --source-dir=/pod/directory -- myspecial-command.sh

Options:
      --dest-dir='': Set a specific directory on the local machine to write gathered data to.
      --image=[]: Specify a must-gather plugin image to run. If not specified, OpenShift's default must-gather image
will be used.
      --image-stream=[]: Specify an image stream (namespace/name:tag) containing a must-gather plugin image to run.
      --node-name='': Set a specific node to use - by default a random master will be used
      --source-dir='/must-gather/': Set the specific directory on the pod copy the gathered data from.
      --timeout='10m': The length of time to gather data, like 5s, 2m, or 3h, higher than zero. Defaults to 10 minutes.

Use "oc adm options" for a list of global command-line options (applies to all commands).
```

### oc adm global options

Probabaly nothing, but interesting could be:

- server

```
$ oc adm options
The following options can be passed to any command:

      --add-dir-header=false: If true, adds the file directory to the header of the log messages
      --add_dir_header=false: If true, adds the file directory to the header of the log messages
      --alsologtostderr=false: log to standard error as well as files
      --as='': Username to impersonate for the operation
      --as-group=[]: Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --cache-dir='/home/maufart/.kube/cache': Default cache directory
      --certificate-authority='': Path to a cert file for the certificate authority
      --client-certificate='': Path to a client certificate file for TLS
      --client-key='': Path to a client key file for TLS
      --cluster='': The name of the kubeconfig cluster to use
      --context='': The name of the kubeconfig context to use
      --insecure-skip-tls-verify=false: If true, the server's certificate will not be checked for validity. This will
make your HTTPS connections insecure
      --kubeconfig='': Path to the kubeconfig file to use for CLI requests.
      --log-backtrace-at=:0: when logging hits line file:N, emit a stack trace
      --log-dir='': If non-empty, write log files in this directory
      --log-file='': If non-empty, use this log file
      --log-file-max-size=1800: Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0,
the maximum file size is unlimited.
      --log-flush-frequency=5s: Maximum number of seconds between log flushes
      --log_backtrace_at=:0: when logging hits line file:N, emit a stack trace
      --log_dir='': If non-empty, write log files in this directory
      --log_file='': If non-empty, use this log file
      --log_file_max_size=1800: Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0,
the maximum file size is unlimited.
      --loglevel=0: Set the level of log output (0-10)
      --logtostderr=true: log to standard error instead of files
      --match-server-version=false: Require server version to match client version
  -n, --namespace='': If present, the namespace scope for this CLI request
      --one-output=false: If true, only write logs to their native severity level (vs also writing to each lower
severity level
      --one_output=false: If true, only write logs to their native severity level (vs also writing to each lower
severity level
      --request-timeout='0': The length of time to wait before giving up on a single server request. Non-zero values
should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests.
  -s, --server='': The address and port of the Kubernetes API server
      --skip-headers=false: If true, avoid header prefixes in the log messages
      --skip-log-headers=false: If true, avoid headers when opening log files
      --skip_headers=false: If true, avoid header prefixes in the log messages
      --skip_log_headers=false: If true, avoid headers when opening log files
      --stderrthreshold=2: logs at or above this threshold go to stderr
      --tls-server-name='': Server name to use for server certificate validation. If it is not provided, the hostname
used to contact the server is used
      --token='': Bearer token for authentication to the API server
      --user='': The name of the kubeconfig user to use
  -v, --v=0: number for the log level verbosity
      --vmodule=: comma-separated list of pattern=N settings for file-filtered logging
```