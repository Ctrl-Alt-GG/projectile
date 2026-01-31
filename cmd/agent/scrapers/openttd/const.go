package openttd

// Constants / limits (mirror OpenTTD limits where useful)
// Mostly defined in https://github.com/OpenTTD/OpenTTD/blob/1dd3d655747df41cda76bad54d06e86d6efa35ed/src/network/core/config.h
const (
	NST_GRFID_MD5      = 0
	NST_GRFID_MD5_NAME = 1
	NST_LOOKUP_ID      = 2

	MD5_SIZE = 16

	// conservative limits to protect from malformed packets
	MAX_NAME_LEN        = 256  // server_name
	MAX_REVISION_LEN    = 128  // server_revision
	MAX_GAMESCRIPT_NAME = 9000 // large as game script name may be long in some uses
	MAX_GRF_NAME_LEN    = 512
)
