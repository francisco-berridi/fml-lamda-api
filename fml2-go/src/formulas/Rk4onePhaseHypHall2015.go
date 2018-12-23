package formulas

import (
	"enums/Gender"
	"fmt"
	"formulas/analytics/params"
	"math"
)

type Rk4onePhaseHypHall2015Result struct {
	G                    []float64
	P                    []float64
	Fa                   []float64
	ECF_F                []float64
	ECF_S                []float64
	BW                   []float64
	SM                   []float64
	DT                   []float64
	SMHyp                []float64
	BWHyp                []float64
	FFMHyp               []float64
	Psig                 []float64
	L_diet               []float64
	Therm                []float64
	UserDependparamaters InitializationHall2015ODESResult
}

func Rk4onePhaseHypHall2015(gender Gender.Enum, heightCMS, initialFMkg, initialBWkg, pal, rmr, propC, propF, propP, prop_pi, newC, newF, newP, psi, S_F, ExtraExCal, hypLevel float64, age, days int64, step float64) Rk4onePhaseHypHall2015Result {

	/*if (1/step)%1 != 0 {
		// 1/step must be an integer, this means that step can only be: 0.1, 0.2, 0.25, 0.5
		panic(fmt.Sprintf("step not 0: %d", step)) //TODO: Do not panic
	}*/

	if step != 0.2 && step != 0.25 && step != 0.5 {
		panic(fmt.Sprintf("step not 0: %d", step))
	}

	userDependparamaters, _ := InitializationHall2015ODES(gender, age, heightCMS, pal, rmr, initialFMkg, initialBWkg, propC, propF, propP, prop_pi, newC, newF, newP, psi, ExtraExCal)

	GValue := []float64{userDependparamaters.Ginit}
	PsigValue := []float64{0.0}
	PValue := []float64{userDependparamaters.Pinit}
	LdietValue := []float64{1.0}
	FaValue := []float64{userDependparamaters.Finit}
	ECFFValue := []float64{0.0}
	ECFSValue := []float64{0.0}
	ThermValue := []float64{0.0}


	cycles := int64(float64(days) / step)

	var kG [4]float64
	var kPsig [4]float64
	var kP [4]float64
	var kLdiet [4]float64
	var kFa [4]float64
	var kECFF [4]float64
	var kECFS [4]float64
	var kTherm [4]float64

	for i := int64(0); i < cycles; i++ {

		//dG, dPsig, dP, dL_diet, dFa, dECF_F, dECF_S, dTherm
		kG[0], kPsig[0], kP[0], kLdiet[0], kFa[0], kECFF[0], kECFS[0], kTherm[0] = Hall2015ODESonly(gender, userDependparamaters, S_F, GValue[i],
			PsigValue[i], PValue[i], LdietValue[i], FaValue[i], ECFFValue[i], ECFSValue[i], ThermValue[i])

		for k := 1; k < 4; k++ {

			h := 0.5
			// half
			if k == 3 {
				// the last one doest gets 0.5, it a 1.0
				h = 1.0
			}

			kG[k], kPsig[k], kP[k], kLdiet[k], kFa[k], kECFF[k], kECFS[k], kTherm[k] = Hall2015ODESonly(gender, userDependparamaters, S_F, GValue[i]+h*float64(step)*kG[k-1],
				PsigValue[i]+h*float64(step)*kPsig[k-1], PValue[i]+h*float64(step)*kP[k-1], LdietValue[i]+h*float64(step)*kLdiet[k-1],
				FaValue[i]+h*float64(step)*kFa[k-1], ECFFValue[i]+h*float64(step)*kECFF[k-1], ECFSValue[i]+h*float64(step)*kECFS[k-1], ThermValue[i]+h*float64(step)*kTherm[k-1])
		}

		GValue = append(GValue, GValue[i]+(1.0/6.0)*(kG[0]+2.0*kG[1]+2.0*kG[2]+kG[3])*float64(step))
		PsigValue = append(PsigValue, PsigValue[i]+(1.0/6.0)*(kPsig[0]+2.0*kPsig[1]+2.0*kPsig[2]+kPsig[3])*float64(step))
		PValue = append(PValue, PValue[i]+(1.0/6.0)*(kP[0]+2.0*kP[1]+2.0*kP[2]+kP[3])*float64(step))
		LdietValue = append(LdietValue, LdietValue[i]+(1.0/6.0)*(kLdiet[0]+2.0*kLdiet[1]+2.0*kLdiet[2]+kLdiet[3])*float64(step))
		FaValue = append(FaValue, FaValue[i]+(1.0/6.0)*(kFa[0]+2.0*kFa[1]+2.0*kFa[2]+kFa[3])*float64(step))
		ECFFValue = append(ECFFValue, ECFFValue[i]+(1.0/6.0)*(kECFF[0]+2.0*kECFF[1]+2.0*kECFF[2]+kECFF[3])*float64(step))
		ECFSValue = append(ECFSValue, ECFSValue[i]+(1.0/6.0)*(kECFS[0]+2.0*kECFS[1]+2.0*kECFS[2]+kECFS[3])*float64(step))
		ThermValue = append(ThermValue, ThermValue[i]+(1.0/6.0)*(kTherm[0]+2.0*kTherm[1]+2.0*kTherm[2]+kTherm[3])*float64(step))
	}

	advance := 1 / step
	var G []float64
	var P []float64
	var Fa []float64
	var ECF_F []float64
	var ECF_S []float64
	var Psig []float64
	var L_diet []float64
	var Therm []float64
	var BW []float64
	var SM []float64
	var FFMHyp []float64
	var dT []float64

	coeffGender := gender.CoeffGender()

	for i := int64(0); i <= days; i++ {

		var j = i * int64(advance)
		G = append(G, GValue[j])
		P = append(P, PValue[j])
		Fa = append(Fa, FaValue[j])
		ECF_F = append(ECF_F, ECFFValue[j])
		ECF_S = append(ECF_S, ECFSValue[j])
		Psig = append(Psig, PsigValue[j])
		L_diet = append(L_diet, LdietValue[j])
		Therm = append(Therm, ThermValue[j])
		BW = append(BW, Fa[i]+userDependparamaters.BM+userDependparamaters.ECP+userDependparamaters.ECF_b+ECF_F[i]+ECF_S[i]+
			userDependparamaters.KICW+P[i]*(1+params.Hp)+G[i]*(1+params.Hg)+userDependparamaters.ICS)
		SM = append(SM, 0.45*(P[i]+userDependparamaters.ECP)+0.8*G[i]+0.35*userDependparamaters.ICS+
			coeffGender*(userDependparamaters.KICW+P[i]*params.Hp+G[i]*params.Hg+userDependparamaters.ECF_b+ECF_F[i]+ECF_S[i]))
		FFMHyp = append(FFMHyp, BW[i]-Fa[i])
		dT = append(dT, float64(i))
	}

	var SMHyp []float64
	var BWHyp []float64

	var Amax = 145 * 1.101977
	var consts = 48.33
	var amax = (Amax-145)*hypLevel + 145
	var x []float64

	for i := int64(0); i < days+1; i++ {

		xValue := amax/145 - (amax/145-1)*math.Exp(-(consts)/(amax-1)*(dT[i]/7))
		xValue = (xValue - 0.04) / 0.96
		x = append(x, xValue)

		SMHyp = append(SMHyp, SM[i]*xValue)
		BWHyp = append(BWHyp, BW[i]+(SMHyp[i]-SM[i]))
		FFMHyp = append(FFMHyp, BW[i]-Fa[i]+(SMHyp[i]-SM[i]))
	}

	return Rk4onePhaseHypHall2015Result{
		G:                    G,
		P:                    P,
		Fa:                   Fa,
		ECF_F:                ECF_F,
		ECF_S:                ECF_S,
		Psig:                 Psig,
		L_diet:               L_diet,
		Therm:                Therm,
		BW:                   BW,
		SM:                   SM,
		FFMHyp:               FFMHyp,
		DT:                   dT,
		UserDependparamaters: userDependparamaters,
	}
}
