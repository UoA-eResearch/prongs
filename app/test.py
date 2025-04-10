import ipaddress

from config import Config
from scanners import accessible_db
from scanners import accessible_rdp
from scanners import password_ssh


def main():
    Config.pretty_print = False

    # Set target, based on app/config.py
    target_hosts: list[ipaddress.IPv4Address] = Config.test_host

    # Enabled/disable specific checks
    accessible_db.run(target_hosts)
    accessible_rdp.run(target_hosts)
    password_ssh.run(target_hosts)


if __name__ == "__main__":
    main()
