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

## Example testing API calls

### Trigger must-gather

```
$ curl -X POST -H "Content-Type: application/json" -d '{"image": "quay.io/konveyor/forklift-must-gather"}' http://localhost:8080/must-gather
{"ID":11,"CreatedAt":"2021-06-22T15:32:11.204910808+02:00","UpdatedAt":"2021-06-22T15:32:11.204910808+02:00","DeletedAt":null,"status":"new","image":"quay.io/konveyor/forklift-must-gather","CustomCommand":"","ArchivePath":"","ExecOutput":""}
```

### Get must-gathering

```
$ curl http://localhost:8080/must-gather/11 | jq
{
  "ID": 11,
  "CreatedAt": "2021-06-22T15:32:11.204910808+02:00",
  "UpdatedAt": "2021-06-22T15:33:27.478623067+02:00",
  "DeletedAt": null,
  "status": "completed",
  "image": "quay.io/konveyor/forklift-must-gather",
  "CustomCommand": "",
  "ArchivePath": "/tmp/must-gather-result-11/quay-io-konveyor-forklift-must-gather-sha256-a80be15dfc0bb358a56b7d11d2c77f3dbd78b51c27b7e24c4ea604c7697e3b43/must-gather.tar.gz\n",
  "ExecOutput": "[must-gather      ] OUT Using must-gather plug-in image: quay.io/konveyor/forklift-must-gather\n[must-gather      ] OUT namespace/openshift-must-gather-cfqzw created\n[must-gather      ] OUT clusterrolebinding.rbac.authorization.k8s.io/must-gather-8j5zh created\n[must-gather      ] OUT pod for plug-in image quay.io/konveyor/forklift-must-gather created\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.645547576Z Gathering Forklift and Kubevirt related CRs for namespaces [openshift-mtv]\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.646037388Z Gathering logs for namespaces [openshift-mtv]\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.646351980Z Waiting for background gather tasks to finish\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.884914879Z [ns=openshift-mtv][pod=forklift-controller-7b8bc54cbd-7nzr5] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.885111752Z [ns=openshift-mtv][pod=forklift-controller-7b8bc54cbd-7nzr5] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.891176427Z [ns=openshift-mtv][pod=forklift-operator-56c68ffdf5-qg6pz] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.891406364Z [ns=openshift-mtv][pod=forklift-operator-56c68ffdf5-qg6pz] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.898434710Z [ns=openshift-mtv][pod=forklift-ui-766b879d9-crdnf] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.898624338Z [ns=openshift-mtv][pod=forklift-ui-766b879d9-crdnf] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.920465335Z [ns=openshift-mtv][pod=forklift-validation-644b94856c-jp5tt] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.920644948Z [ns=openshift-mtv][pod=forklift-validation-644b94856c-jp5tt] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.981588561Z [ns=openshift-mtv][pod=virt-launcher-marek-vm-rocky-c8tz4] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.982010207Z [ns=openshift-mtv][pod=virt-launcher-marek-vm-rocky-c8tz4] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:47.329025393Z Starting collection of: [forkliftcontrollers.forklift.konveyor.io hooks.forklift.konveyor.io hosts.forklift.konveyor.io migrations.forklift.konveyor.io networkmaps.forklift.konveyor.io plans.forklift.konveyor.io providers.forklift.konveyor.io provisioners.forklift.konveyor.io storagemaps.forklift.konveyor.io]\n[must-gather-8kxcw] POD 2021-06-22T13:30:47.329062924Z Collecting forkliftcontrollers.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.500310051Z [ns=openshift-cnv][pod=cdi-apiserver-7cc5c9ff96-d5r84] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.500310051Z [ns=openshift-cnv][pod=cdi-apiserver-7cc5c9ff96-d5r84] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.507008640Z [ns=openshift-cnv][pod=cdi-deployment-787fd94f86-wxhc4] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.507222982Z [ns=openshift-cnv][pod=cdi-deployment-787fd94f86-wxhc4] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.513519506Z [ns=openshift-cnv][pod=cdi-operator-95bbd576b-fjdt6] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.513770140Z [ns=openshift-cnv][pod=cdi-operator-95bbd576b-fjdt6] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.538349228Z [ns=openshift-cnv][pod=cdi-uploadproxy-85c6b74dcc-s4dwh] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.538388261Z [ns=openshift-cnv][pod=cdi-uploadproxy-85c6b74dcc-s4dwh] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.928562055Z Collecting hooks.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:49.010198696Z [ns=openshift-cnv][pod=vm-import-operator-5847c7f656-ld6js] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:49.010466149Z [ns=openshift-cnv][pod=vm-import-operator-5847c7f656-ld6js] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:50.105634416Z Collecting hosts.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:51.686522689Z Collecting migrations.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:52.634266103Z Collecting networkmaps.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:54.009785229Z Collecting plans.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:55.262569392Z Collecting providers.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:56.552077343Z Collecting provisioners.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:01.875836141Z Collecting storagemaps.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:03.505230045Z Starting collection of: [forkliftcontrollers.forklift.konveyor.io hooks.forklift.konveyor.io hosts.forklift.konveyor.io migrations.forklift.konveyor.io networkmaps.forklift.konveyor.io plans.forklift.konveyor.io providers.forklift.konveyor.io provisioners.forklift.konveyor.io storagemaps.forklift.konveyor.io cdiconfigs.cdi.kubevirt.io cdis.cdi.kubevirt.io datavolumes.cdi.kubevirt.io hostpathprovisioners.hostpathprovisioner.kubevirt.io hyperconvergeds.hco.kubevirt.io kubevirts.kubevirt.io networkaddonsconfigs.networkaddonsoperator.network.kubevirt.io nodemaintenances.nodemaintenance.kubevirt.io ovirtproviders.v2v.kubevirt.io ssps.ssp.kubevirt.io v2vvmwares.v2v.kubevirt.io virtualmachineimports.v2v.kubevirt.io virtualmachineinstancemigrations.kubevirt.io virtualmachineinstancepresets.kubevirt.io virtualmachineinstancereplicasets.kubevirt.io virtualmachineinstances.kubevirt.io virtualmachinerestores.snapshot.kubevirt.io virtualmachines.kubevirt.io virtualmachinesnapshotcontents.snapshot.kubevirt.io virtualmachinesnapshots.snapshot.kubevirt.io vmimportconfigs.v2v.kubevirt.io]\n[must-gather-8kxcw] POD 2021-06-22T13:31:03.505387499Z Collecting forkliftcontrollers.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:04.882659332Z Collecting hooks.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:06.083539576Z Collecting hosts.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:07.189612576Z Collecting migrations.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:08.357078919Z Collecting networkmaps.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:09.630711612Z Collecting plans.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:10.848809461Z Collecting providers.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:12.211409068Z Collecting provisioners.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:17.664170854Z Collecting storagemaps.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:18.748687665Z Collecting cdiconfigs.cdi.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:20.041747643Z Collecting cdis.cdi.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:21.446744482Z Collecting datavolumes.cdi.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:22.786888526Z Collecting hostpathprovisioners.hostpathprovisioner.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:24.603953498Z Collecting hyperconvergeds.hco.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:25.665915425Z Collecting kubevirts.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:27.262465812Z Collecting networkaddonsconfigs.networkaddonsoperator.network.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:29.280646306Z Collecting nodemaintenances.nodemaintenance.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:30.787414576Z Collecting ovirtproviders.v2v.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:31.745572061Z Collecting ssps.ssp.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:32.761109960Z Collecting v2vvmwares.v2v.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:34.242837733Z Collecting virtualmachineimports.v2v.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:35.393366026Z Collecting virtualmachineinstancemigrations.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:37.145331775Z Collecting virtualmachineinstancepresets.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:38.411380451Z Collecting virtualmachineinstancereplicasets.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:39.340117397Z Collecting virtualmachineinstances.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:40.489982652Z Collecting virtualmachinerestores.snapshot.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:41.827559432Z Collecting virtualmachines.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:42.917761088Z Collecting virtualmachinesnapshotcontents.snapshot.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:44.572503985Z Collecting virtualmachinesnapshots.snapshot.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:45.505710157Z Collecting vmimportconfigs.v2v.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:46.541691736Z Tarring must-gather artifacts...\n[must-gather-8kxcw] POD 2021-06-22T13:31:46.547138473Z tar: 2021-06-22T13:31:46.547209576Z Removing leading `/' from member names2021-06-22T13:31:46.547235042Z \n[must-gather-8kxcw] POD 2021-06-22T13:31:47.333120540Z Created /must-gather/must-gather.tar.gz\n[must-gather-8kxcw] POD 2021-06-22T13:31:47.333340523Z Waiting for copy phase...\n[must-gather-8kxcw] OUT waiting for gather to complete\n[must-gather-8kxcw] OUT downloading gather output\n[must-gather-8kxcw] OUT receiving incremental file list\n[must-gather-8kxcw] OUT ./\n[must-gather-8kxcw] OUT must-gather.tar.gz\n[must-gather-8kxcw] OUT \n[must-gather-8kxcw] OUT sent 46 bytes  received 2,748,363 bytes  366,454.53 bytes/sec\n[must-gather-8kxcw] OUT total size is 2,747,106  speedup is 1.00\n[must-gather      ] OUT clusterrolebinding.rbac.authorization.k8s.io/must-gather-8j5zh deleted\n[must-gather      ] OUT namespace/openshift-must-gather-cfqzw deleted\n"
}
```

### List must-gatherings

```
$ curl http://localhost:8080/must-gather | jq
[
  {
    "ID": 11,
    "CreatedAt": "2021-06-22T15:32:11.204910808+02:00",
    "UpdatedAt": "2021-06-22T15:33:27.478623067+02:00",
    "DeletedAt": null,
    "status": "completed",
    "image": "quay.io/konveyor/forklift-must-gather",
    "CustomCommand": "",
    "ArchivePath": "/tmp/must-gather-result-11/quay-io-konveyor-forklift-must-gather-sha256-a80be15dfc0bb358a56b7d11d2c77f3dbd78b51c27b7e24c4ea604c7697e3b43/must-gather.tar.gz\n",
    "ExecOutput": "[must-gather      ] OUT Using must-gather plug-in image: quay.io/konveyor/forklift-must-gather\n[must-gather      ] OUT namespace/openshift-must-gather-cfqzw created\n[must-gather      ] OUT clusterrolebinding.rbac.authorization.k8s.io/must-gather-8j5zh created\n[must-gather      ] OUT pod for plug-in image quay.io/konveyor/forklift-must-gather created\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.645547576Z Gathering Forklift and Kubevirt related CRs for namespaces [openshift-mtv]\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.646037388Z Gathering logs for namespaces [openshift-mtv]\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.646351980Z Waiting for background gather tasks to finish\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.884914879Z [ns=openshift-mtv][pod=forklift-controller-7b8bc54cbd-7nzr5] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.885111752Z [ns=openshift-mtv][pod=forklift-controller-7b8bc54cbd-7nzr5] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.891176427Z [ns=openshift-mtv][pod=forklift-operator-56c68ffdf5-qg6pz] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.891406364Z [ns=openshift-mtv][pod=forklift-operator-56c68ffdf5-qg6pz] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.898434710Z [ns=openshift-mtv][pod=forklift-ui-766b879d9-crdnf] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.898624338Z [ns=openshift-mtv][pod=forklift-ui-766b879d9-crdnf] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.920465335Z [ns=openshift-mtv][pod=forklift-validation-644b94856c-jp5tt] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.920644948Z [ns=openshift-mtv][pod=forklift-validation-644b94856c-jp5tt] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.981588561Z [ns=openshift-mtv][pod=virt-launcher-marek-vm-rocky-c8tz4] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:46.982010207Z [ns=openshift-mtv][pod=virt-launcher-marek-vm-rocky-c8tz4] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:47.329025393Z Starting collection of: [forkliftcontrollers.forklift.konveyor.io hooks.forklift.konveyor.io hosts.forklift.konveyor.io migrations.forklift.konveyor.io networkmaps.forklift.konveyor.io plans.forklift.konveyor.io providers.forklift.konveyor.io provisioners.forklift.konveyor.io storagemaps.forklift.konveyor.io]\n[must-gather-8kxcw] POD 2021-06-22T13:30:47.329062924Z Collecting forkliftcontrollers.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.500310051Z [ns=openshift-cnv][pod=cdi-apiserver-7cc5c9ff96-d5r84] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.500310051Z [ns=openshift-cnv][pod=cdi-apiserver-7cc5c9ff96-d5r84] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.507008640Z [ns=openshift-cnv][pod=cdi-deployment-787fd94f86-wxhc4] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.507222982Z [ns=openshift-cnv][pod=cdi-deployment-787fd94f86-wxhc4] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.513519506Z [ns=openshift-cnv][pod=cdi-operator-95bbd576b-fjdt6] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.513770140Z [ns=openshift-cnv][pod=cdi-operator-95bbd576b-fjdt6] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.538349228Z [ns=openshift-cnv][pod=cdi-uploadproxy-85c6b74dcc-s4dwh] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.538388261Z [ns=openshift-cnv][pod=cdi-uploadproxy-85c6b74dcc-s4dwh] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:48.928562055Z Collecting hooks.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:49.010198696Z [ns=openshift-cnv][pod=vm-import-operator-5847c7f656-ld6js] Collecting Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:49.010466149Z [ns=openshift-cnv][pod=vm-import-operator-5847c7f656-ld6js] Collecting previous Pod logs...\n[must-gather-8kxcw] POD 2021-06-22T13:30:50.105634416Z Collecting hosts.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:51.686522689Z Collecting migrations.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:52.634266103Z Collecting networkmaps.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:54.009785229Z Collecting plans.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:55.262569392Z Collecting providers.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:30:56.552077343Z Collecting provisioners.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:01.875836141Z Collecting storagemaps.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:03.505230045Z Starting collection of: [forkliftcontrollers.forklift.konveyor.io hooks.forklift.konveyor.io hosts.forklift.konveyor.io migrations.forklift.konveyor.io networkmaps.forklift.konveyor.io plans.forklift.konveyor.io providers.forklift.konveyor.io provisioners.forklift.konveyor.io storagemaps.forklift.konveyor.io cdiconfigs.cdi.kubevirt.io cdis.cdi.kubevirt.io datavolumes.cdi.kubevirt.io hostpathprovisioners.hostpathprovisioner.kubevirt.io hyperconvergeds.hco.kubevirt.io kubevirts.kubevirt.io networkaddonsconfigs.networkaddonsoperator.network.kubevirt.io nodemaintenances.nodemaintenance.kubevirt.io ovirtproviders.v2v.kubevirt.io ssps.ssp.kubevirt.io v2vvmwares.v2v.kubevirt.io virtualmachineimports.v2v.kubevirt.io virtualmachineinstancemigrations.kubevirt.io virtualmachineinstancepresets.kubevirt.io virtualmachineinstancereplicasets.kubevirt.io virtualmachineinstances.kubevirt.io virtualmachinerestores.snapshot.kubevirt.io virtualmachines.kubevirt.io virtualmachinesnapshotcontents.snapshot.kubevirt.io virtualmachinesnapshots.snapshot.kubevirt.io vmimportconfigs.v2v.kubevirt.io]\n[must-gather-8kxcw] POD 2021-06-22T13:31:03.505387499Z Collecting forkliftcontrollers.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:04.882659332Z Collecting hooks.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:06.083539576Z Collecting hosts.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:07.189612576Z Collecting migrations.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:08.357078919Z Collecting networkmaps.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:09.630711612Z Collecting plans.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:10.848809461Z Collecting providers.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:12.211409068Z Collecting provisioners.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:17.664170854Z Collecting storagemaps.forklift.konveyor.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:18.748687665Z Collecting cdiconfigs.cdi.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:20.041747643Z Collecting cdis.cdi.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:21.446744482Z Collecting datavolumes.cdi.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:22.786888526Z Collecting hostpathprovisioners.hostpathprovisioner.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:24.603953498Z Collecting hyperconvergeds.hco.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:25.665915425Z Collecting kubevirts.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:27.262465812Z Collecting networkaddonsconfigs.networkaddonsoperator.network.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:29.280646306Z Collecting nodemaintenances.nodemaintenance.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:30.787414576Z Collecting ovirtproviders.v2v.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:31.745572061Z Collecting ssps.ssp.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:32.761109960Z Collecting v2vvmwares.v2v.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:34.242837733Z Collecting virtualmachineimports.v2v.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:35.393366026Z Collecting virtualmachineinstancemigrations.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:37.145331775Z Collecting virtualmachineinstancepresets.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:38.411380451Z Collecting virtualmachineinstancereplicasets.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:39.340117397Z Collecting virtualmachineinstances.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:40.489982652Z Collecting virtualmachinerestores.snapshot.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:41.827559432Z Collecting virtualmachines.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:42.917761088Z Collecting virtualmachinesnapshotcontents.snapshot.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:44.572503985Z Collecting virtualmachinesnapshots.snapshot.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:45.505710157Z Collecting vmimportconfigs.v2v.kubevirt.io\n[must-gather-8kxcw] POD 2021-06-22T13:31:46.541691736Z Tarring must-gather artifacts...\n[must-gather-8kxcw] POD 2021-06-22T13:31:46.547138473Z tar: 2021-06-22T13:31:46.547209576Z Removing leading `/' from member names2021-06-22T13:31:46.547235042Z \n[must-gather-8kxcw] POD 2021-06-22T13:31:47.333120540Z Created /must-gather/must-gather.tar.gz\n[must-gather-8kxcw] POD 2021-06-22T13:31:47.333340523Z Waiting for copy phase...\n[must-gather-8kxcw] OUT waiting for gather to complete\n[must-gather-8kxcw] OUT downloading gather output\n[must-gather-8kxcw] OUT receiving incremental file list\n[must-gather-8kxcw] OUT ./\n[must-gather-8kxcw] OUT must-gather.tar.gz\n[must-gather-8kxcw] OUT \n[must-gather-8kxcw] OUT sent 46 bytes  received 2,748,363 bytes  366,454.53 bytes/sec\n[must-gather-8kxcw] OUT total size is 2,747,106  speedup is 1.00\n[must-gather      ] OUT clusterrolebinding.rbac.authorization.k8s.io/must-gather-8j5zh deleted\n[must-gather      ] OUT namespace/openshift-must-gather-cfqzw deleted\n"
  }
]
```

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