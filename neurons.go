package go_deep

type inputLayer interface {
	forward(setItem []float64) (output float64)
}

type firstHiddenLayer interface {
	activation
	cost
	init()
	forward(float64) [][]float64
	forwardMeasure([]float64, []float64) ([][]float64, float64)
	backward([][]float64)
	applyCorrections(float64)
}

//type hiddenLayer interface {
	//activation
	//cost
	//forward(arg) return_val
// TODO: compute actual correction with hidden layer Delata rule
// https://theclevermachine.wordpress.com/2014/09/06/derivation-error-backpropagation-gradient-descent-for-neural-networks/
	//backward(arg) return_val
	//init(arg) return_val
	//applyCorrection()
//}

type outputLayer interface {
	activation
	cost
	forward(rowInput [][]float64) []float64
	forwardMeasure([]float64, []float64) ([][]float64, float64)
	backward(labels []float64) [][]float64
}

type inputDense struct {}

// Optimize input layer and create firstHidden layer and extraHidden layer with different input vector shape
func (l *inputDense) farward(setItem []float64) (output float64) {
	/*
	The Input nodes provide information from the outside world to the 
	network and are together referred to as the “Input Layer”. No computation
	is performed in any of the Input nodes – they just pass on the information to the hidden nodes.
	*/
	for _, i := range setItem {
		output += i
	}
	// TODO: don't do so in nn, prepare data outside. Raise an error instead if sum of signals is InF
	output *= .00001
	return
}

type hiddenDenseFirst struct {
	activation
	synapseInitializer
	currLayerSize, nextLayerSize int // Length of neurons sequence - 1
	learningRate float64
	corrections, synapses [][]float64
}

func (l *hiddenDenseFirst) init() {
	l.synapses = l.synapseInitializer.init()
}

func (l *hiddenDenseFirst) forward(input float64) (output [][]float64) {
	// Each neuron of a first hidden layer receives a sum of all input signals from an input later and activates it.
	// Computation of first hidden layer cost value has no sense because before multiplication of activated sum on
	// synapses all neurons have the same value - activated sum of incoming signal. It is true because input layer
	// has no weights.
	activated := l.activate(input)

	for i := 0; i < l.nextLayerSize; i++ {
		for j := 0; j < l.currLayerSize - 1; j++ {
			if output[i] == nil {
				output[i] = make([]float64, l.currLayerSize)
			}
			// Transition between layers is a matrix reshape. Way or another reshape matrix is required on step of multiplication or sum.
			output[i][j] = l.synapses[j][i] * activated
		}
		output[i][l.currLayerSize] = l.synapses[l.currLayerSize][i] // Add i bias to the sum of weighted output. Bias doesn't use signal, bias is a weight without input.
	}
	return output
}

// A high-grade i.e. extra hidden layer collects corrections (incoming errors)
// then sum per neuron incoming errors (alongside) and computes errors for
// a next hidden layer.
// First hidden layer doesn't have a previous hidden layer so it doesn't compute
// errors (corrections) for synapses between previous layer and an actual one.
// Instead of it, it is just collects errors for correction synapses between itself
// and a next layer (possibly) output.
func (l *hiddenDenseFirst) backward(eRRors [][]float64) {
	for i, eRR := range eRRors {
		for j, c := range eRR {
			l.corrections[i][j] += c
		}
	}
}

func (l *hiddenDenseFirst) applyCorrections(batchSize float64) {
	for i, corr := range l.corrections {
		for j, c := range corr {
			l.synapses[i][j] += l.learningRate * c / batchSize
		}
	}
}


type outputDense struct {
	activation
	// Cost function exists only in output layer and in hidden layers used indirectly
	// as a sum of wighted errors. Thus cost function is global for a network.
	cost
	prevLayerSize int
	input [][]float64
}

func (l *outputDense) forward(rowInput [][]float64) (output []float64) {
	l.input = rowInput
	var iSum float64

	for _, raw := range rowInput {
		iSum = 0
		for _, item := range raw {
			iSum += item
		}
		output = append(output, l.activate(iSum))
	}
	return
}

func (l *outputDense) forwardMeasure(rowInput [][]float64, labels []float64) (prediction []float64, cost float64) {
	prediction = l.forward(rowInput)
	cost = l.countCost(prediction, labels)
	return
}


func (l *outputDense) backward(prediction []float64, labels []float64) (corrections [][]float64) {
	var cost, zk float64

	for i, ak := range prediction {
		zk = 0
		for _, aj := range l.input[i] {
			zk += aj // Sum current layer input
		}
		// Delta rule
		cost = l.costDerivative(ak, labels[i]) * l.actDerivative(zk)
		for k := 0; k < l.prevLayerSize; k++ {
			// Corrections vector of the same shape as synapses vector
			corrections[k][i] = cost * l.input[i][k]
		}
		// Add bias correction
		corrections[l.prevLayerSize][i] = cost
	}
	return
}