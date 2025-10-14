from enum import Enum

# ----------------------------
# ENTITIES
# ----------------------------
class VehicleType(Enum):
    CAR = 1
    BIKE = 2

class Vehicle:
    def __init__(self, vtype: VehicleType, reg: str):
        self.vtype = vtype
        self.reg = reg

class Spot:
    def __init__(self, spot_id: int, vtype: VehicleType):
        self.id = spot_id
        self.vtype = vtype
        self.is_occupied = False
        self.vehicle = None

class Level:
    def __init__(self, level_id: int, car_spots: int, bike_spots: int):
        self.id = level_id
        self.spots = {}  # spot_id -> Spot
        self.free_spots = {VehicleType.CAR: set(), VehicleType.BIKE: set()}  # vehicle_type -> free spot IDs
        self.free_count = {VehicleType.CAR: 0, VehicleType.BIKE: 0}           # vehicle_type -> free count

        # Initialize car spots
        for i in range(1, car_spots + 1):
            spot = Spot(i, VehicleType.CAR)
            self.spots[i] = spot
            self.free_spots[VehicleType.CAR].add(i)
            self.free_count[VehicleType.CAR] += 1

        # Initialize bike spots
        for i in range(car_spots + 1, car_spots + bike_spots + 1):
            spot = Spot(i, VehicleType.BIKE)
            self.spots[i] = spot
            self.free_spots[VehicleType.BIKE].add(i)
            self.free_count[VehicleType.BIKE] += 1

class ParkingLot:
    def __init__(self):
        self.levels = []

    def add_level(self, level: Level):
        self.levels.append(level)

    # ----------------------------
    # CORE METHODS
    # ----------------------------
    def park(self, vehicle: Vehicle):
        for level in self.levels:
            free_set = level.free_spots[vehicle.vtype]
            if not free_set:
                continue

            # Pick any free spot (O(1))
            spot_id = free_set.pop()
            spot = level.spots[spot_id]
            spot.is_occupied = True
            spot.vehicle = vehicle
            level.free_count[vehicle.vtype] -= 1

            return level.id, spot.id

        return None, None  # No available spot

    def unpark(self, level_id: int, spot_id: int):
        level = next((l for l in self.levels if l.id == level_id), None)
        if not level:
            raise ValueError("Invalid level ID")

        spot = level.spots.get(spot_id)
        if not spot:
            raise ValueError("Invalid spot ID")
        if not spot.is_occupied:
            raise ValueError("Spot already empty")

        vehicle_type = spot.vtype
        spot.is_occupied = False
        spot.vehicle = None
        level.free_spots[vehicle_type].add(spot_id)
        level.free_count[vehicle_type] += 1

    def display_availability(self):
        print("----- Parking Availability -----")
        for level in self.levels:
            print(f"Level {level.id}: CAR={level.free_count[VehicleType.CAR]}, BIKE={level.free_count[VehicleType.BIKE]}")
        print("--------------------------------")

# ----------------------------
# DEMO
# ----------------------------
if __name__ == "__main__":
    lot = ParkingLot()
    lot.add_level(Level(1, 2, 2))  # Level 1: 2 CAR, 2 BIKE spots
    lot.add_level(Level(2, 1, 3))  # Level 2: 1 CAR, 3 BIKE spots

    lot.display_availability()

    # Park vehicles
    car = Vehicle(VehicleType.CAR, "CAR123")
    bike = Vehicle(VehicleType.BIKE, "BIKE456")

    l1, s1 = lot.park(car)
    print(f"Parked {car.reg} at Level {l1}, Spot {s1}")

    l2, s2 = lot.park(bike)
    print(f"Parked {bike.reg} at Level {l2}, Spot {s2}")

    lot.display_availability()

    # Unpark vehicle
    lot.unpark(l1, s1)
    print(f"Unparked vehicle from Level {l1}, Spot {s1}")

    lot.display_availability()
