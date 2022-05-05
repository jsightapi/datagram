package jse

import (
	"github.com/google/uuid"
	"net"
	"testing"
)

func TestEncode_Encode(t *testing.T) {
	uid, err := uuid.Parse("31c14807-084b-4884-bdd9-f32d3eaddc12")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		event   Event
		want    string
		wantErr bool
	}{
		{
			Event{
				Version:      1,
				AppType:      AppEditor,
				AppHost:      "host",
				ClientID:     uid,
				ClientIPv4:   net.ParseIP("0.0.0.0").To4(),
				ProjectTitle: "aaa",
				ProjectSize:  1000,
				ParseError:   nil,
			},
			"",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.event.ProjectTitle, func(t *testing.T) {
			got := tt.event.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
