#!/bin/sh

set -e

PAT="pat"

# Check if command is available
has() {
    command -v "$1" >/dev/null 2>&1
}

add_to_path() {
    if has fish;
    then
        fish -c "set -U fish_user_paths $fish_user_paths $1"
    fi

    if has zsh && ! grep -q "$1" ~/.zshenv;
    then
        echo "export PATH=\"$1:\$PATH\"" >> ~/.zshenv
    fi

    if has bash && ! grep -q "$1" ~/.bashrc;
    then
        echo "export PATH=\"$1:\$PATH\"" >> ~/.bashrc
    fi
}

install_completions() {
    test "$OS" = "Darwin" && return 0

    COMPLETIONS="$1"

    if has zsh;
    then
        info "Installing zsh completions"
        TARGET="/usr/share/zsh/vendor-completions"
        sudo mkdir -p "$TARGET"
        sudo cp "$COMPLETIONS/$PAT.zsh" "$TARGET/_$PAT"
    fi

    if has fish;
    then
        info "Installing fish completions"
        TARGET="/usr/share/fish/completions"
        sudo mkdir -p "$TARGET"
        sudo cp "$COMPLETIONS/$PAT.fish" "$TARGET/$PAT.fish"
    fi

    if has bash;
    then
        info "Installing bash completions"
        TARGET="/usr/share/bash-completion/completions"

        sudo mkdir -p $TARGET
        sudo cp "$COMPLETIONS/$PAT.bash" "$TARGET/$PAT"
    fi
}

info() {
    printf "\033[1;35m%s\033[0m %s\n" "info:" "$1"
}

warn() {
    printf "\033[1;33m%s\033[0m %s\n" "warning:" "$1"
}

err() {
    printf "\033[1;31m%s\033[0m %s\n" "error:" "$1" >&2
}

die() {
    err "$1"
    exit 1
}

# Download file. first argument is out path, second is URL
download() {
    if has curl; then
        curl -sfLo "$1" "$2"
    elif has wget; then
        wget -qO "$1" "$2"
    else
        die "No download program (curl, wget) found, please install one."
    fi
}

verify_checksums() {
    info "Verifying checksums"
    if has sha256sum; then
        OK=$(sha256sum --ignore-missing --quiet --check checksums.txt)
    else
        OK=$(shasum -a 256 --ignore-missing --quiet --check checksums.txt)
    fi

    $OK || die "Checksums did not match! Abort"
}


install_binary() {
    case "$ARCH" in
    aarch64)
        ARCH="arm64"
        ;;
    armv*)
        ARCH="armv6"
        ;;
    amd64)
        ARCH="x86_64"
        ;;
    i*)
        ARCH="i386"
        ;;
    esac

    TAR_NAME="${FILE_BASENAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
    TAR_FILE="$TMPDIR/$TAR_NAME"

    export TAR_NAME TAR_FILE

    (
        cd "$TMPDIR"

        info "Downloading $TAR_NAME"
        download "$TAR_FILE" "$RELEASES_URL/download/$TAG/$TAR_NAME"

        info "Downloading checksums"
        download "checksums.txt" "$RELEASES_URL/download/$TAG/checksums.txt"

        verify_checksums
    )

    tar -xf "$TAR_FILE" -C "$TMPDIR"

    if has pat;
    then
        # get directory name of installed pat
        OUT=$(dirname "$(command -v pat)")
    else
        OUT="$HOME/.$PAT/bin"
    fi

    mkdir -p "$OUT"
    info "Moving to ${OUT}"
    sh -c "install -m755 '$TMPDIR/$FILE_BASENAME' '${OUT}'"
    add_to_path "$OUT"

    install_completions "$TMPDIR/completions"
}

pre_install() {
    RELEASES_URL="https://github.com/metafates/$PAT/releases"
    FILE_BASENAME="$PAT"

    info "Fetching latest version"
    TAG="$(curl -sfL -o /dev/null -w "%{url_effective}" "$RELEASES_URL/latest" |
        rev |
        cut -f1 -d'/' |
        rev)"

    test -z "$TAG" && {
        die "Unable to get $PAT version."
    }

    # test if tag is a semver
    echo "$TAG" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$' || {
        die "Unable to get $PAT version."
    }

    OS=$(uname -s)
    ARCH=$(uname -m)
    VERSION=${TAG#?}

    info "Latest version is $VERSION"

    TMPDIR="$(mktemp -d)"
}

post_install() {
    printf "\n🎉 \033[1;32m%s was installed successfully\033[0m\n\n" "$PAT"
}

install_pat() {
    pre_install
    install_binary
    post_install
}

install_pat
