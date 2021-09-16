package main

type Estimator interface{
	predict() float32
	update(bool)
}

type Context8 struct{
	b byte
	c byte
}

type Counter struct{
	n0, n1 uint8
}

type Estimator8 struct{
	freq map[Context8]Counter
	ctxt Context8
	nbits int
}

func newEstimator8() *Estimator8 {
	var E Estimator8
	E.freq = make(map[Context8]Counter)
	return &E
}

func (E *Estimator8) predict() float32{
	ctr, ok := E.freq[E.ctxt]
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

func (E *Estimator8) update(bit bool) {
	ctr, ok := E.freq[E.ctxt]
	if !ok {
		ctr = Counter{0,0}
	}
	ctr.increment(bit)
	E.freq[E.ctxt] = ctr

	if E.nbits >= 8{
		c := E.ctxt.c
		E.ctxt = Context8{c,0}
		E.nbits = 0
	}
	E.ctxt.c = 2*E.ctxt.c + btoi(bit)
	E.nbits++
}

func (ctr *Counter) increment(bit bool){
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
