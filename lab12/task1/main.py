import random


class Router:
    def __init__(self, ip):
        self.ip = ip
        self.routing_table = {self.ip: (self.ip, 0)}
        self.neighbors = []

    def add_neighbor(self, neighbor):
        if neighbor not in self.neighbors:
            self.neighbors.append(neighbor)

    def send_rip_update(self):
        routes = []
        for dest_ip, (next_hop, metric) in self.routing_table.items():
            routes.append((dest_ip, metric))
        return routes

    def receive_rip_update(self, sender_ip, routes):
        updated = False
        for dest_ip, metric in routes:
            new_metric = metric + 1
            if dest_ip == self.ip:
                continue
            current = self.routing_table.get(dest_ip)
            if current is None or new_metric < current[1]:
                self.routing_table[dest_ip] = (sender_ip, new_metric)
                updated = True
            elif current[0] == sender_ip and new_metric != current[1]:
                self.routing_table[dest_ip] = (sender_ip, new_metric)
                updated = True
        return updated


def generate_random_network(num_routers=5):
    routers = []
    ips = set()
    while len(ips) < num_routers:
        ip = ".".join(str(random.randint(1, 254)) for _ in range(4))
        ips.add(ip)
    ips = list(ips)

    for ip in ips:
        routers.append(Router(ip))

    p = 0.3
    for router in routers:
        for potential_neighbor in routers:
            if potential_neighbor != router:
                if random.random() < p:
                    router.add_neighbor(potential_neighbor)
                    potential_neighbor.add_neighbor(router)

    return routers


def simulate_rip(routers, max_iterations=20):
    for _ in range(max_iterations):
        updates_occurred = False

        messages = {}
        for router in routers:
            messages[router.ip] = []

        for router in routers:
            routes_to_send = router.send_rip_update()
            for neighbor in router.neighbors:
                messages.setdefault(neighbor.ip, []).append((router.ip, routes_to_send))

        for router in routers:
            received_messages = messages.get(router.ip, [])
            changed_in_this_router = False
            for sender_ip, routes in received_messages:
                if router.receive_rip_update(sender_ip, routes):
                    changed_in_this_router = True
            if changed_in_this_router:
                updates_occurred = True

        if not updates_occurred:
            break


def print_routing_tables(routers):
    for router in routers:
        print(f"\nFinal state of router {router.ip} table:")
        print(f"{'[Source IP]':<15} {'[Destination IP]':<20} {'[Next Hop]':<15} {'[Metric]':<7}")

        sorted_routes = sorted(router.routing_table.items(), key=lambda x: x[0])

        for dest_ip, (next_hop, metric) in sorted_routes:
            if dest_ip != router.ip:
                print(f"{router.ip:<15} {dest_ip:<20} {next_hop:<15} {metric:<7}")


def print_network(routers):
    for router in routers:
        neighbor_ips = [neighbor.ip for neighbor in router.neighbors]
        print(f"{router.ip} -> [{', '.join(neighbor_ips)}]")


def main():
    num_routers_input = input("Num routers: ").strip()
    num_routers = int(num_routers_input) if num_routers_input.isdigit() else 5
    routers = generate_random_network(num_routers)

    print_network(routers)
    simulate_rip(routers)

    print("\nResult:")
    print_routing_tables(routers)


if __name__ == "__main__":
    main()
