# Zugang CLI Documentation

Zugang is a command-line interface (CLI) tool designed to securely connect to remote hosts via SSH using credentials stored in your Bitwarden vault.

## Installation

To install Zugang, follow these steps:

1. Clone the repository:
   ```
   git clone <repository_url>
   ```

2. Navigate to the project directory:
   ```
   cd zugang
   ```

3. Build the application:
   ```
   go build
   ```

4. Optionally, you can move the binary to a directory in your PATH to make it globally accessible.

## Usage

Zugang provides two main commands: `login` and `sync`.

### Login

The `login` command enables you to connect to a remote host using credentials from your Bitwarden vault.

```
zugang login <host> [flags]
```

#### Flags

- `--user`, `-u`: Specify a specific username when connecting to the remote host.
- `--hostKeyCheck`: Enable or disable host key checks when connecting to the remote host.

#### Examples

- To connect to a remote host named "example.com":
  ```
  zugang login example.com
  ```

- To specify a specific username when connecting to a remote host:
  ```
  zugang login example.com --user myusername
  ```

- To disable host key checks when connecting to a remote host (use with caution):
  ```
  zugang login example.com --hostKeyCheck=false
  ```

### Sync

The `sync` command pulls the latest vault data from the Bitwarden server.

```
zugang sync
```

## License

Zugang is licensed under the [GLWTS License](LICENSE).
