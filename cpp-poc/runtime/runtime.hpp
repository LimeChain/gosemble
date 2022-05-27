#ifndef RUNTIME_H
#define RUNTIME_H
#define WASM_EXPORT __attribute__((visibility("default"))) extern "C"

#include <stdint.h>

#include "env.hpp"
#include "core.hpp"
#include "block_builder.hpp"
#endif // RUNTIME_H