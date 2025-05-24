import tkinter as tk
from tkinter import ttk
import socket
import threading
import time


class TCPServer:
    def __init__(self, root):
        self.root = root
        self.root.title("TCP Server")
        ttk.Label(root, text="IP для прослушки:").grid(row=0, column=0, sticky='w')
        self.ip_entry = ttk.Entry(root)
        self.ip_entry.grid(row=0, column=1)
        self.ip_entry.insert(0, "127.0.0.1")
        ttk.Label(root, text="Порт для прослушки:").grid(row=1, column=0, sticky='w')
        self.port_entry = ttk.Entry(root)
        self.port_entry.grid(row=1, column=1)
        self.port_entry.insert(0, "12345")
        self.get_button = ttk.Button(root, text="Получить", command=self.toggle_listening)
        self.get_button.grid(row=4, column=0, columnspan=2, pady=10)
        ttk.Label(root, text="Получено пакетов:").grid(row=2, column=0, sticky='w')
        self.received_count_label = ttk.Label(root, text="0")
        self.received_count_label.grid(row=2, column=1)
        ttk.Label(root, text="Измеренная скорость (байт/с):").grid(row=3, column=0, sticky='w')
        self.speed_label = ttk.Label(root, text="0")
        self.speed_label.grid(row=3, column=1)

        self.packet_count = 0
        self.total_bytes = 0
        self.start_time = None
        self.listen_thread = None
        self.sock = None
        self.listening = False

    def toggle_listening(self):
        if not self.listening:
            self.start_listening()
            self.get_button.config(text="Стоп")
            self.listening = True
        else:
            self.stop_listening()
            self.get_button.config(text="Получить")
            self.listening = False

    def start_listening(self):
        def reset_ui():
            self.received_count_label.config(text="0")
            self.speed_label.config(text="0")

        self.root.after(0, reset_ui)

        ip = self.ip_entry.get()
        port_str = self.port_entry.get()
        try:
            port = int(port_str)
        except:
            print("Некорректный номер порта")
            return

        sock_type = socket.AF_INET6 if ':' in ip else socket.AF_INET
        self.sock = socket.socket(sock_type, socket.SOCK_STREAM)
        try:
            self.sock.bind((ip, port))
            self.sock.listen(1)
            self.listening = True
            threading.Thread(target=self.accept_connections, daemon=True).start()
        except Exception as e:
            print(f"Ошибка привязки сокета: {e}")

    def stop_listening(self):
        self.listening = False
        if hasattr(self, 'sock') and self.sock:
            try:
                self.sock.close()
            except:
                pass

    def accept_connections(self):
        while getattr(self, 'listening', False):
            try:
                conn, addr = self.sock.accept()
                threading.Thread(target=self.handle_client, args=(conn,), daemon=True).start()
            except:
                break

    def handle_client(self, conn):
        start_time = time.time()
        packet_count_local = 0
        total_bytes_local = 0
        try:
            while getattr(self, 'listening', False):
                data = conn.recv(1024)
                if not data:
                    break
                packet_count_local += 1
                total_bytes_local += len(data)
                elapsed_time = time.time() - start_time
                speed = int(total_bytes_local / elapsed_time) if elapsed_time > 0 else 0

                def update():
                    if hasattr(self, 'sock') and getattr(self, 'listening', False):
                        self.received_count_label.config(text=str(packet_count_local))
                        self.speed_label.config(text=str(speed))

                if hasattr(self.root, 'after'):
                    self.root.after(0, update)
            conn.close()
        except:
            pass


if __name__ == "__main__":
    root = tk.Tk()
    app = TCPServer(root)
    root.mainloop()
