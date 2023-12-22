package models

import (
	"encoding/json"
	"strconv"
	"time"
)

type StringifiedInt uint64

type ListSessionsResponse struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	// docs indicate that keys come back as snake case but they definitely come back camelcase...

	SessionID       string         `json:"sessionId"`              // confirmed
	RoomName        string         `json:"roomName"`               // confirmed
	CreatedAt       time.Time      `json:"createdAt"`              // confirmed
	LastActive      time.Time      `json:"lastActive"`             // confirmed
	BandwidthIn     StringifiedInt `json:"bandwidthIn,omitempty"`  // confirmed
	BandwidthOut    StringifiedInt `json:"bandwidthOut,omitempty"` // confirmed
	NumParticipants int            `json:"numParticipants"`        // confirmed

	// TODO KNOWN BUG: does not parse/store the egress key
	// I just don't have any streams with egress to test against

	// docs indicate there should be a "connection_counts" key but I have yet to actually see one in practice
	// (nor have I seen connectionCounts)
	// TODO: reach out to LiveKit team and ask about this; would be useful to have

	// docs indicate there should be a "num_active_participants" key but I have yet to actually see one in practice
	// (nor have I seen numActiveParticipants)
	// TODO: reach out to LiveKit team and ask about this; would be useful to have
}

type SessionDetailsResponse struct {
	// docs indicate that keys come back as snake case but they definitely come back camelcase...

	RoomID          string        `json:"roomId"`          // confirmed
	RoomName        string        `json:"roomName"`        // confirmed
	StartTime       time.Time     `json:"startTime"`       // confirmed
	EndTime         time.Time     `json:"endTime"`         // confirmed
	Participants    []Participant `json:"participants"`    // confirmed
	NumParticipants int           `json:"numParticipants"` // confirmed

	// docs indicate there should be a "status" key but I have yet to actually see one in practice
	// TODO: reach out to LiveKit team and ask about this; would be useful to have

	// docs indicate there should be a "bandwidth" key but I have yet to actually see one in practice
	// TODO: reach out to LiveKit team and ask about this; would be useful to have
}

type Participant struct {
	ParticipantIdentity string               `json:"participantIdentity"` // confirmed
	ParticipantName     string               `json:"participantName"`     // confirmed
	RoomID              string               `json:"roomId"`              // confirmed
	JoinedAt            time.Time            `json:"joinedAt"`            // confirmed
	LeftAt              time.Time            `json:"leftAt"`              // confirmed
	PublishedSources    PublishedSources     `json:"publishedSources"`    // confirmed
	Sessions            []ParticipantSession `json:"sessions"`            // confirmed

	// docs indicate there should be a "isActive" key but I have yet to actually see one in practice
	// possible it only gets sent if isActive=true? but unlikely that I wouldn't have seen that by now
	// TODO: reach out to LiveKit team and ask about this; would be useful to have
}

type PublishedSources struct {
	CameraTrack      bool `json:"cameraTrack,omitempty"`      // confirmed
	MicrophoneTrack  bool `json:"microphoneTrack,omitempty"`  // confirmed
	ScreenShareTrack bool `json:"screenShareTrack,omitempty"` // confirmed
	ScreenShareAudio bool `json:"screenShareAudio,omitempty"` // confirmed
}

type ParticipantSession struct {
	ParticipantSessionID string    `json:"participantId"` // confirmed
	JoinedAt             time.Time `json:"joinedAt"`      // confirmed
	LeftAt               time.Time `json:"leftAt"`        // confirmed
}

// JSON API returns bandwidthIn/bandwidthOut as strings
func (si *StringifiedInt) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}

	*si = StringifiedInt(i)
	return nil
}
