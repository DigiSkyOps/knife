package global

var (
	ServerSig chan int
	UpdateSig chan int
)

func init() {
	ServerSig = make(chan int)
	UpdateSig = make(chan int)
}
