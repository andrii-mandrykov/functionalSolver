package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func main() {
	processResult(solveEquation)(receiveCoefficients)
}

/*
 * Higher order function that receive function as argument and return function as result.
 * Receives coefficients, solves the equation and returns results.
 */
func processResult(solveEquation func(func() ([]float64, error)) ([]string, error)) func(receiveCoefficients func() ([]float64, error)) ([]string, error) {
	return func(receiveCoefficients func() ([]float64, error)) ([]string, error) {
		roots, err := solveEquation(receiveCoefficients)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(roots)
		}

		return roots, err
	}
}

/*
 * Higher order function that receive function as argument.
 * Returns real roots of equation.
 */
func solveEquation(receiveCoefficients func() ([]float64, error)) ([]string, error) {
	coefficients, err := receiveCoefficients()

	if err != nil {
		return []string{}, err
	}

	/*
	 * Recursion used
	 */
	if coefficients[0] == 0 {
		return solveEquation(func() ([]float64, error) {
			return coefficients[1:], nil
		})
	}

	switch len(coefficients) {
	case 2:
		return solveLinearEquation(coefficients[0], coefficients[1])
	case 3:
		return solveQuadraticEquation(coefficients[0], coefficients[1], coefficients[2])
	case 4:
		return solveCubicEquation(coefficients[0], coefficients[1], coefficients[2], coefficients[3])
	case 5:
		return solveQuarticEquation(coefficients[0], coefficients[1], coefficients[2], coefficients[3], coefficients[4])
	default:
		return []string{}, errors.New("there should be at least 2 coefficients in the equation, 1 provided")
	}
}

/*
 * First class function that used as argument.
 * Returns entered coefficients.
 */
func receiveCoefficients() ([]float64, error) {
	coefficients := [6]float64{}

	fmt.Print("Enter coefficients of your equation = (a, b, etc...) ")
	argsCounter, err := fmt.Scanln(&coefficients[0], &coefficients[1], &coefficients[2], &coefficients[3], &coefficients[4], &coefficients[5])

	if err != nil && err.Error() != "unexpected newline" {
		if argsCounter > 5 {
			return []float64{}, errors.New("there is no algebraic solutions of general quintic equations and more complex ones")
		}

		return []float64{}, errors.New("an error occured while reading input" + (string(argsCounter)))
	}

	if argsCounter < 2 {
		return []float64{}, errors.New("there should be at least 2 coefficients in the equation, " + (string(argsCounter)) + "provided")
	}

	if argsCounter > 5 {
		return []float64{}, errors.New("there is no algebraic solutions of general quintic equations and more complex ones")
	}

	return  coefficients[:argsCounter], nil
}

/*
 * Pure function.
 * Returns real roots of linear equation of type ax + b = 0
 */
func solveLinearEquation(a, b float64) ([]string, error) {
	if a== 0 {
		return []string{}, errors.New("the first argument of linear equation could not be zero")
	}

	return []string{formatFloatValue( -b/a)}, nil
}

/*
 * Pure function.
 * Returns real roots of quadratic equation of type ax2 + bx + c = 0
 */
func solveQuadraticEquation(a, b, c float64) ([]string, error) {
	discriminant := (b * b) - (4 * a * c)

	if discriminant > 0 {
		return []string{
			formatFloatValue( (-b + math.Sqrt(discriminant))/(2*a)),
			formatFloatValue( (-b - math.Sqrt(discriminant))/(2*a)),
		}, nil
	}

	if discriminant == 0 {
		return []string{
			formatFloatValue( -b/2*a),
		}, nil
	}

	return []string{}, errors.New("there are no roots for this equation")
}

/*
 * Pure function.
 * Returns real roots of cubic equation of type ax3 + bx2 + cx +d = 0
 */
func solveCubicEquation(a, b, c, d float64) ([]string, error) {
	invertedA := 1.0/a
	invertedB := invertedA * b
	invertedC := invertedA * c
	invertedD := invertedA * d

	Q := (3 * invertedC - invertedB * invertedB) / 9
	R := (9 * invertedB * invertedC - 27 * invertedD - 2 * invertedB * invertedB * invertedB) / 54
	D := Q * Q * Q + R * R

	if Q == 0 {
		if R == 0 {
			return []string{
				formatFloatValue(- invertedB / 3.0),
			}, nil
		} else {
			return []string{
				formatFloatValue(math.Pow(2 * R, 1.0/3) - invertedB / 3),
			}, nil
		}
	}

	if D <= 0 {
		T := math.Acos(R / math.Sqrt(-Q * Q * Q))

		return []string{
			formatFloatValue(2 * math.Sqrt(-Q) * math.Cos(T / 3.0) - invertedB / 3),
			formatFloatValue(2 * math.Sqrt(-Q) * math.Cos((T + 2 * math.Pi) / 3.0) - invertedB / 3),
			formatFloatValue(2 * math.Sqrt(-Q) * math.Cos((T + 4 * math.Pi) / 3.0) - invertedB / 3),
		}, nil
	}

	if R == 0 {
		return []string{
			formatFloatValue(- invertedB / 3.0),
		}, nil
	} else if R > 0 {
		AD := math.Pow(math.Abs(R) + math.Sqrt(D), 1.0/3)

		return []string{
			formatFloatValue(AD - Q / AD - invertedB / 3.0),
		}, nil
	} else {
		AD := -1 * math.Pow(math.Abs(R) + math.Sqrt(D), 1.0/3)

		return []string{
			formatFloatValue(AD - Q / AD - invertedB / 3.0),
		}, nil
	}
}

/*
 * Pure function.
 * Returns real roots of quartic equation of type ax4 + bx3 + cx2 +dx +e = 0
 */
func solveQuarticEquation(a, b, c, d, e float64) ([]string, error) {
	invertedA := 1.0/a
	invertedB := invertedA * b
	invertedC := invertedA * c
	invertedD := invertedA * d
	invertedE := invertedA * e

	cubicResult, _ := solveCubicEquation(1, -invertedC, invertedD * invertedB - 4 * invertedE, 4 * invertedC * invertedE - invertedD * invertedD - invertedB * invertedB * invertedE)

	if len(cubicResult) == 0 {
		return []string{}, errors.New("there are no roots for this equation")
	}

	firstCubicResult, _ := strconv.ParseFloat(cubicResult[0], 64)
	R2 := 0.25 * invertedB * invertedB - invertedC + firstCubicResult

	if R2 < 0 {
		return []string{}, errors.New("there are no real roots for this equation")
	}

	R := math.Sqrt(R2)
	invertedR := 1.0/R

	var D2, E2 float64

	if R < math.Pow(10, -12) {
		T := firstCubicResult * firstCubicResult - 4 * invertedE

		if T < 0 {
			D2 = -1
			E2 = -1
		} else {
			D2 = 0.75 * invertedB * invertedB - 2 * invertedC + 2 * math.Sqrt(T)
			E2 = D2 - 4 * math.Sqrt(T)
		}
	} else {
		u := 0.75 * invertedB * invertedB - 2 * invertedC - R2
		v := 0.25 * invertedR * (4 * invertedB * invertedC - 8 * invertedD - invertedB * invertedB * invertedB)

		D2 = u + v
		E2 = u - v
	}

	if E2 >= 0 {
		if D2 >= 0 {
			RR := 0.5 * R + 0.5 * math.Sqrt(D2) - 0.25 * invertedB
			EE := - 0.5 * R + math.Sqrt(E2) * 0.5 - 0.25 * invertedB

			return []string{
				formatFloatValue(RR),
				formatFloatValue(RR - math.Sqrt(D2)),
				formatFloatValue(EE),
				formatFloatValue(EE - math.Sqrt(E2)),
			}, nil
		} else {
			EE := - 0.5 * R + math.Sqrt(E2) * 0.5 - 0.25 * invertedB

			return []string{
				formatFloatValue(EE),
				formatFloatValue(EE - math.Sqrt(E2)),
			}, nil
		}
	} else {
		if D2 >= 0 {
			RR := 0.5 * R + 0.5 * math.Sqrt(D2) - 0.25 * invertedB

			return []string{
				formatFloatValue(RR),
				formatFloatValue(RR - math.Sqrt(D2)),
			}, nil
		}
	}

	return []string{}, errors.New("there are no real roots for this equation")
}

/*
 * Pure function.
 * Returns formatted string value of number.
 */
func formatFloatValue(value float64) string {
	if value == float64(int(value)) {
		return fmt.Sprintf("%.0f", value)
	}

	if math.Abs(value - float64(int(value))) <= 0.001 {
		return fmt.Sprintf("%.0f", value)
	}

	return fmt.Sprintf("%.2f", value)
}