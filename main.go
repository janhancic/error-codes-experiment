// An experiment to see how one might encode 3 different bits of information in one 32 bit integer.
// The idea is that we produce an error code with 3 different pieces of information:
// - service name
// - generic error code (db, disk, redis, ...)
// - service specific error code
// The values are numbers which map to something human readable (via an enum not shown here).
// The code is stored as a uint32, but
// is formatted in HEX when displayed as a string and separated with a dot.
// For example; "F5E.96.5E9" could tell us that Service A encountered a DB error while trying to
// insert some record.
// An alternative would be to represent the error code as a byte slice, which would allow us to
// store more error codes (in case we want to chain them).
package main

import (
	"fmt"
	"strconv"
	"strings"
)

/** Returns an uint32 error code such that:
- the first 12 bits represent the service (4096 maximum values)
- the second 8 bits represent a generic error (256 maximum values)
- the last 12 bits represent a specialized error (4096 maximum values)

The function takes in arguments in the smallest types that can hold these values. They could carry
more, but we use bit-shift operations to shift the value to the "beginning" of the value, so that we
can later easily embed the bits in a single uint32, which represents the final error code.
This function doesn't actually do any checking to see if the values are in bounds. If you pass in
values out of the supported range you will get undefined behavior.
*/
func getErrorCode(serviceName uint16, generalErrorCode byte, subErrorCode uint16) uint32 {
	var errorCode uint32

	// Convert to a uint32 and shift for 20 bits (4 to get the original value to the beggining of
	// the uint16 and 16 to get the uint16 to the beginning of a uint32). Then OR the bits together
	// into their final position in the error code.
	errorCode = errorCode | (uint32(serviceName) << 20)

	// Convert to uint32 so it can be bit OR-ed and shift it left 12 places so the byte value sits
	// in the right place of the final error code.
	errorCode = errorCode | (uint32(generalErrorCode) << 12)

	// These bits are already in the right position so we just add them to the error code for the
	// final result.
	errorCode = errorCode | uint32(subErrorCode)

	return errorCode
}

func decodeErrorCode(errorCode uint32) (uint16, byte, uint16) {
	serviceName := uint16(errorCode >> 20)
	// 1044480 coresponds to a uint32 with bits 12-20 set to 1 and all others set to 0
	generalErrorCode := byte((errorCode & 1044480) >> 12)
	// 4095 corresponds to a uint32 with all but the last 12 bits set to 0
	subErrorCode := uint16(errorCode & 4095)
	return serviceName, generalErrorCode, subErrorCode
}

func errorCodeToStr(errorCode uint32) string {
	serviceName, generalErrorCode, subErrorCode := decodeErrorCode(errorCode)
	return fmt.Sprintf("%03X.%02X.%03X", serviceName, generalErrorCode, subErrorCode)
}

func strErrorCodeToErrorCode(strErrorCode string) uint32 {
	parts := strings.Split(strErrorCode, ".")

	serviceName, _ := strconv.ParseUint(parts[0], 16, 16)
	generalErrorCode, _ := strconv.ParseUint(parts[1], 16, 8)
	subErrorCode, _ := strconv.ParseUint(parts[2], 16, 16)

	return getErrorCode(uint16(serviceName), byte(generalErrorCode), uint16(subErrorCode))
}

func main() {
	// print out some random values
	for serviceNameInput := 0; serviceNameInput <= 4095; serviceNameInput += 7 {
		for subErrorCodeInput := 0; subErrorCodeInput <= 4095; subErrorCodeInput += 89 {
			for generalErrorCodeInput := 0; generalErrorCodeInput <= 255; generalErrorCodeInput += 30 {
				errorCode := getErrorCode(uint16(serviceNameInput), byte(generalErrorCodeInput), uint16(subErrorCodeInput))
				strErrorCode := errorCodeToStr(errorCode)
				fmt.Println(strErrorCode)
			}
		}
	}
}
