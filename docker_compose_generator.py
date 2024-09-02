import sys

def write_initial_config(file):
    file.write("name: tp0\n")
    file.write("services:\n")


def write_server_config(file):
    file.write("  server:\n")
    file.write("    container_name: server\n")
    file.write("    image: server:latest\n")
    file.write("    volumes:\n")
    file.write("      - ./server/config.ini:/config.ini\n")
    file.write("    entrypoint: python3 /main.py\n")
    file.write("    environment:\n")
    file.write("      - PYTHONUNBUFFERED=1\n")
    file.write("    networks:\n")
    file.write("      - testing_net\n\n")


def write_client_config(file, client_number):
    file.write(f"  client{client_number}:\n")
    file.write(f"    container_name: client{client_number}\n")
    file.write("    image: client:latest\n")
    file.write("    volumes:\n")
    file.write("      - ./client/config.yaml:/config.yaml\n")
    file.write("    entrypoint: /client\n")
    file.write("    environment:\n")
    file.write(f"      - CLI_ID={client_number}\n")
    file.write("    networks:\n")
    file.write("      - testing_net\n")
    file.write("    depends_on:\n")
    file.write("      - server\n\n")


def write_network_config(file):
    file.write("networks:\n")
    file.write("  testing_net:\n")
    file.write("    ipam:\n")
    file.write("      driver: default\n")
    file.write("      config:\n")
    file.write("        - subnet: 172.25.125.0/24\n")


def generate_docker_compose(output_filename, amount_clients):
    with open(output_filename, "w") as file:
        write_initial_config(file)
        
        write_server_config(file)

        for i in range(1, amount_clients + 1):
            write_client_config(file, i)

        write_network_config(file)


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(
            "The number of parameters is incorrect. The correct format is: \
            ./generar-compose.sh <output_filename> <amount_clients>"
        )
        sys.exit(1)

    output_filename = sys.argv[1]
    amount_clients = int(sys.argv[2])
    generate_docker_compose(output_filename, amount_clients)
