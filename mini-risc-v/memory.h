#ifndef __RISC_V_MEM__
#define __RISC_V_MEM__

#include <stdbool.h>
#include <inttypes.h>

bool init(void);
void cleanup(void);

bool readByte(int64_t addr, int64_t *x);
bool readHalf(int64_t addr, int64_t *x);
bool readWord(int64_t addr, int64_t *x);
bool readDouble(int64_t addr, int64_t *x);

bool writeByte(int64_t addr, int8_t x);
bool writeHalf(int64_t addr, int16_t x);
bool writeWord(int64_t addr, int32_t x);
bool writeDouble(int64_t addr, int64_t x);


#endif
