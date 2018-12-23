package net.fitmyife.core.formulas

import net.fitmyife.core.formulas.params.Hg
import net.fitmyife.core.formulas.params.Hp
import net.fitmylife.core.enums.Gender

class Rk4onePhaseHypHall2015 {

    fun run(input: Rk4onePhaseHypHall2015Input): Rk4onePhaseHypHall2015Result {

        val userDependantParameters = InitializationHall2015ODES().run(
            InitializationHall2015ODESInput(
                input.gender,
                input.age,
                input.heightCMS,
                input.pal,
                input.rmr,
                input.initialFMKg,
                input.initialBWKg,
                input.propC,
                input.propF,
                input.propP,
                input.prop_pi,
                input.newC,
                input.newF,
                input.newP,
                input.psi,
                input.ExtraExCal
            )
        )

        val GValue = arrayListOf(userDependantParameters.Ginit)
        val PsigValue = arrayListOf(0.0)
        val PValue = arrayListOf(userDependantParameters.Pinit)
        val LdietValue = arrayListOf(0.0)
        val FaValue = arrayListOf(userDependantParameters.Finit)
        val ECFFValue = arrayListOf(0.0)
        val ECFSValue = arrayListOf(0.0)
        val ThermValue = arrayListOf(0.0)

        val cycles = (input.predictionDays / input.step).toInt()


        val kG = DoubleArray(4)
        val kPsig = DoubleArray(4)
        val kP = DoubleArray(4)
        val kLdiet = DoubleArray(4)
        val kFa = DoubleArray(4)
        val kECFF = DoubleArray(4)
        val kECFS = DoubleArray(4)
        val kTherm = DoubleArray(4)


        for (i in 0 until cycles) {

            //dG, dPsig, dP, dL_diet, dFa, dECF_F, dECF_S, dTherm
            val hall2015ODESonlyResult = Hall2015ODESonly().run(
                Hall2015ODESonlyInput(
                    userDependantParameters, input.gender, input.S_F, GValue[i],
                    PsigValue[i], PValue[i], LdietValue[i], FaValue[i], ECFFValue[i], ECFSValue[i], ThermValue[i]
                )
            )
            kG[0] = hall2015ODESonlyResult.dG
            kPsig[0] = hall2015ODESonlyResult.dPsig
            kP[0] = hall2015ODESonlyResult.dP
            kLdiet[0] = hall2015ODESonlyResult.dL_diet
            kFa[0] = hall2015ODESonlyResult.dFa
            kECFF[0] = hall2015ODESonlyResult.dECF_F
            kECFS[0] = hall2015ODESonlyResult.dECF_S
            kTherm[0] = hall2015ODESonlyResult.dTherm

            for (k in 1..3) {

                var h = 0.5
                // half
                if (k == 3) {
                    // the last one doest gets 0.5, it a 1.0
                    h = 1.0
                }

                val extraHall2015ODESonlyResult = Hall2015ODESonly().run(
                    Hall2015ODESonlyInput(
                        userDependantParameters,
                        input.gender,
                        input.S_F,
                        GValue[i] + h * input.step * kG[k - 1],
                        PsigValue[i] + h * input.step * kPsig[k - 1],
                        PValue[i] + h * input.step * kP[k - 1],
                        LdietValue[i] + h * input.step * kLdiet[k - 1],
                        FaValue[i] + h * input.step * kFa[k - 1],
                        ECFFValue[i] + h * input.step * kECFF[k - 1],
                        ECFSValue[i] + h * input.step * kECFS[k - 1],
                        ThermValue[i] + h * input.step * kTherm[k - 1]
                    )
                )


                kG[k] = extraHall2015ODESonlyResult.dG
                kPsig[k] = extraHall2015ODESonlyResult.dPsig
                kP[k] = extraHall2015ODESonlyResult.dP
                kLdiet[k] = extraHall2015ODESonlyResult.dL_diet
                kFa[k] = extraHall2015ODESonlyResult.dFa
                kECFF[k] = extraHall2015ODESonlyResult.dECF_F
                kECFS[k] = extraHall2015ODESonlyResult.dECF_S
                kTherm[k] = extraHall2015ODESonlyResult.dTherm
            }

            GValue.add(GValue[i] + (1.0 / 6.0) * (kG[0] + 2.0 * kG[1] + 2.0 * kG[2] + kG[3]) * input.step)
            PsigValue.add(PsigValue[i] + (1.0 / 6.0) * (kPsig[0] + 2.0 * kPsig[1] + 2.0 * kPsig[2] + kPsig[3]) * input.step)
            PValue.add(PValue[i] + (1.0 / 6.0) * (kP[0] + 2.0 * kP[1] + 2.0 * kP[2] + kP[3]) * input.step)
            LdietValue.add(LdietValue[i] + (1.0 / 6.0) * (kLdiet[0] + 2.0 * kLdiet[1] + 2.0 * kLdiet[2] + kLdiet[3]) * input.step)
            FaValue.add(FaValue[i] + (1.0 / 6.0) * (kFa[0] + 2.0 * kFa[1] + 2.0 * kFa[2] + kFa[3]) * input.step)
            ECFFValue.add(ECFFValue[i] + (1.0 / 6.0) * (kECFF[0] + 2.0 * kECFF[1] + 2.0 * kECFF[2] + kECFF[3]) * input.step)
            ECFSValue.add(ECFSValue[i] + (1.0 / 6.0) * (kECFS[0] + 2.0 * kECFS[1] + 2.0 * kECFS[2] + kECFS[3]) * input.step)
            ThermValue.add(ThermValue[i] + (1.0 / 6.0) * (kTherm[0] + 2.0 * kTherm[1] + 2.0 * kTherm[2] + kTherm[3]) * input.step)
        }

        val advance = (1 / input.step).toInt()
        val G = arrayListOf<Double>()
        val P = arrayListOf<Double>()
        val Fa = arrayListOf<Double>()
        val ECF_F = arrayListOf<Double>()
        val ECF_S = arrayListOf<Double>()
        val Psig = arrayListOf<Double>()
        val L_diet = arrayListOf<Double>()
        val Therm = arrayListOf<Double>()
        val BW = arrayListOf<Double>()
        val SM = arrayListOf<Double>()
        val FFMHyp = arrayListOf<Double>()
        val dT = arrayListOf<Double>()

        for (i in 0..input.predictionDays) {

            val j = i * advance
            G.add(GValue[j])
            P.add(PValue[j])
            Fa.add(FaValue[j])
            ECF_F.add(ECFFValue[j])
            ECF_S.add(ECFSValue[j])
            Psig.add(PsigValue[j])
            L_diet.add(LdietValue[j])
            Therm.add(ThermValue[j])
            BW.add(
                Fa[i] + userDependantParameters.BM + userDependantParameters.ECP + userDependantParameters.ECF_b + ECF_F[i] + ECF_S[i] +
                        userDependantParameters.KICW + P[i] * (1 + Hp) + G[i] * (1 + Hg) + userDependantParameters.ICS
            )
            SM.add(
                0.45 * (P[i] + userDependantParameters.ECP) + 0.8 * G[i] + 0.35 * userDependantParameters.ICS +
                        input.gender.CoeffGender * (userDependantParameters.KICW + P[i] * Hp + G[i] * Hg + userDependantParameters.ECF_b + ECF_F[i] + ECF_S[i])
            )
            FFMHyp.add(BW[i] - Fa[i])
            dT.add(i.toDouble())

        }


        val SMHyp = arrayListOf<Double>()
        val BWHyp = arrayListOf<Double>()

        val Amax = 145 * 1.101977
        val consts = 48.33
        val amax = (Amax - 145) * input.hypLevel + 145
        val x = arrayListOf<Double>()

        for (i in 1..input.predictionDays) {

            var xValue = amax / 145 - (amax / 145 - 1) * Math.exp(-(consts) / (amax - 1) * (dT[i] / 7))
            xValue = (xValue - 0.04) / 0.96
            x.add(xValue)

            SMHyp.add(SM[i] * xValue)
            BWHyp.add(BW[i] + (SMHyp[i] - SM[i]))
            FFMHyp.add(BW[i] - Fa[i] + (SMHyp[i] - SM[i]))
        }

        return Rk4onePhaseHypHall2015Result(
            G = G.toDoubleArray(),
            P = P.toDoubleArray(),
            Fa = Fa.toDoubleArray(),
            ECF_F = ECF_F.toDoubleArray(),
            ECF_S = ECF_S.toDoubleArray(),
            Psig = Psig.toDoubleArray(),
            L_diet = L_diet.toDoubleArray(),
            Therm = Therm.toDoubleArray(),
            BW = BW.toDoubleArray(),
            SM = SM.toDoubleArray(),
            FFMHyp = FFMHyp.toDoubleArray(),
            DT = dT.toDoubleArray(),
            SMHyp = SMHyp.toDoubleArray(),
            BWHyp = BWHyp.toDoubleArray(),
            initializationHall2015ODESResult = userDependantParameters
        )
    }
}

data class Rk4onePhaseHypHall2015Input(
    val gender: Gender,
    val heightCMS: Double,
    val initialFMKg: Double,
    val initialBWKg: Double,
    val pal: Double,
    val rmr: Double,
    val propC: Double,
    val propF: Double,
    val propP: Double,
    val prop_pi: Double,
    val newC: Double,
    val newF: Double,
    val newP: Double,
    val psi: Double,
    val S_F: Double,
    val ExtraExCal: Double,
    val hypLevel: Double,
    val age: Int,
    val predictionDays: Int,
    val step: Double
)

data class Rk4onePhaseHypHall2015Result(
    val G: DoubleArray,
    val P: DoubleArray,
    val Fa: DoubleArray,
    val ECF_F: DoubleArray,
    val ECF_S: DoubleArray,
    val BW: DoubleArray,
    val SM: DoubleArray,
    val DT: DoubleArray,
    val SMHyp: DoubleArray,
    val BWHyp: DoubleArray,
    val FFMHyp: DoubleArray,
    val Psig: DoubleArray,
    val L_diet: DoubleArray,
    val Therm: DoubleArray,
    val initializationHall2015ODESResult: InitializationHall2015ODESResult
) {
    /**
     * This is recommened for data classes that
     * have arrays in them
     */
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false

        other as Rk4onePhaseHypHall2015Result

        if (!G.contentEquals(other.G)) return false
        if (!P.contentEquals(other.P)) return false
        if (!Fa.contentEquals(other.Fa)) return false
        if (!ECF_F.contentEquals(other.ECF_F)) return false
        if (!ECF_S.contentEquals(other.ECF_S)) return false
        if (!BW.contentEquals(other.BW)) return false
        if (!SM.contentEquals(other.SM)) return false
        if (!DT.contentEquals(other.DT)) return false
        if (!SMHyp.contentEquals(other.SMHyp)) return false
        if (!BWHyp.contentEquals(other.BWHyp)) return false
        if (!FFMHyp.contentEquals(other.FFMHyp)) return false
        if (!Psig.contentEquals(other.Psig)) return false
        if (!L_diet.contentEquals(other.L_diet)) return false
        if (!Therm.contentEquals(other.Therm)) return false

        return true
    }

    override fun hashCode(): Int {
        var result = G.contentHashCode()
        result = 31 * result + P.contentHashCode()
        result = 31 * result + Fa.contentHashCode()
        result = 31 * result + ECF_F.contentHashCode()
        result = 31 * result + ECF_S.contentHashCode()
        result = 31 * result + BW.contentHashCode()
        result = 31 * result + SM.contentHashCode()
        result = 31 * result + DT.contentHashCode()
        result = 31 * result + SMHyp.contentHashCode()
        result = 31 * result + BWHyp.contentHashCode()
        result = 31 * result + FFMHyp.contentHashCode()
        result = 31 * result + Psig.contentHashCode()
        result = 31 * result + L_diet.contentHashCode()
        result = 31 * result + Therm.contentHashCode()
        return result
    }
}