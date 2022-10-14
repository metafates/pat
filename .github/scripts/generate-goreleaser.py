import pathlib as pl

APP = "pat"
DESC = "The $PATH manager"

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
      - windows
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
      windows: Windows
      android: Android
      386: i386
      amd64: x86_64
    format_overrides:
      - goos: windows
        format: zip
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

brews:
  - name: {APP}

    tap:
      owner: metafates
      name: homebrew-{APP}
      branch: main
      token: "{{{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}}}"

    commit_author:
      name: goreleaserbot
      email: bot@goreleaser.com

    commit_msg_template: "Brew formula update for {{{{ .ProjectName }}}} version {{{{ .Tag }}}}"
    homepage: "https://github.com/metafates/{APP}"
    description: "{DESC}"
    license: "MIT"
    skip_upload: false

    test: |
      system "#{{bin}}/{APP} -v"

    install: |-
      bin.install "{APP}"
      bash_completion.install "completions/{APP}.bash" => "{APP}"
      zsh_completion.install "completions/{APP}.zsh" => "_{APP}"
      fish_completion.install "completions/{APP}.fish"

scoop:
  bucket:
    owner: metafates
    name: scoop-metafates
    branch: main
    token: "{{{{ .Env.SCOOP_TAP_GITHUB_TOKEN }}}}"

  folder: bucket

  commit_author:
    name: goreleaserbot
    email: bot@goreleaser.com

  commit_msg_template: "Scoop update for {{{{ .ProjectName }}}} version {{{{ .Tag }}}}"
  homepage: "https://github.com/metafates/{APP}"
  description: "{DESC}"
  license: MIT
  skip_upload: false


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

nfpms:
  - file_name_template: "{{{{ .ConventionalFileName }}}}"
    homepage: https://github.com/metafates/{APP}
    maintainer: metafates <fates@duck.com>
    description: |-
{DESC} 

    license: MIT
    formats:
      - deb
      - rpm

    bindir: /usr/local/bin
    section: utils

    deb:
      lintian_overrides:
        - statically-linked-binary
        - changelog-file-missing-in-native-package

    contents:
      - src: ./completions/{APP}.bash
        dst: /usr/share/bash-completion/completions/{APP}
        file_info:
          mode: 0644
      - src: ./completions/{APP}.fish
        dst: /usr/share/fish/completions/{APP}.fish
        file_info:
          mode: 0644
      - src: ./completions/{APP}.zsh
        dst: /usr/share/zsh/vendor-completions/_{APP}
        file_info:
          mode: 0644
      - src: ./LICENSE
        dst: /usr/share/doc/{APP}/copyright
        file_info:
          mode: 0644
'''

with open(pl.Path(project_root, OUT), "w") as file:
    file.write(goreleaser_yaml)
