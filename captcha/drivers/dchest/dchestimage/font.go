package dchestimage

type Font struct {
	Width     int
	Height    int
	BlackChar byte
	Data      [][]byte
}

var DefaultFont = &Font{
	Width:     11,
	Height:    18,
	BlackChar: 1,
	Data: [][]byte{
		{ // 0
			0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
			0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0,
			1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1,
			0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
		},
		{ // 1
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0,
			0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0,
			0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
		{ // 2
			0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0,
			0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0,
			0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		},
		{ // 3
			0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0,
			0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
			1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0,
		},
		{ // 4
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 1, 1, 0, 1, 1, 0, 0,
			0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0,
			0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0,
			1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
		},
		{ // 5
			0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0,
		},
		{ // 6
			0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0,
			0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0,
			1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0,
			1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 0,
			1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1,
			0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
		},
		{ // 7
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
			0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0,
			0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0,
		},
		{ // 8
			0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0,
			0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0,
			0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0,
			1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0,
			0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
		},
		{ // 9
			0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0,
			0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0,
			1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1,
			0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1,
			0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1,
			0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0,
			0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0,
			0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0,
			0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0,
		},
	},
}
