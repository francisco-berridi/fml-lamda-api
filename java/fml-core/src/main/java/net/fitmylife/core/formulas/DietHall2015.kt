package net.fitmylife.core.formulas

import net.fitmylife.core.enums.Gender

/**
 * Note tha None of the variables
 * respect the naming standards, this is due the fact that this is a translation
 * form the origin Javascript code.
 */
object DietHall2015 {

    fun run(input: DietHall205Input): DietHall205Result {

        var averageEI = input.rmr * input.pal

        var NewC = input.carbsProportionNewDiet * averageEI
        var NewF = input.fatProportionNewDiet * averageEI
        var NewP = input.proteinProportionNewDiet * averageEI

        val ExtraExCal = 0.0
        val hypLevel = 0.0

        var result = Rk4onePhaseHypHall2015().run(
            Rk4onePhaseHypHall2015Input(
                input.gender,
                input.heightCMS,
                input.initialFMKg,
                input.initialBWKg,
                input.pal,
                input.rmr,
                input.carbsProportionNewDiet,
                input.fatProportionNewDiet,
                input.proteinProportionNewDiet,
                input.prop_pi,
                NewC,
                NewF,
                NewP,
                0.4,
                0.0,
                ExtraExCal,
                hypLevel,
                input.age,
                input.daysBetweenInitialAndFinal,
                input.step
            )
        )


        var BWlast = result.BW.last()
        var errorFactor = input.finalBW * 1000 - BWlast

        var averageEI_UP: Double
        var averageEI_LOW: Double

        averageEI_UP = if (errorFactor > 0) {
            4 * averageEI
        } else {
            1.2 * averageEI
        }

        averageEI_LOW = if (errorFactor > 0) {
            0.8 * averageEI
        } else {
            averageEI / 5
        }

        while (Math.abs(errorFactor) > 50) {
            averageEI = randomInt(averageEI_LOW.toInt(), averageEI_UP.toInt()).toDouble()

            NewC = input.carbsProportionNewDiet * averageEI
            NewF = input.fatProportionNewDiet * averageEI
            NewP = input.proteinProportionNewDiet * averageEI

            result = Rk4onePhaseHypHall2015().run(
                Rk4onePhaseHypHall2015Input(
                    input.gender,
                    input.heightCMS,
                    input.initialFMKg,
                    input.initialBWKg,
                    input.pal,
                    input.rmr,
                    input.carbsProportionNewDiet,
                    input.fatProportionNewDiet,
                    input.proteinProportionNewDiet,
                    input.prop_pi,
                    NewC,
                    NewF,
                    NewP,
                    0.4,
                    0.0,
                    ExtraExCal,
                    hypLevel,
                    input.age,
                    input.daysBetweenInitialAndFinal,
                    input.step
                )
            )

            // del result busca  el ultimo BW y lo compara con el final BW para (en gramos) para tratar de encontrar "errorFacotor" < 50
            BWlast = result.BW.last()
            errorFactor = input.finalBW * 1000 - BWlast

            if (errorFactor > 0) {
                averageEI_LOW = averageEI
            } else {
                averageEI_UP = averageEI
            }

            if ((averageEI_UP - averageEI_LOW) < 1) {
                throw Exception("There was an issue processing your weight data. Showing results based on your most recent verified entry")
            }
        }

        val dietResult = DietHall205Result(
            BW = DoubleArray(0),
            Fa = DoubleArray(0),
            SM = DoubleArray(0),
            RMR = input.rmr,
            AverageEnergyIntake = averageEI,
            PropC = input.carbsProportionNewDiet,
            PropF = input.fatProportionNewDiet,
            PropP = input.proteinProportionNewDiet,
            UserDependantParameters = result.initializationHall2015ODESResult,
            LastValuesODEFittedModel = ODEFittedModelValues(
                G = result.G.last(),
                Psig = result.Psig.last(),
                P = result.P.last(),
                L_diet = result.L_diet.last(),
                Fa = result.Fa.last(),
                ECF_F = result.ECF_F.last(),
                ECF_S = result.ECF_S.last(),
                Therm = result.Therm.last()

            )
        )

        when (input.predictionDays) {
            0 -> {
                dietResult.BW = doubleArrayOf(result.BW.last())
                dietResult.Fa = doubleArrayOf(result.Fa.last())
                dietResult.SM = doubleArrayOf(result.SM.last())
            }
            else -> {
                val predictionResult = Rk4onePhaseHypHall2015().run(
                    Rk4onePhaseHypHall2015Input(
                        input.gender,
                        input.heightCMS,
                        input.initialFMKg,
                        input.initialBWKg,
                        input.pal,
                        input.rmr,
                        input.carbsProportionNewDiet,
                        input.fatProportionNewDiet,
                        input.proteinProportionNewDiet,
                        input.prop_pi,
                        NewC,
                        NewF,
                        NewP,
                        0.4,
                        0.0,
                        ExtraExCal,
                        hypLevel,
                        input.age,
                        input.daysBetweenInitialAndFinal + input.predictionDays,
                        input.step
                ))

                dietResult.BW = predictionResult.BW.sliceArray(0..input.daysBetweenInitialAndFinal)
                dietResult.Fa = predictionResult.Fa.sliceArray(0..input.daysBetweenInitialAndFinal)
                dietResult.SM = predictionResult.SM.sliceArray(0..input.daysBetweenInitialAndFinal)
            }
        }

        return dietResult
    }
}

data class DietHall205Input(
    val age: Int,
    val gender: Gender,
    val heightCMS: Double,
    val pal: Double,
    val prop_pi: Double,
    val rmr: Double,
    val initialBWKg: Double,
    val initialFMKg: Double,
    val finalBW: Double,
    val carbsProportionNewDiet: Double,
    val fatProportionNewDiet: Double,
    val proteinProportionNewDiet: Double,
    val daysBetweenInitialAndFinal: Int,
    val predictionDays: Int,
    val step: Double
)

data class DietHall205Result(
        val AverageEnergyIntake: Double,
        val PropC: Double,
        val PropF: Double,
        val PropP: Double,
        val RMR: Double,
        var BW: DoubleArray,
        var Fa: DoubleArray,
        var SM: DoubleArray,
        val UserDependantParameters: InitializationHall2015ODESResult,
        val LastValuesODEFittedModel: ODEFittedModelValues
) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false

        other as DietHall205Result

        if (AverageEnergyIntake != other.AverageEnergyIntake) return false
        if (PropC != other.PropC) return false
        if (PropF != other.PropF) return false
        if (PropP != other.PropP) return false
        if (RMR != other.RMR) return false
        if (!BW.contentEquals(other.BW)) return false
        if (!Fa.contentEquals(other.Fa)) return false
        if (!SM.contentEquals(other.SM)) return false
        if (UserDependantParameters != other.UserDependantParameters) return false
        if (LastValuesODEFittedModel != other.LastValuesODEFittedModel) return false

        return true
    }

    override fun hashCode(): Int {
        var result = AverageEnergyIntake.hashCode()
        result = 31 * result + PropC.hashCode()
        result = 31 * result + PropF.hashCode()
        result = 31 * result + PropP.hashCode()
        result = 31 * result + RMR.hashCode()
        result = 31 * result + BW.contentHashCode()
        result = 31 * result + Fa.contentHashCode()
        result = 31 * result + SM.contentHashCode()
        result = 31 * result + UserDependantParameters.hashCode()
        result = 31 * result + LastValuesODEFittedModel.hashCode()
        return result
    }
}

data class ODEFittedModelValues(
    val G: Double,
    val Psig: Double,
    val P: Double,
    val L_diet: Double,
    val Fa: Double,
    val ECF_F: Double,
    val ECF_S: Double,
    val Therm: Double
)