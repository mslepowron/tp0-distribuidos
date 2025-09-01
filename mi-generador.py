import sys

def generate_docker_compose(file, clients):
    compose = """name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - CLIENTS={clients}
    networks:
      - testing_net
    volumes:
      - ./server/config.ini:/config.ini
"""

    for client_number in range(1, clients + 1):
        compose += f"""  client{client_number}:
    container_name: client{client_number}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={client_number}
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{client_number}.csv:/agency.csv
    depends_on:
      - server
"""

    compose += """networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""
        
    with open(file, 'w') as f:
        f.write(compose)

    print(f"Docker Compose file '{file}' generado con {clients} clientes.")

if __name__ == "__main__":
    
    file = sys.argv[1]
    clients = int(sys.argv[2])

    generate_docker_compose(file, clients)