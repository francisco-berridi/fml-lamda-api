package formulas

import (
	"enums/Gender"
	"fmt"
	"math"
)

type ValuesODEFittedModel struct {
	G      float64
	Psig   float64
	P      float64
	L_diet float64
	Fa     float64
	ECF_F  float64
	ECF_S  float64
	Therm  float64
}

type GetDietHall2015Result struct {
	AverageEnergyIntake      float64
	PropC                    float64
	PropF                    float64
	PropP                    float64
	RMR                      float64
	BW                       []float64
	Fa                       []float64
	SM                       []float64
	UserDependparamaters     InitializationHall2015ODESResult
	LastValuesODEFittedModel ValuesODEFittedModel
}

func GetDietHall2015(age int64, gender Gender.Enum, heightCMS, pal, prop_pi, rmr, initialBWKg,
	initialFMKg, finalBW, carbsProportionNewDiet, fatProportionNewDiet, proteinProportionNewDiet float64, daysBetweenInitialAndFinal, predictionDays int64, step float64) (GetDietHall2015Result, error) {

	var averageEI = rmr * pal

	var NewC = carbsProportionNewDiet * averageEI
	var NewF = fatProportionNewDiet * averageEI
	var NewP = proteinProportionNewDiet * averageEI

	var ExtraExCal = 0.0
	var hypLevel = 0.0

	result := Rk4onePhaseHypHall2015(gender, heightCMS, initialFMKg, initialBWKg, pal, rmr, carbsProportionNewDiet, fatProportionNewDiet, proteinProportionNewDiet, prop_pi,
		NewC, NewF, NewP, 0.4, 0, ExtraExCal, hypLevel, age, daysBetweenInitialAndFinal, step)

	var BWlast = result.BW[len(result.BW)-1]
	var errorFactor = finalBW*1000 - BWlast

	var averageEI_UP float64
	var averageEI_LOW float64

	if errorFactor > 0 {
		averageEI_UP = 4 * averageEI
	} else {
		averageEI_UP = 1.2 * averageEI
	}

	if errorFactor > 0 {
		averageEI_LOW = 0.8 * averageEI
	} else {
		averageEI_LOW = averageEI / 5
	}

	for math.Abs(errorFactor) > 50 {

		averageEI = float64(GetRandomInt(averageEI_LOW, averageEI_UP))
		NewC = carbsProportionNewDiet * averageEI
		NewF = fatProportionNewDiet * averageEI
		NewP = proteinProportionNewDiet * averageEI

		// Esta es la que va intentando differentes trajectorias, para tratar de encontrar el valor adecuado


		result = Rk4onePhaseHypHall2015(gender, heightCMS, initialFMKg, initialBWKg, pal, rmr, carbsProportionNewDiet, fatProportionNewDiet,
			proteinProportionNewDiet, prop_pi, NewC, NewF, NewP, 0.4, 0, ExtraExCal, hypLevel, age, daysBetweenInitialAndFinal, step)

		// del result busca  el ultimo BW y lo compara con el final BW para (en gramos) para tratar de encontrar "errorFacotor" < 50
		BWlast = result.BW[len(result.BW)-1]
		errorFactor = finalBW*1000 - BWlast

		if errorFactor > 0 {
			averageEI_LOW = averageEI
		} else {
			averageEI_UP = averageEI
		}

		if (averageEI_UP - averageEI_LOW) < 1 {
			errorFactor = 0
			return GetDietHall2015Result{}, fmt.Errorf("There was an issue processing your weight data. Showing results based on your most recent verified entry")
		}

	}

	dietResult := GetDietHall2015Result{
		RMR:                  rmr,
		AverageEnergyIntake:  averageEI,
		PropC:                carbsProportionNewDiet,
		PropF:                fatProportionNewDiet,
		PropP:                proteinProportionNewDiet,
		UserDependparamaters: result.UserDependparamaters,
		LastValuesODEFittedModel: ValuesODEFittedModel{
			G:      result.G[len(result.G)-1],
			Psig:   result.Psig[len(result.Psig)-1],
			P:      result.P[len(result.P)-1],
			L_diet: result.L_diet[len(result.L_diet)-1],
			Fa:     result.Fa[len(result.Fa)-1],
			ECF_F:  result.ECF_F[len(result.ECF_F)-1],
			ECF_S:  result.ECF_S[len(result.ECF_S)-1],
			Therm:  result.Therm[len(result.Therm)-1],
		},
	}

	if predictionDays == 0 {
		dietResult.BW = []float64{result.BW[len(result.BW)-1]}
		dietResult.Fa = []float64{result.Fa[len(result.Fa)-1]}
		dietResult.SM = []float64{result.SM[len(result.SM)-1]}

	} else {

		predictionResult := Rk4onePhaseHypHall2015(gender, heightCMS, initialFMKg, initialBWKg, pal, rmr, carbsProportionNewDiet, fatProportionNewDiet,
			proteinProportionNewDiet, prop_pi, NewC, NewF, NewP, 0.4, 0, ExtraExCal, hypLevel, age, daysBetweenInitialAndFinal+predictionDays, step)

		dietResult.BW = predictionResult.BW[0:daysBetweenInitialAndFinal]
		dietResult.Fa = predictionResult.Fa[0:daysBetweenInitialAndFinal]
		dietResult.SM = predictionResult.Fa[0:daysBetweenInitialAndFinal]

	}

	return dietResult, nil
}
