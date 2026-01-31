package openttd

// NetworkServerGameInfo is a Go representation of the server info.
// Based on https://github.com/OpenTTD/OpenTTD/blob/1dd3d655747df41cda76bad54d06e86d6efa35ed/src/network/core/network_game_info.h#L92-L114
type NetworkServerGameInfo struct {
	// version byte from the wire
	Version uint8

	// version >= 7
	TicksPlaying uint64

	// version >= 6
	NewGRFSerialization uint8

	// version >= 5
	GameScriptVersion uint32
	GameScriptName    string

	// version >= 4 : GRF list
	GRFs []GRFEntry

	// version >= 3
	CalendarDate  uint32
	CalendarStart uint32

	// version >= 2 (and also in v1)
	CompaniesMax uint8
	CompaniesOn  uint8
	// note: clients_max appears multiple times in wire; we store last-seen
	ClientsMax uint8

	// version >= 1 base fields
	ServerName   string
	ServerRev    string
	UsePassword  bool
	ClientsOn    uint8
	SpectatorsOn uint8
	MapWidth     uint16
	MapHeight    uint16
	Landscape    uint8
	Dedicated    bool
}

// GRFEntry represents one GRF identification entry read from wire.
type GRFEntry struct {
	GRFID uint32
	MD5   []byte // MD5_SIZE bytes
	Name  string // optional; present for NST_GRFID_MD5_NAME
}
