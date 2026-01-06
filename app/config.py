import ipaddress

version: str = "0.2.0"
pretty_print: bool = False

scan_types_enabled: dict[str, bool] = {
    "accessible-db": False,
    "accessible-rdp": False,
    "password-ssh": False,
}

test_cidr = ipaddress.ip_network("45.33.32.156/32")  # scanme.nmap.org
test_host: list[ipaddress.IPv4Address] = list(test_cidr.hosts())
