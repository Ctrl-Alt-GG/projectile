package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameServerData_Copy(t *testing.T) {
	type fields struct {
		GameServerStaticData  GameServerStaticData
		GameServerDynamicData GameServerDynamicData
	}
	tests := []struct {
		name string
		gsd  GameServerData
	}{
		{
			name: "empty",
			gsd:  GameServerData{},
		},
		{
			name: "full",
			gsd: GameServerData{
				GameServerStaticData: GameServerStaticData{
					Game:         "hello",
					Name:         "world",
					AgentVersion: "1.2",
					Addresses:    []string{"one", "two"},
					Capabilities: Capabilities{
						PlayerCount: true,
						PlayerNames: true,
						PlayerScore: true,
						PlayerTeam:  true,
					},
				},
				GameServerDynamicData: GameServerDynamicData{
					Info:               "info",
					MaxPlayers:         12,
					OnlinePlayersCount: new(uint32(1)),
					OnlinePlayers: new([]Player{
						{
							Name:  "Jozsi",
							Score: new(int32(12)),
							Team:  new("alpha"),
							Info:  "alive",
						},
					}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.gsd.Copy()
			if !reflect.DeepEqual(got, tt.gsd) {
				t.Errorf("Copy() = %v, want %v", got, tt.gsd)
			}
		})
	}
}

func TestCapabilities_IsValid(t *testing.T) {
	type fields struct {
		PlayerCount bool
		PlayerNames bool
		PlayerScore bool
		PlayerTeam  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid__all_false",
			fields: fields{
				PlayerCount: false,
				PlayerNames: false,
				PlayerScore: false,
				PlayerTeam:  false,
			},
			want: true,
		}, {
			name: "valid__1",
			fields: fields{
				PlayerCount: true,
				PlayerNames: false,
				PlayerScore: false,
				PlayerTeam:  false,
			},
			want: true,
		}, {
			name: "valid__2",
			fields: fields{
				PlayerCount: true,
				PlayerNames: true,
				PlayerScore: false,
				PlayerTeam:  false,
			},
			want: true,
		}, {
			name: "valid__3",
			fields: fields{
				PlayerCount: true,
				PlayerNames: true,
				PlayerScore: true,
				PlayerTeam:  false,
			},
			want: true,
		}, {
			name: "valid__4",
			fields: fields{
				PlayerCount: true,
				PlayerNames: true,
				PlayerScore: true,
				PlayerTeam:  true,
			},
			want: true,
		}, {
			name: "invalid__no_count",
			fields: fields{
				PlayerCount: false,
				PlayerNames: true,
				PlayerScore: true,
				PlayerTeam:  true,
			},
			want: false,
		}, {
			name: "invalid__no_names_but_score",
			fields: fields{
				PlayerCount: true,
				PlayerNames: false,
				PlayerScore: true,
				PlayerTeam:  false,
			},
			want: false,
		}, {
			name: "invalid__no_names_but_team",
			fields: fields{
				PlayerCount: true,
				PlayerNames: false,
				PlayerScore: false,
				PlayerTeam:  true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Capabilities{
				PlayerCount: tt.fields.PlayerCount,
				PlayerNames: tt.fields.PlayerNames,
				PlayerScore: tt.fields.PlayerScore,
				PlayerTeam:  tt.fields.PlayerTeam,
			}
			if got := c.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProtobufConversionMirror(t *testing.T) {
	original := GameServerData{
		Game:         "hello",
		Name:         "world",
		AgentVersion: "1.2",
		Addresses:    []string{"asd", "dasd"},
		Capabilities: Capabilities{
			PlayerCount: true,
			PlayerNames: true,
			PlayerScore: true,
			PlayerTeam:  true,
		},
		Info:               "info",
		MaxPlayers:         12,
		OnlinePlayersCount: new(uint32(1)),
		OnlinePlayers: new([]Player{
			{
				Name:  "Jozsi",
				Score: new(int32(12)),
				Team:  new("alpha"),
				Info:  "alive",
			},
		}),
	}

	middle := original.ToProtobuf()

	converted, ok := GameServerDataFromProtobuf(middle)
	assert.True(t, ok)

	if !reflect.DeepEqual(converted, original) {
		t.Errorf("Copy() = %v, want %v", converted, original)
	}
}
