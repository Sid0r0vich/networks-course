import socket

def start_client(host='::1', port=12345):
    with socket.socket(socket.AF_INET6, socket.SOCK_STREAM) as client_socket:
        client_socket.connect((host, port))
        
        while True:
            message = input("Введите сообщение для отправки (или 'exit' для выхода): ")
            if message.lower() == 'exit':
                break
            
            client_socket.sendall(message.encode('utf-8'))

            response = client_socket.recv(1024)
            print(f"Ответ от сервера: {response.decode('utf-8')}")

if __name__ == "__main__":
    start_client()