package calc

func SolveExample(ex Example) (float64, error) {
	if ex.Second_value == 0 {
		return 0, ErrDivideByZero
	}

	switch ex.Operation {
	case Plus:
		return ex.First_value + ex.Second_value, nil
	case Minus:
		return ex.First_value - ex.Second_value, nil
	case Multiply:
		return ex.First_value * ex.Second_value, nil
	case Division:
		return ex.First_value / ex.Second_value, nil
	case Equals:
		return ex.First_value, nil
	}
	return 0, ErrExpressionIncorrect
}
