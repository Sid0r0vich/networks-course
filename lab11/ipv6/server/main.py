import socket

def start_server(host='::1', port=12345):
    with socket.socket(socket.AF_INET6, socket.SOCK_STREAM) as server_socket:
        server_socket.bind((host, port))
        server_socket.listen()

        print(f"Сервер запущен на {host}:{port}. Ожидание подключения...")

        while True:
            client_socket, client_address = server_socket.accept()
            with client_socket:
                print(f"Подключено к {client_address}")
                while True:
                    data = client_socket.recv(1024)
                    if not data:
                        break

                    response = data.decode('utf-8').upper()
                    client_socket.sendall(response.encode('utf-8'))

if __name__ == "__main__":
    start_server()