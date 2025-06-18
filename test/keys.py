import random
import string

def generateRandomKeys(count, length):
    keys = set()
    
    while len(keys) < count:
        key = ''.join(random.choices(string.ascii_letters + string.digits, k=length))
        keys.add(key)
        
    return list(keys)

def saveKeys(keys, filename):
    with open(filename, 'w') as f:
        for key in keys:
            f.write(key + '\n')

if __name__ == "__main__":
    import sys

    if len(sys.argv) != 4:
        print("Usage: python3 keys.py <count> <length> <outputFile>")
        sys.exit(1)
    
    count      = int(sys.argv[1])
    length     = int(sys.argv[2])
    outputFile = sys.argv[3]
    
    keys = generateRandomKeys(count, length)
    saveKeys(keys, outputFile)
    
    print(f"Generated {count} unique keys of length {length} and saved to {outputFile}")
