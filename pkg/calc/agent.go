package calc

func SolveExample(ex Example) (float64, error) {
	if ex.SecondArgument.Value == 0 && ex.Operation == Division {
		return 0, ErrDivideByZero
	}

	switch ex.Operation {
	case Plus:
		return ex.FirstArgument.Value + ex.SecondArgument.Value, nil
	case Minus:
		return ex.FirstArgument.Value - ex.SecondArgument.Value, nil
	case Multiply:
		return ex.FirstArgument.Value * ex.SecondArgument.Value, nil
	case Division:
		return ex.FirstArgument.Value / ex.SecondArgument.Value, nil
	case Equals:
		return ex.FirstArgument.Value, nil
	}
	return 0, ErrExpressionIncorrect
}
