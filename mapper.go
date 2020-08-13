package main

import "log"

type Mapper interface {
	Read(uint16) byte
	Write(uint16, byte)
}

// mapper0
type Mapper0 struct {
	*GameData
	isMirrored bool
}

func NewMapper0(gd *GameData) Mapper {
	return &Mapper0{gd, len(gd.PRG) == 16 * 1024}
}

func (m Mapper0) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return m.CHR[address]
	case address >= 0x8000:
		if m.isMirrored {
			return m.PRG[address & 0xBFFF - 0x8000]
		} else {
			return m.PRG[address - 0x8000]
		}
	case address >= 0x6000:
		return m.SRAM[address - 0x6000]
	default:
		log.Fatal("mapper0 read error at: %v\n", address)
	}
	return 0
}

func (m Mapper0) Write(address uint16, data byte) {
	switch {
	case address < 0x2000:
		m.CHR[address] = data
	case address >= 0x8000:
		if m.isMirrored {
			m.PRG[address & 0xBFFF - 0x8000] = data
		} else {
			m.PRG[address - 0x8000] = data
		}
	case address >= 0x6000:
		m.SRAM[address - 0x6000] = data
	default:
		log.Fatal("mapper0 write error at: %v\n", address)
	}
}
