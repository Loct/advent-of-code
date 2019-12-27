import math

#data parsing
data = []
with open('data', 'r') as f:
   data = f.readlines()

data = data[0].split(',')
i = 0
for k in data:
  data[i] = int(k)
  i+=1

org = data[:]


def gen_op(op, a, b, values):
    if values[a] == None:
        return
    if values[b] == None:
        return    
    if op == 1:
        return values[a] + values[b]
    if op == 2:
        return values[a] * values[b]

def brute_force(a, b, data):
    data[1] = a
    data[2] = b
    for i in range(0, len(data), 4):
        if data[i] == 99:
            break
        res = gen_op(data[i], data[i+1], data[i+2], data)
        if res == None:
            return data
        data[data[i+3]] = res
    return data

brute_force(12, 0, data)
print data[0]