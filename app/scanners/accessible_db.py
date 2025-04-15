from concurrent.futures import ThreadPoolExecutor
from datetime import datetime
import ipaddress
from queue import Queue
import socket
import threading
import time

from config import Config


PROGRESS_COUNTER = 0


def check_db_exposed(ip: str, port: int, result_queue: Queue) -> bool:
    global PROGRESS_COUNTER

    db_ports = {
        3306,  # mysql
        5432,  # postgres
    }

    for port in db_ports:
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(2)
            sock.connect((ip, port))
            result_queue.put((ip, port, True))
            sock.close()
            PROGRESS_COUNTER += 1
            return
        except socket.error:
            continue

    result_queue.put((ip, port, False))
    PROGRESS_COUNTER += 1
    return


def run(hosts: list[ipaddress.IPv4Network]) -> None:
    result_queue = Queue()

    total_hosts = len(hosts)
    max_workers = 100
    progress_lock = threading.Lock()

    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(check_db_exposed, str(ip), None, result_queue) for ip in hosts]

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
                print(f"{datetime.now().isoformat()}\t{ip}\taccessible-db\t{port}")

    if Config.pretty_print:
        print(f"Total hosts/checks: {total_hosts}/{PROGRESS_COUNTER}")


def main():
    Config.pretty_print = True
    run(Config.test_host)


if __name__ == "__main__":
    main()
