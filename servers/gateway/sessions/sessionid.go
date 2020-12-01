package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

//InvalidSessionID represents an empty, invalid session ID
const InvalidSessionID SessionID = ""

//idLength is the length of the ID portion
const idLength = 32

//signedLength is the full length of the signed session ID
//(ID portion plus signature)
const signedLength = idLength + sha256.Size

//SessionID represents a valid, digitally-signed session ID.
//This is a base64 URL encoded string created from a byte slice
//where the first `idLength` bytes are crytographically random
//bytes representing the unique session ID, and the remaining bytes
//are an HMAC hash of those ID bytes (i.e., a digital signature).
//The byte slice layout is like so:
//+-----------------------------------------------------+
//|...32 crypto random bytes...|HMAC hash of those bytes|
//+-----------------------------------------------------+
type SessionID string

//ErrInvalidID is returned when an invalid session id is passed to ValidateID()
var ErrInvalidID = errors.New("Invalid Session ID")

//NewSessionID creates and returns a new digitally-signed session ID,
//using `signingKey` as the HMAC signing key. An error is returned only
//if there was an error generating random bytes for the session ID
func NewSessionID(signingKey string) (SessionID, error) {
	// if `signingKey` is zero-length, return InvalidSessionID
	// and an error indicating that it may not be empty
	if len(signingKey) == 0 {
		return InvalidSessionID, errors.New("Signing key may not be empty")
	}
	// Generate a new digitally-signed SessionID that follows the spec of SessionID
	randomBytes := make([]byte, idLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return InvalidSessionID, err
	}
	// create a new HMAC hasher
	remaining := hmac.New(sha256.New, []byte(signingKey))
	// write the SessionID into the hasher
	remaining.Write(randomBytes)
	// caculate the HMAC signature
	remainingBytes := remaining.Sum(nil)
	// https://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go
	// append HMAC signature to the SessionID
	finalByteSlice := append(randomBytes, remainingBytes...)
	// Encode the byteslice to Base64 URL Encoded string
	finalSessionID := SessionID(base64.URLEncoding.EncodeToString(finalByteSlice))

	return finalSessionID, nil
}

//ValidateID validates the string in the `id` parameter
//using the `signingKey` as the HMAC signing key
//and returns an error if invalid, or a SessionID if valid
func ValidateID(id string, signingKey string) (SessionID, error) {
	// decode the id parameter
	decodedID, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		return InvalidSessionID, err
	}

	// separate the SessionID (first 32) and the hmac hash (remaining characters)
	idPortion := decodedID[0:idLength]
	compare := decodedID[idLength:]

	// create new hmac hasher using signing key and sha256 has algorithm
	// hash the SessionID
	remaining := hmac.New(sha256.New, []byte(signingKey))
	_, writeErr := remaining.Write(idPortion)
	if writeErr != nil {
		return InvalidSessionID, writeErr
	}
	// calculate the final hash
	remainingBytes := remaining.Sum(nil)
	// compare the final hash (from hmac hasher) and the digital signature
	if hmac.Equal(compare, remainingBytes) {
		// return sessionID if signature is correct
		return SessionID(id), nil
	}

	// return in hashes do not match
	return InvalidSessionID, ErrInvalidID
}

//String returns a string representation of the sessionID
func (sid SessionID) String() string {
	return string(sid)
}
