import tkinter as tk
from tkinter import ttk
import socket
import threading
import time


class UDPClient:
    def __init__(self, root):
        self.root = root
        self.root.title("UDP Client")

        ttk.Label(root, text="IP получателя:").grid(row=0, column=0, sticky='w')
        self.ip_entry = ttk.Entry(root)
        self.ip_entry.grid(row=0, column=1)
        self.ip_entry.insert(0, "127.0.0.1")

        ttk.Label(root, text="Порт получателя:").grid(row=1, column=0, sticky='w')
        self.port_entry = ttk.Entry(root)
        self.port_entry.grid(row=1, column=1)
        self.port_entry.insert(0, "12345")

        ttk.Label(root, text="Количество пакетов:").grid(row=2, column=0, sticky='w')
        self.packet_count_entry = ttk.Entry(root)
        self.packet_count_entry.grid(row=2, column=1)
        self.packet_count_entry.insert(0, "10")

        self.send_button = ttk.Button(root, text="Отправить", command=self.start_sending)
        self.send_button.grid(row=4, column=0, columnspan=2, pady=10)

        self.sent_packets = 0

    def start_sending(self):
        try:
            total_packets = int(self.packet_count_entry.get())
            if total_packets <= 0:
                raise ValueError
        except ValueError:
            print("Введите корректное число пакетов")
            return

        print(f"Запуск отправки на IP: {self.ip_entry.get()}, порт: {self.port_entry.get()}")

        threading.Thread(target=self.send_packets, args=(total_packets,), daemon=True).start()

    def send_packets(self, total_packets):
        sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        ip = self.ip_entry.get()
        port = int(self.port_entry.get())

        for _ in range(total_packets):
            try:
                sock.sendto(b"Test", (ip, port))
                time.sleep(0.01)
            except Exception as e:
                print(f"Ошибка при отправке: {e}")
                break


if __name__ == "__main__":
    root = tk.Tk()
    app = UDPClient(root)
    root.mainloop()