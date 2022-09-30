package kurento

type RembParams struct {
	PacketsRecvIntervalTop int
	ExponentialFactor      float64
	LinealFactorMin        int
	LinealFactorGrade      float64
	DecrementFactor        float64
	ThresholdFactor        float64
	UpLosses               int
	RembOnConnect          int
}
