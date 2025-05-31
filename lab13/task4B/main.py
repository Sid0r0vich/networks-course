import asyncio
import threading
import ipaddress
import socket
import tkinter as tk
from tkinter import ttk, scrolledtext

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


async def get_network_generator_async(ip_list, progress_callback):
    tasks = [scan_ip(ip) for ip in ip_list]
    total_ips = len(ip_list)
    responses_list = await asyncio.gather(*tasks)

    cnt = 0
    for responses in responses_list:
        cnt += 1
        progress_callback(cnt, total_ips)
        for psrc, hwsrc in responses:
            yield psrc, hwsrc


def start_scan():
    threading.Thread(target=run_scan_async).start()


def run_scan_async():
    asyncio.run(scan_network())


async def scan_network():
    def progress_callback(current, total):
        progress_var.set(current / total * 100)
        progress_bar.update_idletasks()

    result_text.delete(1.0, tk.END)

    local_ip = get_local_ip()
    ip_mask = f'{".".join(local_ip.split(".")[:-1])}.1/24'

    network = ipaddress.IPv4Network(ip_mask, strict=False)
    ip_list = [str(ip) for ip in network.hosts() if str(ip) != local_ip]
    total_ips = len(ip_list)
    progress_callback(0, total_ips)

    progress_callback(1, total_ips)
    result_text.insert(tk.END, f"{'IP':<15} {'MAC':<17}\n")
    my_responses = await scan_ip(local_ip)
    if my_responses:
        my_ip, my_mac = my_responses[0]
        result_text.insert(tk.END, f"{my_ip:<15} {my_mac:<17} (my device)\n")

    async for ip, mac in get_network_generator_async(ip_list, progress_callback):
        if ip == local_ip:
            continue
        result_text.insert(tk.END, f"{ip:<15} {mac:<17}\n")
        result_text.see(tk.END)
    result_text.insert(tk.END, f"Scanning completed successfully!\n")


if __name__ == "__main__":
    root = tk.Tk()
    root.title("Network Scanner")

    start_button = ttk.Button(root, text="Start scanning", command=start_scan)
    start_button.pack(pady=10)

    progress_var = tk.DoubleVar()
    progress_bar = ttk.Progressbar(root, variable=progress_var, maximum=100)
    progress_bar.pack(fill=tk.X, padx=10)

    result_text = scrolledtext.ScrolledText(root, width=60, height=20)
    result_text.pack(padx=10, pady=10)

    root.mainloop()