import math
data = []
with open('data', 'r') as f:
   data = f.readlines()


def cal_fuel_mass(input, divisor, subtractor):
    result = math.floor(input / divisor) - subtractor
    if result < 0:
        return 0
    return result

sum = 0
for i in data:
    sum += cal_fuel_mass(int(i), 3, 2)
print sum
