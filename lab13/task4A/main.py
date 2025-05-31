import asyncio
import ipaddress
import socket

import scapy.all as sc


def get_local_ip():
    st = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        st.connect(('10.255.255.255', 1))
        ip_l = st.getsockname()[0]
    except Exception:
        ip_l = '127.0.0.1'
    finally:
        st.close()
    return ip_l


async def scan_ip(ip):
    answered = await asyncio.get_event_loop().run_in_executor(
        None,
        lambda: sc.srp(
            sc.Ether(dst='ff:ff:ff:ff:ff:ff') / sc.ARP(pdst=ip),
            timeout=1,
            verbose=False
        )[0]
    )
    results = []
    for sent, received in answered:
        results.append((received.psrc, received.hwsrc))
    return results


async def get_network_generator_async(ip_list):
    tasks = [scan_ip(ip) for ip in ip_list]
    responses_list = await asyncio.gather(*tasks)

    for responses in responses_list:
        for psrc, hwsrc in responses:
            yield psrc, hwsrc


async def main():
    local_ip = get_local_ip()
    ip_mask = f'{".".join(local_ip.split(".")[:-1])}.1/24'

    network = ipaddress.IPv4Network(ip_mask, strict=False)
    ip_list = [str(ip) for ip in network.hosts() if ip != local_ip]

    print(f"{'IP':<15} {'MAC':<17}")
    print('-' * 32)
    my_ip, my_mac = (await scan_ip(local_ip))[0]
    print(f"{my_ip:<15} {my_mac:<17} (my device)")
    async for ip, mac in get_network_generator_async(ip_list):
        if ip == local_ip:
            continue
        print(f"{ip:<15} {mac:<17}")


if __name__ == "__main__":
    asyncio.run(main())
