import sys
import threading
import serial
import serial.tools.list_ports
import time
from logger import setup_logger  # Import the setup_logger function

# Initialize the logger
logger = setup_logger(log_file="logs/securewalletOpr.log", log_level="DEBUG")

BAUD_RATE = 115200
READ_TIMEOUT = 1

serial_lock = threading.Lock()

class ESP32Device:
    def __init__(self, port):
        self.port = port
        try:
            with serial_lock:
                self.ser = serial.Serial(port, BAUD_RATE, timeout=READ_TIMEOUT)
                logger.info(f"Opened serial connection on {port}")
        except Exception as e:
            logger.error(f"Failed to open serial port {port}: {e}")
            sys.exit(1)

    def send_command(self, command):
        """Send a command string terminated by newline."""
        full_command = command + "\n"
        self.ser.write(full_command.encode())
        logger.info(f"Sent: {command}")

    def read_line(self):
        """Read one line from the serial port."""
        try:
            line = self.ser.readline().decode(errors="replace").strip()
            return line
        except Exception as e:
            logger.error(f"Error reading from serial: {e}")
            return ""

    def close(self):
        self.ser.close()


def authenticate_device(device, password):
    """Send password to the ESP32 and wait for authentication response."""
    device.send_command(f"PASS {password}")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            # response = response.decode('utf-8').strip()
            if response == "PASSWORD_OK":
                return True
            elif response == "FAIL":
                return False
        time.sleep(0.5)
    return False


def logout_device(device):
    """Send logout command to the ESP32."""
    device.send_command("LOGOUT")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            if response == "Logged out":
                logger.info("Logged out successfully.")
                return True
            else:
                logger.error("Logout failed.")
                return False
        time.sleep(0.5)
    return False

def getreqnft(device):
    device.send_command("GET_NFT_REQ")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received NFT request response: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to get NFT request response.")
    return None

def getauthnft(device):
    device.send_command("GET_NFT_AUTH")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received NFT auth response: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to get NFT auth response.")
    return None

def setreqnft(device, nft):
    device.send_command(f"SET_NFT_REQ {nft}")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received NFT request set response: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to set NFT request response.")
    return None

def setauthnft(device,nft):
    device.send_command(f"SET_NFT_AUTH {nft}")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received NFT auth set response: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to set NFT auth response.")
    return None

def signauthnft(device,msg):
    device.send_command(f"SIGN_MSG_AUTH {msg}")
    print(f"SIGN_MSG_AUTH {msg}")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received NFT auth sign response: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to sign NFT auth response.")
    return None

def Signreqnft(device,msg):
    device.send_command(f"SIGN_MSG_REQ {msg}")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received NFT request sign response: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to sign NFT request response.")
    return None

def RemoveReqNFT(device):
    device.send_command("REMOVE_NFT_REQ")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received NFT request remove response: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to remove NFT request response.")
    return None

def RemoveAuthNFT(device):
    device.send_command("REMOVE_NFT_AUTH")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received NFT auth remove response: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to remove NFT auth response.")
    return None


def GetReqAddr(device):
    device.send_command("GET_ADDR_REQ")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received request address: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to get request address.")
    return None

def GetAuthAddr(device):
    device.send_command("GET_ADDR_AUTH")
    start_time = time.time()
    while time.time() - start_time < 5:
        response = device.read_line()
        if response:
            response = response.strip()
            logger.info(f"Received auth address: {response}")
            return response
        time.sleep(0.5)
    logger.error("Failed to get auth address.")
    return None