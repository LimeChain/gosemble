#include "primitives.hpp"

struct BlockHeader
{
  Blake2bHash ParentHash;
  uint64_t Number;
  Hash StateRoot;
  Hash ExtrinsicsRoot;
  Digest Digest;
};