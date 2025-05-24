import tkinter as tk
from tkinter import ttk
import socket
import threading
import time


class TCPClient:
    def __init__(self, root):
        self.root = root
        self.root.title("TCP Client")

        ttk.Label(self.root, text="IP получателя:").grid(row=0, column=0, sticky='w')
        self.ip_entry = ttk.Entry(self.root)
        self.ip_entry.grid(row=0, column=1)
        self.ip_entry.insert(0, "127.0.0.1")

        ttk.Label(self.root, text="Порт получателя:").grid(row=1, column=0, sticky='w')
        self.port_entry = ttk.Entry(self.root)
        self.port_entry.grid(row=1, column=1)
        self.port_entry.insert(0, "12345")

        ttk.Label(self.root, text="Количество пакетов:").grid(row=2, column=0, sticky='w')
        self.packet_count_entry = ttk.Entry(self.root)
        self.packet_count_entry.grid(row=2, column=1)
        self.packet_count_entry.insert(0, "10")

        ttk.Button(self.root, text="Отправить", command=self.start_sending).grid(row=4, columnspan=2, pady=10)

    def start_sending(self):
        try:
            total_packets = int(self.packet_count_entry.get())
            if total_packets <= 0:
                raise ValueError
        except:
            print("Введите корректное число пакетов")
            return

        print(f"Запуск отправки на IP:{self.ip_entry.get()}, порт:{self.port_entry.get()}")
        threading.Thread(target=self.send_packets, args=(total_packets,), daemon=True).start()

    def send_packets(self, total_packets):
        ip = self.ip_entry.get()
        port = int(self.port_entry.get())

        print(f"Попытка подключения к {ip}:{port}")
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.connect((ip, port))
            print("Подключение успешно")
        except Exception as e:
            print(f"Ошибка подключения: {e}")
            return

        for i in range(total_packets):
            try:
                sock.sendall(b"Test")
                print(f"Отправлен пакет {i + 1}/{total_packets}")
                time.sleep(0.01)
            except Exception as e:
                print(f"Ошибка при отправке пакета {i + 1}: {e}")
                break

        try:
            sock.close()
            print("Соединение закрыто")
        except:
            pass


if __name__ == "__main__":
    root = tk.Tk()
    app = TCPClient(root)
    root.mainloop()