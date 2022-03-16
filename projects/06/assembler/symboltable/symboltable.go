package symboltable

type SymboolTable struct {
	count int
	data  map[string]int
}

var (
	symboltable = SymboolTable{
		count: 0,
		data:  make(map[string]int),
	}
)

func GetSymbolTable() *SymboolTable {
	return &symboltable
}

// シンボルテーブルに(symbol, address)のペアを追加する。
func (st *SymboolTable) AddEntry(symbol string, address int) {
	st.data[symbol] = address
}

// シンボルテーブルは与えらたsymbolを含むか。
func (st *SymboolTable) Contains(symbol string) bool {
	_, ok := st.data[symbol]
	return ok
}

// symbolに結び付けられたアドレスを返す。
func (st *SymboolTable) GetAddress(symbol string) int {
	return st.data[symbol]
}

func (st *SymboolTable) GetCount() int {
	return st.count
}

func (st *SymboolTable) InitCount() {
	st.count = 0
}

func (st *SymboolTable) Increment() {
	st.count++
}
