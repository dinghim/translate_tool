package tabfile

var (
	cr byte = 0x0d //回车CR
	lf byte = 0x0a //换行LF
	tb byte = 0x09 //tab制表符
)

const (
	state_normal        = iota //正常状态
	state_double_quotes        //"双引号"字符串
)

type tabfile struct{}

func New() *tabfile {
	return &tabfile{}
}

func (t *tabfile) filter(text []byte) bool {
	for i := 0; i < len(text); i++ {
		if text[i]&0x80 != 0 {
			return false
		}
	}
	return true
}

func (t *tabfile) GetString(text []byte) ([][]byte, []int, []int, error) {
	var entryStart []int
	var entryEnd []int
	var entryTotal [][]byte
	frecord := func(nStart, nEnd int) {
		slice := text[nStart:nEnd]
		if !t.filter(slice) {
			entryStart = append(entryStart, nStart)
			entryEnd = append(entryEnd, nEnd)
			entryTotal = append(entryTotal, slice)
		}
	}
	nStart := 0
	length := len(text)
	for i := 0; i < length; i++ {
		if text[i] == tb {
			frecord(nStart, i)
			nStart = i + 1
		} else if i+1 < length && text[i] == cr && text[i] == lf {
			frecord(nStart, i)
			nStart = i + 2
		} else if text[i] == cr || text[i] == lf {
			frecord(nStart, i)
			nStart = i + 1
		}
	}
	return entryTotal, entryStart, entryEnd, nil
}

func (t *tabfile) Pretreat(trans []byte) []byte {
	return trans
}
