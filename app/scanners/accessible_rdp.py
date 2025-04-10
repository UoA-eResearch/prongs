from concurrent.futures import ThreadPoolExecutor
from datetime import date
import ipaddress
from queue import Queue
import socket
import threading
import time

from config import Config


PROGRESS_COUNTER = 0


def check_rdp_exposed(ip: str, port: int, result_queue: Queue) -> bool:
    global PROGRESS_COUNTER

    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(2)
        sock.connect((ip, port))
        result_queue.put((ip, port, True))
        sock.close()
        PROGRESS_COUNTER += 1
        return
    except socket.error:
        result_queue.put((ip, port, False))
        PROGRESS_COUNTER += 1
        return


def run(hosts: list[ipaddress.IPv4Network]) -> None:
    result_queue = Queue()

    total_hosts = len(hosts)
    max_workers = 100
    progress_lock = threading.Lock()

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(check_rdp_exposed, str(ip), 3389, result_queue) for ip in hosts]

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
                print(f"🚨 {ip}:{port} is open")
            else:
                print(f"{date.today()}\t{ip}\taccessible-rdp\t3389")

    if Config.pretty_print:
        print(f"Total hosts/checks: {total_hosts}/{PROGRESS_COUNTER}")


def main():
    Config.pretty_print = True
    run(Config.test_host)


if __name__ == "__main__":
    main()
