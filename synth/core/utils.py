import logging
import datetime
import random
import string
import ssl
import netifaces
import os.path
import defusedxml.ElementTree as ET
import synth
from base64 import b64decode
from termcolor import colored
from functools import wraps
from cryptography import x509
from cryptography.hazmat.backends import default_backend
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives.asymmetric import rsa#, ec
from cryptography.hazmat.primitives import serialization, hashes

class CmdError(Exception):
    def __init__(self, message):
        logging.error(message)
        super().__init__(message)

def get_path_in_package(path):
    return os.path.join(os.path.dirname(synth.__file__), path.lstrip('/'))

def get_data_folder():
    return os.path.expanduser("~/.synth")

def get_path_in_data_folder(path):
    return os.path.join(get_data_folder(), path.lstrip('/'))

def shellcode_to_int_byte_array(data):
    return  ','.join(list(map(str, map(int, data))))

def shellcode_to_hex_byte_array(shellcode):
    byte_array = []
    shellcode_hex = shellcode.hex()
    for i in range(0, len(shellcode_hex), 2):
        byte = shellcode_hex[i:i+2]
        byte_array.append(f"0x{byte.upper()}")

    return ','.join(byte_array)

def shellcode_to_hex_string(shellcode):
    byte_array = []
    shellcode_hex = shellcode.hex()
    for i in range(0, len(shellcode_hex), 2):
        byte = shellcode_hex[i:i + 2]
        byte_array.append(f"\\x{byte.upper()}")

    return ''.join(byte_array)

# Stolen and adapted from https://github.com/zerosum0x0/koadic/blob/master/core/plugin.py
def convert_shellcode(shellcode):
    shellcode = shellcode.translate({ord(c): None for c in '\\x'}).rstrip('\n')
    decis = []
    count = 0
    for i in range(0, len(shellcode), 2):
        count += 1
        hexa = shellcode[i:i + 2]
        deci = int(hexa, 16)

        if count % 25 == 0:
            decis.append(" _\\n" + str(deci))
        else:
            decis.append(str(deci))

    return ",".join(decis)

class PastebinPaste:
    def __init__(self, paste_xml):
        paste_xml = "\n".join(paste_xml.strip().split('\r\n')) + "\n</paste>"

        root = ET.fromstring(paste_xml)
        for child in root.getchildren():
            setattr(self, child.tag.split('_')[1], child.text)

def get_interfaces():
    return netifaces.interfaces()

def get_ipaddress(interface=None):
    if interface and (interface in get_interfaces()):
        return netifaces.ifaddresses(interface)[netifaces.AF_INET][0]['addr']
    else:
        for iface in netifaces.interfaces():
            try:
                netif = netifaces.ifaddresses(iface)
                if netif[netifaces.AF_INET][0]['addr'] == '127.0.0.1':
                    continue
                return netif[netifaces.AF_INET][0]['addr']
            except (ValueError, KeyError):
                continue

            return ""

def get_ips():
    ips = []
    for iface in netifaces.interfaces():
      try:
          netif = netifaces.ifaddresses(iface)
          if netif[netifaces.AF_INET][0]['addr'] == '127.0.0.1':
              continue
          ips.append(netif[netifaces.AF_INET][0]['addr'])
      except (ValueError, KeyError):
          continue

    return ips

def decode_auth_header(request_headers):
    auth_header = request_headers['Authorization']
    auth_header = b64decode(auth_header)
    username, password_digest = auth_header.decode().split(':')
    return username, password_digest

def gen_random_string(length: int = 10):
    return ''.join([random.choice(string.ascii_letters + string.digits) for n in range(length)])

def gen_random_string_no_digits(length: int = 10):
    return ''.join([random.choice(string.ascii_letters) for n in range(length)])

# Adapted from https://cryptography.io/en/latest/x509/tutorial/#creating-a-self-signed-certificate
def create_self_signed_cert(key_file: str = get_path_in_data_folder("key.pem"), cert_file: str = get_path_in_data_folder("cert.pem"), chain_file: str = get_path_in_data_folder("chain.pem")):
    logging.info('Creating self-signed certificate')
    """
    key = ec.generate_private_key(
        curve=ec.SECP521R1(),
        backend=default_backend()
    )
    """
    key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=4096,
        backend=default_backend()
    )
    with open(chain_file, "wb") as ch:
        with open(key_file, "wb") as k:
            privkey_bytes = key.private_bytes(
                encoding=serialization.Encoding.PEM,
                format=serialization.PrivateFormat.PKCS8,
                encryption_algorithm=serialization.NoEncryption(),
            )
            ch.write(privkey_bytes)
            k.write(privkey_bytes)

    subject = issuer = x509.Name([
        x509.NameAttribute(NameOID.COUNTRY_NAME, u"BE"),
        x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, u"VB"),
        x509.NameAttribute(NameOID.LOCALITY_NAME, u"Louvain"),
        x509.NameAttribute(NameOID.ORGANIZATION_NAME, u"1337 H4CK3R"),
        x509.NameAttribute(NameOID.COMMON_NAME, u"croq-et-chocolat.be"),
    ])

    # Sign our certificate with our private key
    cert = x509.CertificateBuilder().subject_name(
        subject
    ).issuer_name(
        issuer
    ).public_key(
        key.public_key()
    ).serial_number(
        x509.random_serial_number()
    ).not_valid_before(
        datetime.datetime.utcnow()
    ).not_valid_after(
        # Our certificate will be valid for 9999 days
        datetime.datetime.utcnow() + datetime.timedelta(days=9999)
    ).add_extension(
        x509.SubjectAlternativeName([x509.DNSName(u"localhost")]),
        critical=False,
    ).sign(key, hashes.SHA256(), default_backend())

    with open(chain_file, "ab") as ch:
        with open(cert_file, "wb") as crt:
            pubkey_bytes = cert.public_bytes(serialization.Encoding.PEM)
            ch.write(pubkey_bytes)
            crt.write(pubkey_bytes)

    logging.info(f"Self-signed certificate written to {key_file}, {cert_file} and {chain_file}")

def get_remote_cert_fingerprint(host: str, port: int):
    pem_data = ssl.get_server_certificate((host, port))
    cert = x509.load_pem_x509_certificate(pem_data.encode(), default_backend())
    return cert.fingerprint(hashes.SHA256())

def get_cert_fingerprint(cert_path):
    with open(cert_path) as pem_data:
        cert = x509.load_pem_x509_certificate(pem_data.read().encode(), default_backend())
        return cert.fingerprint(hashes.SHA256())

def print_good(msg):
    print(f"{colored('[+]', 'green')} {msg}")

def print_bad(msg):
    print(f"{colored('[-]', 'red')} {msg}")

def print_info(msg):
    print(f"{colored('[*]', 'blue')} {msg}")

def print_banner(codename, version):
    logo = """

 .oooooo..o oooooo   oooo ooooo      ooo ooooooooooooo ooooo   ooooo   .oooooo.     .oooo.   
d8P'    `Y8  `888.   .8'  `888b.     `8' 8'   888   `8 `888'   `888'  d8P'  `Y8b  .dP""Y88b  
Y88bo.        `888. .8'    8 `88b.    8       888       888     888  888                ]8P' 
 `"Y8888o.     `888.8'     8   `88b.  8       888       888ooooo888  888              .d8P'  
     `"Y88b     `888'      8     `88b.8       888       888     888  888            .dP'     
oo     .d8P      888       8       `888       888       888     888  `88b    ooo  .oP     .o 
8""88888P'      o888o     o8o        `8      o888o     o888o   o888o  `Y8bood8P'  8888888888 
                
    """
    banner = """
    
                        use the cyberskeleton if you want to survive.
                                        0xSH4RKS

    """
    version = f"""
                               Codename : {colored(codename, "yellow")}
                               Version  : {colored(version, "yellow")}
    """

    print(logo)
    print(banner)
    print(version)
