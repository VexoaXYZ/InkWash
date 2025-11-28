package animation

import "math"

// EasingFunc is a function that takes a time value (0.0 to 1.0) and returns an eased value
type EasingFunc func(t float64) float64

// Linear easing - no easing, constant rate
func Linear(t float64) float64 {
	return t
}

// EaseInOutCubic - Cubic bezier for smooth, professional motion
// Accelerates until halfway, then decelerates
func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

// EaseOutExpo - For quick, snappy interactions
// Decelerates exponentially
func EaseOutExpo(t float64) float64 {
	if t == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}

// EaseInExpo - Accelerates exponentially
func EaseInExpo(t float64) float64 {
	if t == 0 {
		return 0
	}
	return math.Pow(2, 10*(t-1))
}

// EaseInOutSine - For gentle, organic movement
func EaseInOutSine(t float64) float64 {
	return -(math.Cos(math.Pi*t) - 1) / 2
}

// EaseOutSine - Gentle deceleration
func EaseOutSine(t float64) float64 {
	return math.Sin(t * math.Pi / 2)
}

// EaseInSine - Gentle acceleration
func EaseInSine(t float64) float64 {
	return 1 - math.Cos(t*math.Pi/2)
}

// EaseInOutQuad - Quadratic easing in and out
func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseOutQuad - Quadratic deceleration
func EaseOutQuad(t float64) float64 {
	return t * (2 - t)
}

// EaseInQuad - Quadratic acceleration
func EaseInQuad(t float64) float64 {
	return t * t
}

// EaseInOutQuart - Quartic easing in and out
func EaseInOutQuart(t float64) float64 {
	if t < 0.5 {
		return 8 * t * t * t * t
	}
	t = t - 1
	return 1 - 8*t*t*t*t
}

// EaseOutBounce - Bouncing effect (for playful interactions)
func EaseOutBounce(t float64) float64 {
	const n1 = 7.5625
	const d1 = 2.75

	if t < 1/d1 {
		return n1 * t * t
	} else if t < 2/d1 {
		t = t - 1.5/d1
		return n1*t*t + 0.75
	} else if t < 2.5/d1 {
		t = t - 2.25/d1
		return n1*t*t + 0.9375
	} else {
		t = t - 2.625/d1
		return n1*t*t + 0.984375
	}
}

// EaseInBack - Slight back movement before forward motion
func EaseInBack(t float64) float64 {
	const c1 = 1.70158
	const c3 = c1 + 1

	return c3*t*t*t - c1*t*t
}

// EaseOutBack - Slight overshoot then settle
func EaseOutBack(t float64) float64 {
	const c1 = 1.70158
	const c3 = c1 + 1

	return 1 + c3*math.Pow(t-1, 3) + c1*math.Pow(t-1, 2)
}

// EaseInOutBack - Back easing in and out
func EaseInOutBack(t float64) float64 {
	const c1 = 1.70158
	const c2 = c1 * 1.525

	if t < 0.5 {
		return (math.Pow(2*t, 2) * ((c2+1)*2*t - c2)) / 2
	}

	return (math.Pow(2*t-2, 2)*((c2+1)*(t*2-2)+c2) + 2) / 2
}
