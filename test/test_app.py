import ipaddress
import sys
from pathlib import Path

PROJECT_ROOT = Path(__file__).resolve().parents[1]
APP_DIR = PROJECT_ROOT / "app"
if str(APP_DIR) not in sys.path:
    sys.path.insert(0, str(APP_DIR))

import config
from scanners import accessible_db, accessible_rdp, password_ssh


def test_scanners_smoke(capsys):
    config.pretty_print = False

    target_hosts: list[ipaddress.IPv4Address] = config.test_host

    accessible_db.run(target_hosts)
    accessible_rdp.run(target_hosts)
    password_ssh.run(target_hosts)

    captured = capsys.readouterr()
    stdout_lines = [line for line in captured.out.splitlines() if line.strip()]

    assert captured.err == "", "Scanners should not write to stderr"
    assert stdout_lines, "Expected scanner output for SSH check"
    assert any(
        ("password-ssh" in line) and ("45.33.32.156" in line) and ("\t22" in line)
        for line in stdout_lines
    ), "SSH password auth should be reported open"
    assert all("accessible-db" not in line for line in stdout_lines), ""
    assert all("accessible-rdp" not in line for line in stdout_lines), ""
