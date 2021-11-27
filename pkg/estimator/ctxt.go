package estimator

type Estimator interface{
	Predict() float32
	Update(bool)
}

type Context8 struct{
	Cur byte
	Last byte
}

type Counter struct{
	n0, n1 uint8
}

type Estimator8 struct{
	Freq map[Context8]Counter
	Ctxt Context8
	nbits int
}

func NewEstimator8() *Estimator8 {
	var E Estimator8
	E.Freq = make(map[Context8]Counter)
	return &E
}

func (E *Estimator8) Predict() float32{
	ctr, ok := E.Freq[E.Ctxt]
	if ok {
		return (float32(ctr.n1)+1)/(float32(ctr.n0)+float32(ctr.n1)+2)
	}else{
		return 0.5
	}
}

func btoi(b bool) uint8 {
    if b {
        return 1
    }
    return 0
 }

func (E *Estimator8) Update(bit bool) {
	ctr, ok := E.Freq[E.Ctxt]
	if !ok {
		ctr = Counter{0,0}
	}
	ctr.Increment(bit)
	E.Freq[E.Ctxt] = ctr

	if E.nbits >= 8{
		last_byte := E.Ctxt.Last
		E.Ctxt = Context8{last_byte,0}
		E.nbits = 0
	}
	E.Ctxt.Cur = 2*E.Ctxt.Cur + btoi(bit)
	E.nbits++
}

func (ctr *Counter) Increment(bit bool){
	if bit {
		ctr.n1++
		if ctr.n0 > 2 {
			ctr.n0 = (ctr.n0+1)/2
		}
	} else {
		ctr.n0++
		if ctr.n1 > 2{
			ctr.n1 = (ctr.n1+1)/2
		}
	}
	if ctr.n1 >= 254 || ctr.n0 >= 254 {
		ctr.n1 /= 2
		ctr.n0 /= 2
	}
}
