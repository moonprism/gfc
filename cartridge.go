package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// 读取并保存卡带信息

type GameData struct {
	PRG    []byte
	CHR    []byte
	SRAM   []byte  // 8kb存档信息,带电池的卡带才有
}

// 游戏文件结构
// +----------+----------------+--------------+--------------+
// |  HEADER  |    TRAINER     |     PRG      |     CHR      |
// +----------+----------------+--------------+--------------+
// | 16 bytes | 0 or 512 bytes | n*16k        | n*8k         |
// +----------+----------------+--------------+--------------+
// | ROM Info |                | Program Data | Picture Data |
// +----------+----------------+--------------+--------------+
type Cartridge struct {
	*GameData
	Mapper
	Mirror byte
}

// 游戏文件开头的ROM信息标志, 一共是16字节
// +--------+----------+----------+-------+-------+------+
// |  0-3   |    4     |    5     |   6   |   7   | 8-15 |
// +--------+----------+----------+-------+-------+------+
// | "NES␚" | PRG Size | CHR Size | Flag1 | Flag2 | -    |
// +--------+----------+----------+-------+-------+------+
type CartridgeFileHeader struct {
	Magic    uint32
	NumPRG   byte    // PRG块数量，每块16kb
	NumCHR   byte    // CHR块数量，每块8kb
	Flag1    byte
	Flag2    byte
	_        [8]byte
}

func NewCartridgeByPath(path string) (cart *Cartridge, err error) {
	file, err := os.Open(path)
	if err != nil {
			return
	}
	defer file.Close()

	cart = &Cartridge{}

	header := &CartridgeFileHeader{}

	// NES CPU 架构 6502 为小端序
	if err = binary.Read(file, binary.LittleEndian, header); err != nil {
		return
	}

	fmt.Printf("%v\n", header)

	// parse mapper type from flag
	// mapper编号的低四位存在flag1字节中前四位
	mapperIndex := header.Flag1 >> 4
	// todo 后四位存在flag2字节中前四位

	// todo parse Mirroring
	// Flag1字节最后一位，0: 水平镜像，1：垂直镜像

	gd := &GameData{}

	// todo 如果flag1字节第六位是1的话，还需要先读取512byte的TRAINER区域

	fmt.Printf("Flag1: %v\n", header.Flag1)

	// read PRG data from file
	gd.PRG = make([]byte, int(header.NumPRG)*16384)
	if _, err = io.ReadFull(file, gd.PRG); err != nil {
		return
	}

	// read CHR data from file
	gd.CHR = make([]byte, int(header.NumCHR)*8192)
	if _, err = io.ReadFull(file, gd.CHR); err != nil {
		return
	}

	// 构造Mapper
	switch (mapperIndex) {
	case 0:
		cart.Mapper = NewMapper0(gd)
		break
	}

	return
}
