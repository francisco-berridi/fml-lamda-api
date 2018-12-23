package formulas

import (
	"enums/Gender"
	"formulas/analytics/params"
)

// This function continues the projection from t=lastTime to t=dates
// rk4ContinueTrajectoryHall2015 reproduces the second part (phase two) of rk4TwoPhasesHypHall2015, without running the first phase
// the information from phase one is encapsulated in userDependparamaters and lastG, ..., lastTherm
// due to time dependencies, hypLevel is assumed to be zero
func Rk4ContinueTrajectoryHall2015(NewC2, NewF2, NewP2, ExtraExCal2, lastG, lastPsig, lastP, lastLdiet, lastFa, lastECFF,
	lastECFS, lastTherm, psi, S_F, lastTime, dates, step float64, userDependparamaters InitializationHall2015ODESResult, gender Gender.Enum) Rk4onePhaseHypHall2015Result {

	var GValue = []float64{lastG}
	var PsigValue = []float64{lastPsig}
	var PValue = []float64{lastP}
	var LdietValue = []float64{lastLdiet}
	var FaValue = []float64{lastFa}
	var ECFFValue = []float64{lastECFF}
	var ECFSValue = []float64{lastECFS}
	var ThermValue = []float64{lastTherm}

	var kG [4]float64
	var kPsig [4]float64
	var kP [4]float64
	var kLdiet [4]float64
	var kFa [4]float64
	var kECFF [4]float64
	var kECFS [4]float64
	var kTherm [4]float64

	// Update userDependparamaters for phase 2
	userDependparamaters.CI = NewC2
	userDependparamaters.FI = NewF2
	userDependparamaters.PI = NewP2
	userDependparamaters.EI = userDependparamaters.CI + userDependparamaters.FI + userDependparamaters.PI
	userDependparamaters.DCI = userDependparamaters.CI - userDependparamaters.CI_b
	userDependparamaters.DFI = userDependparamaters.FI - userDependparamaters.FI_b
	userDependparamaters.DPI = userDependparamaters.PI - userDependparamaters.PI_b
	userDependparamaters.TEF = params.Alpha_F*userDependparamaters.FI + params.Alpha_C*userDependparamaters.CI + params.Alpha_P*userDependparamaters.PI
	userDependparamaters.Nu = userDependparamaters.Nu_b + ExtraExCal2/userDependparamaters.BWinit
	userDependparamaters.L_PA = psi * ((userDependparamaters.Activ_tsb+userDependparamaters.Nu)/(userDependparamaters.Activ_b+userDependparamaters.Nu_b) - 1)

	if userDependparamaters.EI < userDependparamaters.EI_b {
		userDependparamaters.Therm_const = params.Lamb_1
	} else {
		userDependparamaters.Therm_const = params.Lamb_2
	}

	if userDependparamaters.PI-userDependparamaters.PI_b < 0 {
		userDependparamaters.Psig_const = params.Sp_neg

	} else {
		userDependparamaters.Psig_const = userDependparamaters.Sp_pos
	}

	// This is phase TWO
	for i := 0; i < int((dates-lastTime)/step); i++ {

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

	advance := int(1 / step)
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
	var dT []float64

	coeffGender := gender.CoeffGender()

	for i := 0; i <= int(dates-lastTime); i++ {

		var j = i * advance
		G = append(G, GValue[j])
		P = append(P, PValue[j])
		Fa = append(Fa, FaValue[j])
		ECF_F = append(ECF_F, ECFFValue[j])
		ECF_S = append(ECF_S, ECFSValue[j])
		Psig = append(Psig, PsigValue[j])
		L_diet = append(L_diet, LdietValue[j])
		Therm = append(Therm, ThermValue[j])
		BW = append(BW, Fa[i]+userDependparamaters.BM+userDependparamaters.ECP+userDependparamaters.ECF_b+ECF_F[i]+ECF_S[i]+
			userDependparamaters.KICW+P[i]*(1.0+params.Hp)+G[i]*(1.0+params.Hg)+userDependparamaters.ICS)
		SM = append(SM, 0.45*(P[i]+userDependparamaters.ECP)+0.8*G[i]+0.35*userDependparamaters.ICS+
			coeffGender*(userDependparamaters.KICW+P[i]*params.Hp+G[i]*params.Hg+userDependparamaters.ECF_b+ECF_F[i]+ECF_S[i]))
		dT = append(dT, float64(i) + lastTime)
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
		DT:                   dT,
		UserDependparamaters: userDependparamaters,
	}

}
