package parachain_info

import sc "github.com/LimeChain/goscale"

// 2000 is considered the start of non-common good parachains.
// See: https://wiki.polkadot.network/docs/learn-parachains-faq#are-there-other-ways-of-acquiring-a-slot-besides-the-candle-auction
var (
	defaultParachainId sc.U32 = 2000
)
