```md
(module
  (type (;0;) (func (param i64 i64 i32) (result i64)))
  (import "env" "memory" (memory (;0;) 20))
  
  ;; (import "env" "ext_logging_log_version_1" (func $ext_logging_log_version_1 (param $0 i32) (param $1 i64) (param $2 i64)))
  (import "env" "ext_logging_log_version_1" (func $ext_logging_log_version_1))

  (global $__data_end i32 (i32.const 7))
  (export "__data_end" (global $__data_end))

  (global $__heap_base i32 (i32.const 13))
  (export "__heap_base" (global $__heap_base))
  
  (func $Core_version (param $0 i32) (param $1 i32) (result i64)
    i64.const 0
  )
  (export "Core_version" (func $Core_version))
)

---

Type: wasm
Size: 295 B
Imports:
  Functions:
    "env"."ext_logging_log_version_1": [I32] -> []
  Memories:
    "env"."memory": not shared (2 pages..)
  Tables:
  Globals:
Exports:
  Functions:
    "Core_version": [I32, I32] -> [I64]
  Memories:
  Tables:
    "__indirect_function_table": FuncRef (1..1)
  Globals:
    "__data_end": I32 (constant)
    "__heap_base": I32 (constant)
```