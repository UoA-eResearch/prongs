import argparse
import ipaddress
import os

from . import config
from .scanners import accessible_db, accessible_rdp, password_ssh


def run(target_hosts: list[ipaddress.IPv4Address]) -> None:
    for scan_type, enabled in config.scan_types_enabled.items():
        if not enabled:
            continue

        if scan_type == "accessible-db":
            accessible_db.run(target_hosts)
        elif scan_type == "accessible-rdp":
            accessible_rdp.run(target_hosts)
        elif scan_type == "password-ssh":
            password_ssh.run(target_hosts)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-p", "--pretty-print",
        default=config.pretty_print,
        action="store_true",
        help="Enable pretty print output",
    )
    parser.add_argument(
        "-v", "--version",
        action="version",
        version=config.version,
        help="Show version number",
    )

    target_group = parser.add_mutually_exclusive_group()
    target_group.add_argument(
        "-f", "--file",
        help="Set target/s CIDR from file (line separated)",
    )
    target_group.add_argument(
        "-t", "--targets",
        help="Set target/s CIDR as argument (comma-separated list)",
    )
    target_group.add_argument(
        "-e", "--envvars",
        action="store_true",
        help="Set target/s CIDR from env vars (comma-separated list)",
    )

    scan_group = parser.add_mutually_exclusive_group()
    scan_group.add_argument(
        "-a", "--enable-all",
        action="store_true",
        help="Enable all scan types",
    )
    scan_group.add_argument(
        "-s", "--enable-scan",
        choices=config.scan_types_enabled.keys(),
        action="append",
        help="Enable specific scan type/s",
    )
    args = parser.parse_args()

    if args.pretty_print:
        config.pretty_print = True

    # Handle required scan types arguments
    if not args.enable_all and not args.enable_scan:
        parser.error("At least one of --enable-all or --enable-scan must be provided.")

    # Handle required targets arguments
    if not args.file and not args.targets and not args.envvars:
        parser.error("At least one of --file or --target must be provided.")

    # Handle scan types
    if args.enable_all:
        for scan_type in config.scan_types_enabled:
            config.scan_types_enabled[scan_type] = True
            if scan_type == "accessible-rdp":
                config.scan_types_enabled[scan_type] = False

    if args.enable_scan:
        for scan_type in args.enable_scan:
            config.scan_types_enabled[scan_type] = True

    # Handle target/s
    if args.file:
        with open(args.file) as f:
            target_hosts: list[ipaddress.IPv4Address] = [
                host for cidr in f for host in ipaddress.ip_network(cidr.strip()).hosts()
            ]
    elif args.targets:
        target_hosts: list[ipaddress.IPv4Address] = [
            host for cidr in args.targets.split(",") for host in ipaddress.ip_network(cidr).hosts()
        ]
    elif args.envvars:
        target_hosts: list[ipaddress.IPv4Address] = [
            host for cidr in os.getenv("TARGET_CIDRS").split(",") for host in ipaddress.ip_network(cidr).hosts()
        ]

    run(target_hosts)


if __name__ == "__main__":
    main()
