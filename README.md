# obsidian-clippings-feed

![publish-docker-image](https://github.com/nakatanakatana/obsidian-clippings-feed/actions/workflows/publish-docker-image.yaml/badge.svg)
![CI](https://github.com/nakatanakatana/obsidian-clippings-feed/actions/workflows/ci.yaml/badge.svg)
![Coverage](https://github.com/nakatanakatana/octocov-central/blob/main/badges/nakatanakatana/obsidian-clippings-feed/coverage.svg?raw=true)
![Code to Test Ratio](https://github.com/nakatanakatana/octocov-central/blob/main/badges/nakatanakatana/obsidian-clippings-feed/ratio.svg?raw=true)
![Test Execution Time](https://github.com/nakatanakatana/octocov-central/blob/main/badges/nakatanakatana/obsidian-clippings-feed/time.svg?raw=true)

## LiveSync Bridge

The `livesync-bridge/` submodule provides the Docker image for synchronising
Obsidian vaults with Self-hosted LiveSync. Its Compose file bind-mounts
`livesync-bridge/data` and `livesync-bridge/dat` into the container.

The current image runs as the unprivileged `deno` user (UID/GID `1993:1993`).
On Linux, prepare the bind-mounted directories before starting the service:

```bash
mkdir -p livesync-bridge/data livesync-bridge/dat
sudo chown -R 1993:1993 livesync-bridge/data livesync-bridge/dat
```

Keep `livesync-bridge/dat/config.json` backed up separately; it contains the
bridge configuration and credentials. The Compose named volume
`bridge_local_storage` stores Deno's offline-scan state and should be retained
when recreating the container.
