#include "memory.h"
#include <stdlib.h>
#include <string.h>

int8_t *mem;

bool init(void) {
    mem = calloc(sizeof(int8_t), 0x7fffffff);
    return mem != NULL;
}

void cleanup(void) {
    free(mem);
}

bool readByte(int64_t addr, int64_t *x) {
    if (mem == NULL || addr < 0x1000000 || addr >= 0x80000000) {
		return false;
	}
	*x = (int64_t) mem[addr-0x10000000];
    return true;
}

bool readHalf(int64_t addr, int64_t *x) {
    if (mem == NULL || addr < 0x1000000 || addr >= 0x80000000) {
		return false;
	}
	*x = (int64_t) ((int16_t *) mem)[addr-0x10000000];
    return true;
}

bool readWord(int64_t addr, int64_t *x) {
    if (mem == NULL || addr < 0x1000000 || addr >= 0x80000000) {
		return false;
	}
	*x = (int64_t) ((int32_t *) mem)[addr-0x10000000];
    return true;
}

bool readDouble(int64_t addr, int64_t *x) {
    if (mem == NULL || addr < 0x1000000 || addr >= 0x80000000) {
		return false;
	}
	*x = ((int64_t *) mem)[addr-0x10000000];
    return true;
}

bool writeByte(int64_t addr, int8_t x) {
    if (mem == NULL || addr < 0x1000000 || addr >= 0x80000000) {
		return false;
	}
	mem[addr-0x10000000] = x;
    return true;
}

bool writeHalf(int64_t addr, int16_t x) {
    if (mem == NULL || addr < 0x1000000 || addr >= 0x80000000) {
		return false;
	}
	((int16_t *) mem)[addr-0x10000000] = x;
    return true;
}

bool writeWord(int64_t addr, int32_t x) {
    if (mem == NULL || addr < 0x1000000 || addr >= 0x80000000) {
		return false;
	}
	((int32_t *) mem)[addr-0x10000000] = x;
    return true;
}

bool writeDouble(int64_t addr, int64_t x) {
    if (mem == NULL || addr < 0x1000000 || addr >= 0x80000000) {
		return false;
	}
	((int64_t *) mem)[addr-0x10000000] = x;
    return true;
}
