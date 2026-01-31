package openttd

import (
	"fmt"
	"io"

	"github.com/Ctrl-Alt-GG/projectile/pkg/utils"
)

// Based on https://github.com/OpenTTD/OpenTTD/blob/1dd3d655747df41cda76bad54d06e86d6efa35ed/src/network/core/network_game_info.cpp#L271-L377

// ParseNetworkGameInfo parses the serialized NetworkServerGameInfo from data.
// Returns parsed struct or an error on malformed input.
func ParseNetworkGameInfo(r io.Reader) (NetworkServerGameInfo, error) {
	var info NetworkServerGameInfo

	// Start parsing in exact order as SerializeNetworkGameInfo.
	structReader := utils.NewStructReader(r)

	var err error

	// 1) version byte
	info.Version, err = structReader.ReadUint8()
	if err != nil {
		return info, fmt.Errorf("reading version: %w", err)
	}

	// version >= 7: uint64 ticks_playing
	if info.Version >= 7 {
		info.TicksPlaying, err = structReader.ReadUint64()
		if err != nil {
			return info, fmt.Errorf("reading ticks_playing: %w", err)
		}
	}

	// version >= 6: newgrf serialization type (uint8)
	if info.Version >= 6 {
		info.NewGRFSerialization, err = structReader.ReadUint8()
		if err != nil {
			return info, fmt.Errorf("reading newgrf serialization type: %w", err)
		}
	} else {
		// older versions default to NST_GRFID_MD5
		info.NewGRFSerialization = NST_GRFID_MD5
	}

	// version >= 5: gamescript version + gamescript name (string)
	if info.Version >= 5 {
		info.GameScriptVersion, err = structReader.ReadUint32()
		if err != nil {
			return info, fmt.Errorf("reading gamescript version: %w", err)
		}

		// gamescript name is a null-terminated string. Use large max.
		info.GameScriptName, err = structReader.ReadNullTerminatedString(MAX_GAMESCRIPT_NAME)
		if err != nil {
			return info, fmt.Errorf("reading gamescript name: %w", err)
		}
	}

	// version >= 4: GRF list (I have no idea what GRFs are lol)
	if info.Version >= 4 {
		count, err := structReader.ReadUint8()
		if err != nil {
			return info, fmt.Errorf("reading grf count: %w", err)
		}

		info.GRFs = make([]GRFEntry, 0, count)
		for i := uint8(0); i < count; i++ {
			grfid, err := structReader.ReadUint32()
			if err != nil {
				return info, fmt.Errorf("reading grf[%d].grfid: %w", i, err)
			}

			var md5 []byte
			md5, err = structReader.ReadBytes(MD5_SIZE)
			if err != nil {
				return info, fmt.Errorf("reading grf[%d].md5: %w", i, err)
			}

			entry := GRFEntry{
				GRFID: grfid,
				MD5:   md5,
			}

			// read optional name if serialization type says so
			if info.NewGRFSerialization == NST_GRFID_MD5_NAME {
				entry.Name, err = structReader.ReadNullTerminatedString(MAX_GRF_NAME_LEN)
				if err != nil {
					return info, fmt.Errorf("reading grf[%d].name: %w", i, err)
				}
			} else if info.NewGRFSerialization == NST_LOOKUP_ID {
				// Not implemented: a lookup-based representation requires prior lookup table
				return info, fmt.Errorf("NST_LOOKUP_ID parsing not implemented")
			}
			info.GRFs = append(info.GRFs, entry)
		}
	}

	// version >= 3: calendar_date, calendar_start (uint32 each)
	if info.Version >= 3 {
		info.CalendarDate, err = structReader.ReadUint32()
		if err != nil {
			return info, fmt.Errorf("reading calendar_date: %w", err)
		}
		info.CalendarStart, err = structReader.ReadUint32()
		if err != nil {
			return info, fmt.Errorf("reading calendar_start: %w", err)
		}
	}

	// version >= 2: companies_max, companies_on, (old) clients_max
	if info.Version >= 2 {
		info.CompaniesMax, err = structReader.ReadUint8()
		if err != nil {
			return info, fmt.Errorf("reading companies_max: %w", err)
		}
		info.CompaniesOn, err = structReader.ReadUint8()
		if err != nil {
			return info, fmt.Errorf("reading companies_on: %w", err)
		}
		info.ClientsMax, err = structReader.ReadUint8() // will be overwritten by the version 1 clients_max read below
		if err != nil {
			return info, fmt.Errorf("reading old clients_max: %w", err)
		}
	}

	// version >= 1: base fields (string fields, booleans, etc)
	if info.Version >= 1 {
		// server_name
		info.ServerName, err = structReader.ReadNullTerminatedString(MAX_NAME_LEN)
		if err != nil {
			return info, fmt.Errorf("reading server_name: %w", err)
		}

		// server_revision
		info.ServerRev, err = structReader.ReadNullTerminatedString(MAX_REVISION_LEN)
		if err != nil {
			return info, fmt.Errorf("reading server_revision: %w", err)
		}

		// use_password (bool)
		info.UsePassword, err = structReader.ReadUint8Bool()
		if err != nil {
			return info, fmt.Errorf("reading use_password: %w", err)
		}

		// clients_max (again)
		info.ClientsMax, err = structReader.ReadUint8()
		if err != nil {
			return info, fmt.Errorf("reading clients_max: %w", err)
		}

		// clients_on
		info.ClientsOn, err = structReader.ReadUint8()
		if err != nil {
			return info, fmt.Errorf("reading clients_on: %w", err)
		}

		// spectators_on
		info.SpectatorsOn, err = structReader.ReadUint8()
		if err != nil {
			return info, fmt.Errorf("reading spectators_on: %w", err)
		}

		// map_width, map_height (uint16)
		info.MapWidth, err = structReader.ReadUint16()
		if err != nil {
			return info, fmt.Errorf("reading map_width: %w", err)
		}
		info.MapHeight, err = structReader.ReadUint16()
		if err != nil {
			return info, fmt.Errorf("reading map_height: %w", err)
		}

		// landscape (uint8)
		info.Landscape, err = structReader.ReadUint8()
		if err != nil {
			return info, fmt.Errorf("reading landscape: %w", err)
		}

		// dedicated (bool)
		info.Dedicated, err = structReader.ReadUint8Bool()
		if err != nil {
			return info, fmt.Errorf("reading dedicated: %w", err)
		}
	}

	return info, nil
}
