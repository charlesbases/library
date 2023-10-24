package nds

import "math"

type level int

const (
	Level0 level = iota
	Level1
	Level2
	Level3
	Level4
	Level5
	Level6
	Level7
	Level8
	Level9
	Level10
	Level11
	Level12
	Level13
	Level14
	Level15
)

// encodeNDS WGS to NDS
func encodeNDS(f float64) int {
	return int(math.Floor(1 << 31 * f / 180))
}

// decodeNDS NDS to WGS
func decodeNDS(tid int, l level) float64 {
	var bin int = 1 << (l*2 + 1)
	return float64(tid) * 180 / float64(bin)
}

// TileID .
type TileID int

// NewTileID .
func NewTileID(x, y float64, l level) TileID {
	xx, yy := encodeNDS(x), encodeNDS(y)

	var bin = xx >> 31
	for i := 1; i <= int(l); i++ {
		bin = bin<<1 + yy>>(31-i)&1
		bin = bin<<1 + xx>>(31-i)&1
	}
	return TileID(bin ^ (1 << (16 + l)))
}

// Decode 返回 TileID 左下角坐标
func (tid TileID) Decode() (float64, float64) {
	var xx, yy = tid.split()
	l := tid.level()
	return decodeNDS(xx, l), decodeNDS(yy, l)
}

// Matrix3 以 tid 为中心的 3 x 3 个 TileID
func (tid TileID) Matrix3() []TileID {
	return tid.matrix(3)
}

// Matrix5 以 tid 为中心的 5 x 5 个 TileID
func (tid TileID) Matrix5() []TileID {
	return tid.matrix(5)
}

// matrix side must be 2k+1
func (tid TileID) matrix(side int) []TileID {
	var ids = make([]TileID, 0, side*side)
	k := (side - 1) / 2
	xx, yy := tid.split()
	for i := -k; i <= k; i++ {
		for j := -k; j <= k; j++ {
			ids = append(ids, tid.merge(xx+i, yy+j))
		}
	}
	return ids
}

// split .
func (tid TileID) split() (int, int) {
	l := int(tid.level())

	var xx, yy TileID = tid >> (2 * l) & 1, 0
	for i := 2*l - 1; i > 0; {
		yy = yy<<1 + tid>>i&1
		i--
		xx = xx<<1 + tid>>i&1
		i--
	}
	return int(xx), int(yy)
}

// merge .
func (tid TileID) merge(xx, yy int) TileID {
	l := int(tid.level())

	var bin = xx >> l
	for i := l - 1; i > -1; i-- {
		bin = bin<<1 + yy>>i&1
		bin = bin<<1 + xx>>i&1
	}
	return TileID(bin ^ (1 << (16 + tid.level())))
}

// level .
func (tid TileID) level() level {
	t := tid >> 16
	for l := 0; l < 16; l++ {
		if t>>l == 1 {
			return level(l)
		}
	}
	return -1
}
