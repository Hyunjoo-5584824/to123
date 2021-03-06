package swarm

import (
	"time"

	addrutil "github.com/libp2p/go-addr-util"
	ma "github.com/multiformats/go-multiaddr"
)

// ListenAddresses returns a list of addresses at which this swarm listens.
func (s *Swarm) ListenAddresses() []ma.Multiaddr {
	s.listeners.RLock()
	defer s.listeners.RUnlock()
	return s.listenAddressesNoLock()
}

func (s *Swarm) listenAddressesNoLock() []ma.Multiaddr {
	addrs := make([]ma.Multiaddr, 0, len(s.listeners.m))
	for l := range s.listeners.m {
		addrs = append(addrs, l.Multiaddr())
	}
	return addrs
}

const ifaceAddrsCacheDuration = 1 * time.Minute

// InterfaceListenAddresses returns a list of addresses at which this swarm
// listens. It expands "any interface" addresses (/ip4/0.0.0.0, /ip6/::) to
// use the known local interfaces.
func (s *Swarm) InterfaceListenAddresses() ([]ma.Multiaddr, error) {
	s.listeners.RLock() // RLock start

	ifaceListenAddres := s.listeners.ifaceListenAddres
	isEOL := time.Now().After(s.listeners.cacheEOL)
	s.listeners.RUnlock() // RLock end

	if ifaceListenAddres != nil && !isEOL {
		// Cache is valid
		return ifaceListenAddres, nil
	}

	// Cache is not valid
	// Perfrom double checked locking

	s.listeners.Lock() // Lock start

	ifaceListenAddres = s.listeners.ifaceListenAddres
	isEOL = time.Now().After(s.listeners.cacheEOL)
	if ifaceListenAddres == nil || isEOL {
		// Cache is still invalid

		var err error
		ifaceListenAddres, err = addrutil.ResolveUnspecifiedAddresses(
			s.listenAddressesNoLock(), nil)

		if err != nil {
			s.listeners.Unlock() // Lock early exit
			return nil, err
		}

		s.listeners.ifaceListenAddres = ifaceListenAddres
		s.listeners.cacheEOL = time.Now().Add(ifaceAddrsCacheDuration)
	}

	s.listeners.Unlock() // Lock end

	return ifaceListenAddres, nil
}
