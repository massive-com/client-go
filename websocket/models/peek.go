package models

import "bytes"

var evNeedle = []byte(`"ev":"`)

// PeekEventType returns the value of the "ev" discriminator field from a raw
// WebSocket message via a byte-scan. It's ~100× faster than json.Unmarshal
// into an EventType struct and is intended for routing frames to the right
// typed decoder on the hot path.
//
// Assumes compact JSON (the WS server emits no whitespace around the colon)
// and that the ev value is a plain ASCII identifier with no escape sequences.
// Both hold for every event type the server emits.
//
// Returns "" if the message is missing the ev field, truncated, or otherwise
// malformed. Callers should treat an empty return as an unknown event and
// skip the message.
func PeekEventType(msg []byte) string {
	i := bytes.Index(msg, evNeedle)
	if i < 0 {
		return ""
	}
	start := i + len(evNeedle)
	if start >= len(msg) {
		return ""
	}
	end := bytes.IndexByte(msg[start:], '"')
	if end < 0 {
		return ""
	}
	return string(msg[start : start+end])
}
