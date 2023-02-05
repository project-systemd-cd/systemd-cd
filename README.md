# systemd-cd

## Usage

```bash
./systemd-cd -f <manifest>.toml
```

You can specify multiple manifest files.

```bash
./systemd-cd -f <manifest>.toml -f <manifest>.toml
```

## Configuration

### Sample `go`

```toml
name = "prometheus_sh_exporter"
git_remote_url = "https://github.com/tingtt/prometheus_sh_exporter.git"
git_target_branch = "main"
git_tag_regex = "v*" # e.g. "v1.0.0"
build_commands = ["/usr/bin/go build"]
binaries = ["prometheus_sh_exporter"]

[[systemd]]
name = "prometheus_sh_exporter"
description = "The shell exporter allows probing with shell scripts."
exec_start = "./prometheus_sh_exporter"
args = "--port 9923"
port = 9923

[[systemd.etc]]
target = "prometheus_sh_exporter.yml"
option = "-config.file"
```

It runs like this

```bash
/usr/local/systemd-cd/bin/prometheus_sh_exporter/prometheus_sh_exporter \
  --port 9923 \
  -config.file /usr/local/systemd-cd/etc/prometheus_sh_exporter/prometheus_sh_exporter.yml
```

with `/usr/local/lib/systemd/system/prometheus_sh_exporter.service`

### Sample `Next.js`

```toml
name = "tingtt_web_site"

git_remote_url = "https://github.com/tingtt/tingtt.git"
git_target_branch = "main"
git_tag_regex = "v*"

build_commands = ["/usr/local/bin/yarn install && /usr/local/bin/yarn build"]

[[systemd]]
name = "tingtt_web_site"
description = "tingtt's portfolio with Next.js"
exec_start = "/usr/local/bin/yarn start"
args = "--port 3000"
opt_files = [".next/", "node_modules/", "package.json", "public/"]
port = 3000
```

It runs like this

```bash
/usr/local/bin/yarn start --port 3000
```

with `/usr/local/lib/systemd/system/tingtt_web_site.service`
