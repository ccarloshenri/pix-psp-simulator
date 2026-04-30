package interfaces

type IDGenerator interface {
	GenerateTxID() string
	GenerateE2EID(ispb string) string
	GenerateRtrID(ispb string) string
}
