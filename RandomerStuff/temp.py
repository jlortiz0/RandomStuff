#!/usr/bin/python3

import io
import os

def encode(data):
    rand = bytes(os.urandom(64))
    primeInd = 0
    for i in rand:
        primeInd ^= i
        for _ in range(8):
            if primeInd & 1:
                primeInd >>= 1
                primeInd ^= 12
            else:
                primeInd >>= 1

    prime = (2, 3, 5, 7, 11, 13, 17, 19)[primeInd & 7]
    codeTable = list(range(256))
    ind2 = 0
    for ind1 in range(256):
        ind2 += codeTable[ind1] + rand[ind1 & 63]
        ind2 &= 255
        temp = codeTable[ind1]
        codeTable[ind1] = codeTable[ind2]
        codeTable[ind2] = temp

    ind1 = 0
    ind2 = 0
    xorer = 0
    xorind = 0
    output = io.BytesIO(rand)
    output.seek(64)
    for b in data:
        ind1 += prime
        ind1 &= 255
        ind2 = xorind + codeTable[(ind2 + codeTable[ind1]) & 255]
        ind2 &= 255
        xorind += ind1 + codeTable[ind1]
        xorind &= 255
        temp = codeTable[ind1]
        codeTable[ind1] = codeTable[ind2]
        codeTable[ind2] = temp
        xorer = codeTable[(ind2 + codeTable[(ind1 + codeTable[(xorer + xorind) & 255]) & 255]) & 255]
        output.write(bytes((b ^ xorer,)))

    return output.getvalue()

def decode(data):
    primeInd = 0
    for i in range(64):
        primeInd ^= data[i]
        for _ in range(8):
            if primeInd & 1:
                primeInd >>= 1
                primeInd ^= 12
            else:
                primeInd >>= 1

    prime = (2, 3, 5, 7, 11, 13, 17, 19)[primeInd & 7]
    codeTable = list(range(256))
    ind2 = 0
    for ind1 in range(256):
        ind2 += codeTable[ind1] + data[ind1 & 63]
        ind2 &= 255
        temp = codeTable[ind1]
        codeTable[ind1] = codeTable[ind2]
        codeTable[ind2] = temp

    ind1 = 0
    ind2 = 0
    xorer = 0
    xorind = 0
    output = io.BytesIO()
    for b in memoryview(data)[64:]:
        ind1 += prime
        ind1 &= 255
        ind2 = xorind + codeTable[(ind2 + codeTable[ind1]) & 255]
        ind2 &= 255
        xorind += ind1 + codeTable[ind1]
        xorind &= 255
        temp = codeTable[ind1]
        codeTable[ind1] = codeTable[ind2]
        codeTable[ind2] = temp
        xorer = codeTable[(ind2 + codeTable[(ind1 + codeTable[(xorer + xorind) & 255]) & 255]) & 255]
        output.write(bytes((b ^ xorer,)))

    return output.getvalue()
