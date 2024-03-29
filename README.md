# systemd-cd

## Build

```bash
go build .
```

## Usage

```bash
./systemd-cd \
  -f <manifest>.toml \
  --webapi.jwt.secret "your jwt secret" \
  --webapi.username "web api username" \ # default is "admin"
  --webapi.password "web api password"
```

You can specify multiple manifest files.

```bash
./systemd-cd -f <manifest>.toml -f <manifest>.toml
```

### GitOps

Configure a GitOps CD system by specifying a repository containing manifest files.

```bash
./systemd-cd \
  -f <manifest>.toml \
  --ops.git-remote "your git-ops repository url" \
  --webapi.jwt.secret "your jwt secret" \
  --webapi.username "web api username" \ # default is "admin"
  --webapi.password "web api password"
```

### Configuration

#### Sample `go`

```toml
name = "prometheus_sh_exporter"
git_remote_url = "https://github.com/tingtt/prometheus_sh_exporter.git"
git_target_branch = "main"
git_tag_regex = "v*" # e.g. "v1.0.0"
build_commands = ["/usr/bin/go build"]
binaries = ["prometheus_sh_exporter"]

[[systemd_services]]
name = "prometheus_sh_exporter"
description = "The shell exporter allows probing with shell scripts."
exec_start = "./prometheus_sh_exporter"
args = "--port 9923"
port = 9923

[[systemd_services.etc]]
target = "prometheus_sh_exporter.yml"
content = """
"""
option = "-config.file"
```

It runs like this

```bash
/usr/local/systemd-cd/bin/prometheus_sh_exporter/prometheus_sh_exporter \
  --port 9923 \
  -config.file /usr/local/systemd-cd/etc/prometheus_sh_exporter/prometheus_sh_exporter.yml
```

with `/usr/local/lib/systemd/system/prometheus_sh_exporter.service`

#### Sample `Next.js`

```toml
name = "nextjs-workspace"

git_remote_url = "https://github.com/tingtt/workspace-nextjs.git"
git_target_branch = "main"
build_commands = ["/root/.local/share/pnpm/pnpm install", "/root/.local/share/pnpm/pnpm build"]

[[systemd_services]]
name = "nextjs-workspace"
description = "Next.js sample"
exec_start_pre = "/root/.local/share/pnpm/pnpm install next"
exec_start = "/root/.local/share/pnpm/pnpm start"
args = "--port 3000"
opt_files = [".next/", "package.json", "public/"]
port = 3000
```

It runs like this

```bash
/usr/local/bin/yarn start --port 3000
```

with `/usr/local/lib/systemd/system/nextjs-workspace.service`
