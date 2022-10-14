import pathlib as pl

APP = "pat"
OUT = ".goreleaser.yaml"

# get script path
script_path = pl.Path(__file__).resolve()

# get project root (3 levels up)
project_root = script_path.parent.parent.parent

goreleaser_yaml = f'''
before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
    goarch:
      - "386"
      - amd64
      - arm
      - arm64
    flags:
      - -trimpath
    ldflags:
      - -s
      - -w

archives:
  - replacements:
      darwin: Darwin
      linux: Linux
      android: Android
      386: i386
      amd64: x86_64
    files:
      - completions/*
      - README.md
      - LICENSE

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{{{ incpatch .Version }}}}-next"

changelog:
  sort: asc
  use: github
  groups:
    - title: Dependency updates
      regexp: "^.*feat\\(deps\\)*:+.*$"
      order: 300
    - title: 'New Features'
      regexp: "^.*feat[(\\w)]*:+.*$"
      order: 100
    - title: 'Bug fixes'
      regexp: "^.*fix[(\\w)]*:+.*$"
      order: 200
    - title: 'Documentation updates'
      regexp: "^.*docs[(\\w)]*:+.*$"
      order: 400
    - title: Other work
      order: 9999
  filters:
    exclude:
      - '^test'
      - '^chore'
      - '^refactor'
      - '^build'
      - 'merge conflict'
      - Merge pull request
      - Merge remote-tracking branch
      - Merge branch
      - go mod tidy


release:
  github:
    owner: metafates
    name: {APP}

  name_template: "{{{{.ProjectName}}}} v{{{{.Version}}}}"
  header: |
    To install:
    ```sh
    curl -sSL {APP}.metafates.one/install | sh
    ```

    ## What's new?

  footer: |

    **Full Changelog**: https://github.com/metafates/{APP}/compare/{{{{ .PreviousTag }}}}...{{{{ .Tag }}}}

    ---

    Bugs? Suggestions? [Open an issue](https://github.com/metafates/{APP}/issues/new/choose)
'''

with open(pl.Path(project_root, OUT), "w") as file:
    file.write(goreleaser_yaml)
