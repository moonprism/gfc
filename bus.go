package main

type Bus interface {
	Read(uint16) byte
	Write(uint16, byte)
}

type cpuBus struct {
	*RAM
}

func NewCPUBus(ram *RAM) Bus {
	return &cpuBus{ram}
}

func (bus *cpuBus) Read(address uint16) byte {
	return 0
}

func (bus *cpuBus) Write(address uint16, value byte) {

}
