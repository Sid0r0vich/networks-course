from __future__ import annotations

import socket
import sys
import time
from collections.abc import Generator
from contextlib import ExitStack


def traceroute(
    dest_addr: str, max_hops: int = 64, timeout: float = 5, packets_per_hop: int = 3
) -> Generator[tuple[str, list[float]], None, None]:
    with ExitStack() as stack:
        rx = stack.enter_context(
            socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_ICMP)
        )
        tx = stack.enter_context(
            socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
        )
        rx.settimeout(timeout)
        rx.bind(("", 0))

        for ttl in range(1, max_hops + 1):
            tx.setsockopt(socket.IPPROTO_IP, socket.IP_TTL, ttl)
            rtt_times = []

            for _ in range(packets_per_hop):
                tx.sendto(b"", (dest_addr, 33434))

                try:
                    start_time = time.perf_counter_ns()
                    _, curr_addr = rx.recvfrom(512)
                    curr_addr = curr_addr[0]
                    end_time = time.perf_counter_ns()
                    elapsed_time = (end_time - start_time) / 1e6
                    rtt_times.append(elapsed_time)
                except socket.error:
                    curr_addr = None
                    rtt_times.append(None)

            yield curr_addr, rtt_times

            if curr_addr == dest_addr:
                break


def main() -> None:
    if len(sys.argv) < 2:
        print("Usage: python traceroute.py <destination> [packets_per_hop]")
        return

    dest_name = sys.argv[1]
    packets_per_hop = int(sys.argv[2]) if len(sys.argv) > 2 else 3
    dest_addr = socket.gethostbyname(dest_name)

    print(f"Traceroute to {dest_name} ({dest_addr})")
    print(f"{'Hop':<5s}{'IP Address':<20s}{'Hostname':<50s}{'Time (ms)':<10s}")
    print("-" * 90)

    for i, (addr, rtt_times) in enumerate(traceroute(dest_addr, packets_per_hop=packets_per_hop)):
        if addr is not None:
            try:
                host = socket.gethostbyaddr(addr)[0]
            except socket.error:
                host = ""
            
            avg_rtt = sum(filter(None, rtt_times)) / len(rtt_times) if any(rtt_times) else '*'
            print(f"{i+1:<5d}{addr:<20s}{host:<50s}{avg_rtt:<10.3f} ms")
        else:
            print(f"{i+1:<5d}{'*':<20s}{'*':<50s}{'*':<10s}")


if __name__ == "__main__":
    main()
