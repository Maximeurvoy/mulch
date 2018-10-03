#!/bin/bash

# This script does all the small things needed to install mulchd and
# mulch-proxy on a Linux system. Have a look to main() function for
# important stuff if you need to create a package.

# defaults (see --help)
ETC="/etc/mulch"
VAR_DATA="/var/lib/mulch"
VAR_STORAGE="/srv/mulch"
FORCE="false"

SOURCE=$(dirname "$0")

# TODO:
# create services?
# API key? (generate a new one?)
# check storage accessibility (minimum: --x) for libvirt?

function main() {
    parse_args "$@"

    check_noroot # show warning if UID 0

    check_libvirt_access

    is_dir_writable "$ETC"
    is_dir_writable "$VAR_DATA"
    is_dir_writable "$VAR_STORAGE"

    check_if_existing_config

    copy_config
    gen_ssh_key
    update_config_ssh
}

function check() {
    if [ $1 -ne 0 ]; then
        echo "error, exiting"
        exit $1
    fi
}

function parse_args() {
    while (( "$#" )); do
        case $1 in
        -e|--etc)
            shift
            ETC="$1"
            shift
            ;;
        -d|--data)
            shift
            VAR_DATA="$1"
            shift
            ;;
        -s|--storage)
            shift
            VAR_STORAGE="$1"
            shift
            ;;
        -f|--force)
            FORCE="true"
            shift
            ;;
        -h|--help)
            echo ""
            echo "** Helper script: install mulchd and mulch-proxy **"
            echo ""
            echo "Note: mulch client is not installed/configured by this script."
            echo ""
            echo "Options and defaults (short options available too):"
            echo "  --etc $ETC (-e, config files)"
            echo "  --data $VAR_DATA (-d, state [small] databases)"
            echo "  --storage $VAR_STORAGE (-s, disks storage)"
            echo "  --force (-f, erase old install)"
            exit 1
            ;;
        *)
            echo "Unknown option $1"
            exit 2
            ;;
        esac
    done
}

function is_dir_writable() {
    echo "checking if $1 is writable…"
    if [ ! -d "$1" ]; then
        echo "error: directory $1 does not exists"
        exit 10
    fi
    test_file="$1/.wtest"
    touch "$test_file"
    check $?
    rm -f "$test_file"
}

function check_noroot() {
    uid=$(id -u)
    if [ "$uid" -eq 0 ]; then
        echo "ROOT PRIVILEGES NOT REQUIRED!"
    fi
}

function check_if_existing_config() {
    if [ -f "$ETC/mulchd.toml" ]; then
        echo "Existing configuration found!"
        if [ $FORCE == "false" ]; then
            echo "This script is intentend to do a new install, not to upgrade an existing one."
            echo "If you know what you are doing, you may use --force option."
            echo "Exiting."
            exit 1
        fi
    fi
}

function copy_config() {
    echo "copying config…"
    cp -Rp $SOURCE/etc/* "$ETC/"
    check $?
    mv "$ETC/mulchd.sample.toml" "$ETC/mulchd.toml"
    check $?
}

function gen_ssh_key() {
    echo "generating SSH key…"

    priv_key="$ETC/ssh/id_rsa_mulchd"
    pub_key="$priv_key.pub"

    mkdir -pm 0700 "$ETC/ssh"
    check $?
    if [ $FORCE == "true" ]; then
        rm -f "$priv_key" "$pub_key"
        check $?
    fi
    ssh-keygen -b 4096 -C "admin@vms" -N "" -q -f "$priv_key"
    check $?
}

function check_libvirt_access() {
    echo "checking libvirt access…"
    virsh -c qemu:///system version
    ret=$?
    if [ "$ret" -ne 0 ]; then
        echo "Failed."
        echo " - check that libvirtd is running"
        echo "   - systemd: systemctl status libvirtd"
        echo "   - sysv: service libvirtd status"
        echo " - check that $USER is allowed to connect to qemu:///system"
        echo "   - check that your user is in 'libvirt' group"
        echo "   - some distributions do this automatically on package install"
        echo "   - you may have to disconnect / reconnect your user"
        echo "   - if needed: 'usermod -aG libvirt \$USER'"
    fi
    check $ret
}

function update_config_ssh() {
    r_priv_key=$(realpath "$priv_key")
    r_pub_key=$(realpath "$pub_key")

    sed -i'' "s|^mulch_ssh_private_key =.*|mulch_ssh_private_key = \"$r_priv_key\"|" "$ETC/mulchd.toml"
    sed -i'' "s|^mulch_ssh_public_key =.*|mulch_ssh_public_key = \"$r_pub_key\"|" "$ETC/mulchd.toml"
}

main "$@"