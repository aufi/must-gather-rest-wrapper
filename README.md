# Must-gather REST wrapper

The wraper should provide HTTP API to allow trigger OpenShift must-gather (full or targeted) and provide output archive again via HTTP.

Under initial development, if this will work well, it should be moved under Forklift organization

Checkout [doc](doc/README.md) directory for more information.

## Notes

### Initial development progress

- <del>prepare HTTP endpoints in gin-gonic</del>
- <del>prepare database/storage for gatherings metadata</del>
- <del>prepare image build scripts</del>
- <del>implement create / get / list gatherings (into db)</del>
- <del>implement raw single must-gather execution based on gathering db record</del>
- <del>implement all needed oc adm must-gather options support</del>
- add ENV variables-driven default must-gather options values
- if needed: handle async/not-blocking gathering with sane limits (e.g. max 10 simul.gatherings)
- prepare serving of gathered archive
- add ocp auth for gin-gonic if needed
- use passed admin token from operator to exec must-gather
- basic tests
- probably: configurable port&db_path&archive storage
