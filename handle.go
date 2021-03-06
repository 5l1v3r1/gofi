// Package gofi provides a super simple API for sending and receiving data-link packets over WiFi.
package gofi

import "fmt"

type ChannelWidth int

const (
	ChannelWidthUnspecified = iota
	ChannelWidth20MHz
	ChannelWidth40MHz
)

// NewChannelWidthMegahertz creates a ChannelWidth which represents
// the provided number of megahertz.
// This currently supports 20 and 40 MHz, and nothing else.
func NewChannelWidthMegahertz(mhz int) ChannelWidth {
	switch mhz {
	case 20:
		return ChannelWidth20MHz
	case 40:
		return ChannelWidth40MHz
	default:
		return ChannelWidthUnspecified
	}
}

// Megahertz returns the approximate number of megahertz represented
// by this channel width.
func (w ChannelWidth) Megahertz() int {
	return map[ChannelWidth]int{
		ChannelWidth20MHz: 20,
		ChannelWidth40MHz: 40,
	}[w]
}

// A Channel specifies information about a WiFi channel's frequency range.
type Channel struct {
	Number int
	Width  ChannelWidth
}

// A DataRate represents a data rate as a multiple of 500Kb/s.
type DataRate int

// String returns a human-readable string, measured in Mb/s.
func (d DataRate) String() string {
	mbps := float64(d) / 2.0
	return fmt.Sprintf("%.1f Mb/s", mbps)
}

// A Handle facilitates raw WiFi interactions like packet injection,
// sniffing, and channel hopping.
type Handle interface {
	// SupportedRates returns a list of supported outgoing data rates
	// in ascending order.
	SupportedRates() []DataRate

	// SupportedChannels returns a list of supported WLAN channels.
	SupportedChannels() []Channel

	// Channel gets the WLAN channel to which the device is tuned.
	Channel() Channel

	// SetChannel tunes the device into a given WLAN channel.
	// If the channel width is unspecified, the handle will automatically
	// choose an appropriate one.
	SetChannel(Channel) error

	// Receive reads the next packet from the device.
	// The returned RadioInfo will be nil if the device does not
	// support radio information.
	Receive() (Frame, *RadioInfo, error)

	// Send sends a packet over the device.
	// If the given DataRate is 0, the lowest supported rate is used.
	Send(Frame, DataRate) error

	// Close closes the handle.
	// You should always close a Handle once you are done with it.
	//
	// Close synchronously terminates pending Receive() calls.
	// While there is no strict time limit for how long this should take,
	// Close is guaranteed to terminate Receive() calls eventually even if
	// no data is read.
	Close()
}
