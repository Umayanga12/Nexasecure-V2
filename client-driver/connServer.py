import sys
import os 
import threading
from command import ESP32Device, RemoveAuthNFT, RemoveReqNFT, Signreqnft, authenticate_device,GetAuthAddr,GetReqAddr, getauthnft, getreqnft, logout_device, setauthnft, setreqnft, signauthnft
from gui import prompt_user_password
from logger import setup_logger
from util import ConnectSocketServer, monitor_wallet_status, wait_for_device


serial_lock = threading.Lock()

logger = setup_logger(log_file="logs/connServer.log", log_level="DEBUG")

def handle_server_commands(client_socket, esp_device):
    command_function_map = {
        "logout": lambda: logout_device(esp_device),
        "getreqnft": lambda: getreqnft(esp_device),
        "getauthnft": lambda: getauthnft(esp_device),
        "setauthnft": lambda data: setauthnft(esp_device, data),
        "setreqnft": lambda data: setreqnft(esp_device, data), 
        "signauthmsg": lambda data: signauthnft(esp_device, data),
        "signreqmsg": lambda data: Signreqnft(esp_device, data),
        "removereqnft": lambda: RemoveReqNFT(esp_device),
        "removeauthnft": lambda: RemoveAuthNFT(esp_device),
        "getauthaddr": lambda: GetAuthAddr(esp_device),
        "getreqaddr": lambda: GetReqAddr(esp_device),
    }

    try:
        while True:
            # Receive command from the socket server
            message = client_socket.recv(1024).decode('utf-8').strip()  # Adjust buffer size as needed
            logger.info(f"Received message from server: {message}")

            # Split the message into command and optional data
            parts = message.split(" ", 1)
            command = parts[0]
            data = parts[1] if len(parts) > 1 else None

            # Validate the command
            if command not in command_function_map:
                logger.warning(f"Invalid command received: {command}")
                client_socket.sendall(f"ERROR: Invalid command '{command}'\n".encode('utf-8'))
                continue

            try:
                # Execute the corresponding function and get the result
                if ("setauthnft" in command or "setreqnft" in command or "signauthmsg" in command or "signreqmsg" in command) and data:
                    result = command_function_map[command](data)
                else:
                    result = command_function_map[command]()
                
                logger.info(f"Command '{command}' executed. Result: {result}")

                # Send the result back to the socket server with a newline character
                client_socket.sendall((result + "\n").encode('utf-8') if isinstance(result, str) else (str(result) + "\n").encode('utf-8'))
                logger.info(f"Sent result back to server: {result}")
            except Exception as e:
                logger.error(f"Error executing command '{command}': {e}")
                client_socket.sendall(f"ERROR: Failed to execute command '{command}'\n".encode('utf-8'))

    except Exception as e:
        logger.error(f"Error while handling server commands: {e}")
    finally:
        client_socket.close()
        logger.info("Socket connection closed.")


# ----- Main Application Flow -----
def main():
    # Wait for the SecureWallet (ESP32 device) to be connected
    port = wait_for_device()
    if not port:
        logger.error("Device not found.")
        return

    esp_device = ESP32Device(port)
    remote_server = "127.0.0.1:10080" 
    client_socket = ConnectSocketServer(remote_server=remote_server)
    #print(client_socket)
    # Prompt user for the password using a GUI
    password = prompt_user_password()
    if not password:
        logger.info("No password entered. Exiting.")
        esp_device.close()
        sys.exit(0)

    if authenticate_device(esp_device, password):
        logger.info("Authentication successful.")

        # Start wallet monitor
        #threading.Thread(target=monitor_wallet_status, args=(port, esp_device, client_socket), daemon=True).start()

        # Get public addresses
        requestPubAddr = GetReqAddr(esp_device)
        authpubAddr = GetAuthAddr(esp_device)

        if client_socket:
            try:
                # Send addresses to server
                message = f"{requestPubAddr},{authpubAddr}\n"
                client_socket.sendall(message.encode('utf-8'))
                logger.info(f"Sent to server: {message}")

                # Wait for server validation
                while True:
                    response = client_socket.recv(1024).decode('utf-8')
                    if response.strip() == "VALIDATED":
                        logger.info("Addresses validated by server. Starting communication...")
                        handle_server_commands(client_socket, esp_device)
                        break
                    elif response.strip() == "RETRY":
                        logger.warning("Server requested retry. Waiting for further instructions...")
                    else:
                        logger.warning(f"Server response: {response}. Authentication failed.")
                        break
            except Exception as e:
                logger.error(f"Communication error: {e}")
            finally:
                client_socket.close()
        else:
            logger.error("Could not connect to the socket server.")
    else:
        logger.warning("Authentication failed.")

    esp_device.close()
    logger.info("ESP32 device connection closed.")

    
if __name__ == "__main__":
    main()
