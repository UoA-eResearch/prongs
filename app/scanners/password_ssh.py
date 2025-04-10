from concurrent.futures import ThreadPoolExecutor
from datetime import date
import ipaddress
from queue import Queue
import socket
import threading
import time

import paramiko
import paramiko.ssh_exception

from config import Config


PROGRESS_COUNTER = 0


def check_ssh_password_auth(ip: str, port: int, result_queue: Queue) -> bool:
    global PROGRESS_COUNTER

    # First check port is open using socket
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(2)
        sock.connect((ip, port))
    except socket.error:
        result_queue.put((ip, port, False))
        PROGRESS_COUNTER += 1
        return

    # Attempt to connect to SSH on port 22
    # Set a timeout for 2 seconds
    # Catch any SSH exception or socket error and quit
    try:
        transport = paramiko.Transport((ip, port))
        transport.start_client(timeout=2)
    except (paramiko.SSHException, socket.error):
        result_queue.put((ip, port, False))
        PROGRESS_COUNTER += 1
        return

    # Try to authenticate with no authentication at all and a "random" username
    # We do this to get a list of available authentication methods
    # Check the available authentication methods for "password"
    try:
        transport.auth_none("cats_are_mythical")
    except paramiko.BadAuthenticationType as e:
        if "password" in e.allowed_types:
            result_queue.put((ip, port, True))
            PROGRESS_COUNTER += 1
            return
    except paramiko.SSHException:
        result_queue.put((ip, port, False))
        PROGRESS_COUNTER += 1
        return
    finally:
        transport.close()

    PROGRESS_COUNTER += 1
    result_queue.put((ip, port, False))


def run(hosts: list[ipaddress.IPv4Network]) -> None:
    result_queue = Queue()

    total_hosts = len(hosts)
    max_workers = 100
    progress_lock = threading.Lock()

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(check_ssh_password_auth, str(ip), 22, result_queue) for ip in hosts]

        while any(future.running() for future in futures):
            with progress_lock:
                if Config.pretty_print:
                    print(f"\rProgress: {PROGRESS_COUNTER}/{total_hosts}", end="")
            time.sleep(1)

    for future in futures:
        future.result()

    if Config.pretty_print:
        print(f"\rProgress: {PROGRESS_COUNTER}/{total_hosts}")

    while not result_queue.empty():
        ip, port, status = result_queue.get()
        if status:
            if Config.pretty_print:
                print(f"🚨 {ip}:{port} allows password authentication")
            else:
                print(f"{date.today()}\t{ip}\tpassword-ssh\t22")

    if Config.pretty_print:
        print(f"Total hosts/checks: {total_hosts}/{PROGRESS_COUNTER}")


def main():
    Config.pretty_print = True
    run(Config.test_host)


if __name__ == "__main__":
    main()
