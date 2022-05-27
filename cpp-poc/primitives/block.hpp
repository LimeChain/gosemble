#include "primitives.hpp"

struct Block
{
	BlockHeader Header;
	Extrinsic Extrinsics[];
};
