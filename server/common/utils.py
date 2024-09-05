import csv
import datetime
import time
import struct


""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574


""" A lottery bet registry. """
class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)

    def __repr__(self):
        return (f"Bet(agency={self.agency}, first_name={self.first_name}, last_name={self.last_name}, "
                f"document={self.document}, birthdate={self.birthdate}, number={self.number}, agency={self.agency})")

    """
    Deserialize a bet from a bytes object.
    """
    @staticmethod
    def deserialize(data):
        bytes_read = 0

        if len(data) < 1:
            return None

        name_length = struct.unpack('B', data[bytes_read:bytes_read+1])[0]
        bytes_read += 1

        if len(data) < bytes_read + name_length:
            return None

        first_name = data[bytes_read:bytes_read+name_length].decode('utf-8')
        bytes_read += name_length

        if len(data) < bytes_read + 1:
            return None

        lastname_length = struct.unpack('B', data[bytes_read:bytes_read+1])[0]
        bytes_read += 1

        if len(data) < bytes_read + lastname_length:
            return None

        last_name = data[bytes_read:bytes_read+lastname_length].decode('utf-8')
        bytes_read += lastname_length

        if len(data) < bytes_read + 4:
            return None

        document = struct.unpack('>I', data[bytes_read:bytes_read+4])[0]
        bytes_read += 4

        if len(data) < bytes_read + 1:
            return None

        birth_day = struct.unpack('B', data[bytes_read:bytes_read+1])[0]
        bytes_read += 1

        if len(data) < bytes_read + 1:
            return None

        birth_month = struct.unpack('B', data[bytes_read:bytes_read+1])[0]
        bytes_read += 1

        if len(data) < bytes_read + 2:
            return None

        birth_year = struct.unpack('>H', data[bytes_read:bytes_read+2])[0]
        bytes_read += 2

        birthdate = f"{birth_year}-{str(birth_month).zfill(2)}-{str(birth_day).zfill(2)}"

        if len(data) < bytes_read + 4:
            return None

        number = struct.unpack('>I', data[bytes_read:bytes_read+4])[0]
        bytes_read += 4

        if len(data) < bytes_read + 1:
            return None

        agency = struct.unpack('B', data[bytes_read:bytes_read+1])[0]

        return Bet(agency, first_name, last_name, document, birthdate, number)

""" Checks whether a bet won the prize or not. """
def has_won(bet: Bet) -> bool:
    return bet.number == LOTTERY_WINNER_NUMBER

"""
Persist the information of each bet in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def store_bets(bets: list[Bet]) -> None:
    with open(STORAGE_FILEPATH, 'a+') as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow([bet.agency, bet.first_name, bet.last_name,
                             bet.document, bet.birthdate, bet.number])

"""
Loads the information all the bets in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def load_bets() -> list[Bet]:
    with open(STORAGE_FILEPATH, 'r') as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])

