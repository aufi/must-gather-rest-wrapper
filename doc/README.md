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

- pending
- completed
- error

## Open questions

Sqlite backend or just in memory? (after container restart there will be no archives anyway so no need for persistent db)
In memory and communicate via chans