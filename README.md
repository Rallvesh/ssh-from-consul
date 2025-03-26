# SSH From Consul CLI Utility

This is a command-line utility that allows you to connect to servers via SSH, with information about the servers retrieved from Consul. The utility supports:

- Listing available nodes in a Consul cluster with table output.
- Connecting to a specified node via SSH.

## Features

- **Table output**: Use the `ls` command to list all nodes in the Consul cluster in a tabular format via the charmbracelet library.
- **Direct SSH connection**: Use the `connect` command to directly SSH into a specified node.
- **Configuration file**: On the first run, a configuration file is created at `~/.config/sfc/sfc.json` with default values. This file can be edited to add custom Consul servers.

## Installation

You can build the binary using the provided `Makefile` or download pre-built binaries for different architectures.

### Build from Source

Clone this repository:

```bash
git clone https://github.com/your-repo/ssh-from-consul.git
cd ssh-from-consul
```

Then, run one of the following commands depending on your platform:

- For Linux (AMD64):
```bash
make build_amd
```

- For macOS (ARM64):
```bash
make build_arm
```

## Run the Utility

After building the binary, you can run the utility with the following commands:

- List nodes:
```bash
./bin/darwin_arm64/sfc <cluster_name> ls
```

- Connect to a node:
```bash
./bin/darwin_arm64/sfc <cluster_name> connect <hostname>
```

## Configuration

The configuration file ~/.config/sfc/sfc.json is automatically created the first time you run the utility. The default configuration looks like this:
```json
[
  {
    "default": {
      "consul_http_addr": "http://127.0.0.1:8500",
      "consul_http_token": "",
      "private_key_path": "",
      "username": ""
    }
  }
]
```

You can edit this file to add or modify Consul clusters. The configuration options are:

    consul_http_addr: The address of your Consul server.

    consul_http_token: The authentication token for Consul.

    private_key_path: The path to your private SSH key (optional).

    username: The SSH username for connecting to the nodes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

```vbnet
Let me know if you'd like to make any changes!
```