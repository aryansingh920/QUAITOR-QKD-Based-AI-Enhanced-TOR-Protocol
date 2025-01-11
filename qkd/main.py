from qiskit import QuantumCircuit, QuantumRegister, ClassicalRegister
from qiskit_aer import AerSimulator
from qiskit.quantum_info import random_statevector
import numpy as np
import random
import argparse
import base64
import json
import os


class QKD:
    def __init__(self):
        self.simulator = AerSimulator()
        self.key_file = 'quantum_key.json'

    def generate_quantum_key(self, key_length):
        """Generate quantum key using BB84 protocol"""
        # Check if we should reuse existing key
        if os.path.exists(self.key_file):
            with open(self.key_file, 'r') as f:
                saved_data = json.load(f)
                return saved_data['key']

        alice_bits = [random.randint(0, 1) for _ in range(key_length)]
        alice_bases = [random.randint(0, 1) for _ in range(
            key_length)]  # 0 for Z, 1 for X
        bob_bases = [random.randint(0, 1) for _ in range(
            key_length)]    # 0 for Z, 1 for X

        qubits = []
        for bit, basis in zip(alice_bits, alice_bases):
            # Create quantum circuit for each qubit
            qr = QuantumRegister(1)
            cr = ClassicalRegister(1)
            qc = QuantumCircuit(qr, cr)

            # Prepare qubit based on Alice's bit and basis
            if bit == 1:
                qc.x(qr)
            if basis == 1:
                qc.h(qr)

            # Bob's measurement
            if bob_bases[len(qubits)] == 1:
                qc.h(qr)
            qc.measure(qr, cr)

            qubits.append(qc)

        # Execute all circuits using the simulator
        bob_measurements = []
        for qc in qubits:
            job = self.simulator.run(qc, shots=1)
            result = job.result()
            counts = result.get_counts(0)
            bob_measurements.append(int(list(counts.keys())[0]))

        # Sift key - keep only bits where Alice and Bob used same basis
        shared_key = ''
        for i in range(key_length):
            if alice_bases[i] == bob_bases[i]:
                shared_key += str(alice_bits[i])

        # Save the key
        with open(self.key_file, 'w') as f:
            json.dump({'key': shared_key}, f)

        return shared_key

    def encrypt(self, message, key):
        """Encrypt message using XOR with quantum-generated key"""
        # Convert message to binary
        binary_message = ''.join(format(ord(c), '08b') for c in message)

        # Extend key if necessary
        extended_key = key * (len(binary_message) // len(key) + 1)
        extended_key = extended_key[:len(binary_message)]

        # XOR message with key
        encrypted_binary = ''.join(str(int(a) ^ int(b))
                                   for a, b in zip(binary_message, extended_key))

        # Convert to base64 for easier transmission
        encrypted_bytes = int(encrypted_binary, 2).to_bytes(
            (len(encrypted_binary) + 7) // 8, byteorder='big')
        encrypted_message = base64.b64encode(encrypted_bytes).decode()

        return encrypted_message

    def decrypt(self, encrypted_message, key):
        """Decrypt message using XOR with quantum-generated key"""
        # Convert from base64
        encrypted_bytes = base64.b64decode(encrypted_message)
        encrypted_binary = bin(int.from_bytes(encrypted_bytes, byteorder='big'))[
            2:].zfill(len(encrypted_bytes) * 8)

        # Extend key if necessary
        extended_key = key * (len(encrypted_binary) // len(key) + 1)
        extended_key = extended_key[:len(encrypted_binary)]

        # XOR encrypted message with key
        decrypted_binary = ''.join(str(int(a) ^ int(b))
                                   for a, b in zip(encrypted_binary, extended_key))

        # Convert binary to text
        decrypted_message = ''
        for i in range(0, len(decrypted_binary), 8):
            byte = decrypted_binary[i:i+8]
            decrypted_message += chr(int(byte, 2))

        return decrypted_message

    def clear_saved_key(self):
        """Clear the saved quantum key"""
        if os.path.exists(self.key_file):
            os.remove(self.key_file)


def main():
    parser = argparse.ArgumentParser(
        description='Quantum Key Distribution Encryption/Decryption')
    parser.add_argument(
        '--mode', choices=['encrypt', 'decrypt', 'clear-key'], required=True, help='Operation mode')
    parser.add_argument('--message', help='Message to encrypt/decrypt')
    parser.add_argument('--key-length', type=int, default=256,
                        help='Length of quantum key to generate')

    args = parser.parse_args()

    qkd = QKD()

    if args.mode == 'clear-key':
        qkd.clear_saved_key()
        print("Saved quantum key cleared.")
        return

    # Generate quantum key
    print("Generating quantum key...")
    quantum_key = qkd.generate_quantum_key(args.key_length)
    print(f"Generated key: {quantum_key}")

    if args.mode == 'encrypt':
        if not args.message:
            print("Error: Message is required for encryption")
            return
        encrypted = qkd.encrypt(args.message, quantum_key)
        print(f"\nEncrypted message: {encrypted}")
    else:
        if not args.message:
            print("Error: Message is required for decryption")
            return
        decrypted = qkd.decrypt(args.message, quantum_key)
        print(f"\nDecrypted message: {decrypted}")


if __name__ == "__main__":
    main()
