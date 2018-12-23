package formulas

import (
	"enums/Gender"
	"formulas/analytics/params"
	"math"
)

// This is what the Javascript code calls `userDependParameters`
type InitializationHall2015ODESResult struct {
	EI_b        float64
	CI_b        float64
	FI_b        float64
	PI_b        float64
	TEFinit     float64
	EI          float64
	CI          float64
	FI          float64
	PI          float64
	TEF         float64
	DCI         float64
	DFI         float64
	DPI         float64
	BWinit      float64
	Finit       float64
	Pinit       float64
	Ginit       float64
	KICW        float64
	ICS         float64
	ECF_b       float64
	BM          float64
	ECP         float64
	Activ_tsb   float64
	Activ_b     float64
	Nu_b        float64
	Nu          float64
	GNGf_init   float64
	W_G         float64
	W_C         float64
	W_F         float64
	S_C         float64
	Sp_pos      float64
	S_A         float64
	K_A         float64
	Ec          float64
	L_PA        float64
	Na          float64
	Therm_const float64
	Psig_const  float64
}

func InitializationHall2015ODES(gender Gender.Enum, age int64, heightCMS, pal, rmr, initialFMkg, initialBWkg, propC, propF, propP, prop_pi, NewC, NewF, NewP, psi, extraExCal float64) (InitializationHall2015ODESResult, error) {

	var F_imbal float64 = 0
	var P_imbal float64 = 0
	var G_imbal float64 = 0

	// Body composition
	var Ginit float64 = 500
	var Finit = 1000 * initialFMkg
	var BWinit = 1000 * initialBWkg
	var FFMinit = BWinit - Finit
	var BM = params.BMfracBW * BWinit // Bone mineral mass in grams at 4% of initial BW

	var ECF_b float64

	if gender == Gender.Male {
		ECF_b = 1000 * (0.191*initialBWkg + 9.570*heightCMS/100 + 0.025*float64(age) - 12.424)
	} else {
		ECF_b = 1000 * (0.167*initialBWkg + 5.987*heightCMS/100 - 4.027)
	}

	var ECP = 0.732*BM + 0.01087*ECF_b      // Extracellular protein, Reference: Wang AJCN 2003 ECP =1877.6
	var CMinit = FFMinit - BM - ECP - ECF_b // initial Cell Mass
	var Pinit = params.PfracCM * CMinit
	var KICW = params.ICWfracCM*CMinit - Pinit*params.Hp - Ginit*params.Hg // constant ICW amount of intracellular water mass
	var ICS = (1-params.PfracCM-params.ICWfracCM)*CMinit - Ginit           // Intracellular solids including salts, nucleic acids, minerals, etc.

	// Baseline intake
	var EI_b = pal * rmr
	var CI_b = propC * EI_b
	var FI_b = propF * EI_b
	var PI_b = propP * EI_b

	// New intake and delta intake
	var CI = NewC
	var PI = NewP
	var FI = NewF
	var EI = CI + FI + PI
	var dCI = CI - CI_b
	var dPI = PI - PI_b
	var dFI = FI - FI_b

	// Thermic Effect of Feeding ... Equation (5)
	var TEFinit = params.Alpha_C*CI_b + params.Alpha_F*FI_b + params.Alpha_P*PI_b
	var TEF = params.Alpha_C*CI + params.Alpha_F*FI + params.Alpha_P*PI

	// Activity parameters
	var activ_b = prop_pi * (EI_b - rmr - TEFinit - F_imbal - P_imbal - G_imbal) / BWinit // activ_b is defined to satisfy Equation (28)
	var nu_b = (1 - prop_pi) * (EI_b - rmr - TEFinit - F_imbal - P_imbal - G_imbal) / BWinit
	if activ_b < 0 {
		activ_b = 0
		nu_b = 0
	}
	var activ_tsb = activ_b
	var nu = nu_b + extraExCal/BWinit

	// Initial DNL ... Equation (24)
	var DNL_init = CI_b / (1 + math.Pow(params.K_DNL, params.D))

	// Initial D_F ... Equation (11)
	var Lipol_init = math.Pow(Finit/params.F_Keys, params.Sl) // Lipol_init = D_F_init/params.D_F_hat, where params.D_F_hat = params.Lipol_b*M_F
	var D_F_init = params.D_F_hat * Lipol_init                // at t=0 L_diet=1 and L_PA=0

	// initial D_P ... Equation (18)
	var Proteol_init = Pinit / params.P_Keys     // Proteol_init = D_P_init/D_P_hat
	var D_P_init = params.D_P_hat * Proteol_init // Proteol_b*AA_mass = D_P_hat

	// G3P_init ... Equation (21) Assumes D_F_init = Synth_F_init
	var G3P_init = params.Rho_C * D_F_init * params.M_G / params.M_F
	// var G3P_init = params.Rho_C*params.Lipol_b*Lipol_init*params.M_G

	//  Equation (14)
	var L_PA = psi * ((activ_tsb+nu)/(activ_b+nu_b) - 1)

	// Initial KTG ... Equation (15)
	//parameter2015.rho_K * D_F_init * (3 * parameter2015.M_FA / parameter2015.M_F) * parameter2015.A_K * (1 / (parameter2015.K_K + 1)) * (Math.exp(-parameter2015.k_G)) * Math.exp(-parameter2015.k_P);
	var KTG_init = params.Rho_K * D_F_init * (3 * params.M_FA / params.M_F) * params.A_K * (1 / (params.K_K + 1)) * (math.Exp(-1*params.K_G)) * math.Exp(-1*params.K_P)

	// Initial KetOx ... Equation (17) & assumes that KU_excretion at baseline is zero (end of page 13)
	var KetOx_init = KTG_init

	// Initial GNGp ... Equation (23)
	var GNGp_init = math.Max(0, params.GNGp_hat*Proteol_init) // Proteol_init = D_P_init/D_P_hat

	// Initial GNGf ... Equation (22)
	var GNGf_init = FI_b*params.Rho_C*params.M_G/(params.Rho_F*params.M_F) + params.Rho_C*D_F_init*(params.M_G/params.M_F)

	// Initial D_G ... Equation (19)
	var D_G_init = params.D_G_hat * Ginit / params.G_Keys // Glycogenol_b*G_mass = D_G_hat

	// Nutrient Balance Parameter Constraints ... Equation (28) ... at baseline gamma_FFM = params.gamma_FFM_hat
	var Ec = EI_b - (TEFinit + (activ_b+nu_b)*BWinit + params.Gamma_B*params.M_B + params.Gamma_FFM_hat*(FFMinit-params.M_B) + params.Gamma_F*Finit + (1-params.Epsilon_d)*DNL_init + (1-params.Epsilon_g)*(GNGp_init+GNGf_init) + params.Pi_K*(1-params.Epsilon_K)*KTG_init + (params.Eta_P+params.Epsilon_P)*D_P_init + params.Eta_F*D_F_init + params.Eta_G*D_G_init + (1+params.Eta_F/params.Rho_F)*F_imbal + (1+params.Eta_P/params.Rho_P)*P_imbal + (1+params.Eta_G/params.Rho_C)*G_imbal + params.Eta_N*(PI_b-P_imbal)/(6.25*params.Rho_P))

	// Nutrient Balance Parameter Constraints ... Equation (31) (35) in 2010 paper
	var Omega = EI_b - (1-params.Epsilon_d)*DNL_init - (1-params.Epsilon_K)*KTG_init - KetOx_init - GNGf_init - GNGp_init + G3P_init - G_imbal - F_imbal - P_imbal

	// Nutrient Balance constraints ... Equation (30) (34 in 2010 paper)
	//           ((1 - (parameter2015.rho_C * parameter2015.M_G) / (parameter2015.rho_F * parameter2015.M_F)) * FI_b + parameter2015.epsilon_d * DNL_init - (1 - parameter2015.epsilon_K) * KTG_init - KetOx_init - F_imbal) / Omega;
	var zeta_F = ((1-(params.Rho_C*params.M_G)/(params.Rho_F*params.M_F))*FI_b + params.Epsilon_d*DNL_init - (1-params.Epsilon_K)*KTG_init - KetOx_init - F_imbal) / Omega
	var zeta_C = (CI_b - DNL_init - G_imbal) / Omega
	var zeta_P = (PI_b - GNGp_init - P_imbal) / Omega

	// Macronutrient Imbalance ... Equation (33) (37 in 2010 paper)
	var w_G = (zeta_C / zeta_P) * (Pinit/params.P_Keys + params.W_P) / (1 + params.W_CG*params.G_Keys/(params.Gmin+Ginit))
	var w_F = zeta_F / (1 - zeta_F) * (1 + zeta_C/zeta_P) * (Pinit/params.P_Keys + params.W_P) * math.Pow(params.F_Keys/Finit, params.Sl)

	// zeta_F 0.26597860528660755
	// zeta_C 0.6398900272971827
	// zeta_P 0.09413136741620974
	// Pinit 8067.762095
	// Finit 27590
	// Omega 2162.7110766638543



	var w_C = w_G * params.W_CG

	//////// Physical inactivity constraints (PIC)

	// DeltaTEE from equation (59) (63 in 2010 paper)
	var DeltaTEE_PIC = -(activ_b+nu_b)*BWinit - params.Eta_N*params.N_bal

	var DeltaG_PIC = DeltaTEE_PIC / (2 * params.Rho_C)

	// Xi ... Equation (63) (67 in 2010 paper)
	var Inact_P = (PI_b - GNGp_init - 6.25*params.Rho_P*params.N_bal) / (EI_b + DeltaTEE_PIC - (1-params.Epsilon_d)*DNL_init - (2-params.Epsilon_K)*KTG_init - GNGp_init - GNGf_init + G3P_init)

	// w_hat from PIC ... Equation (61)
	var w_G_hat_PIC = w_G * (1 + DeltaG_PIC/Ginit)
	var w_C_hat_PIC = w_C * ((Ginit + DeltaG_PIC) / (params.Gmin + Ginit + DeltaG_PIC))
	var w_F_hat_PIC = w_F * (1 - psi) * math.Pow(Finit/params.F_Keys, params.Sl)

	// S_A from PIC ... Equation (62) (66 in 2010 paper)
	var S_A = math.Max(1, (params.P_Keys/Pinit)*((Inact_P/(1-Inact_P))*(w_C_hat_PIC+w_G_hat_PIC+w_F_hat_PIC)-params.W_P)) //Inact_P is capital Xi
	var k_A = math.Log((S_A - params.Activ_min) / (1 - params.Activ_min))

	//////// Carbohydrate Perturbation Constraint (CPC)

	var DeltaG_CPC = params.Kappa_C * 510 / params.Rho_C

	// Delta Thermogenesis from Carbohydrate Perturbation Constraint ... Equation (38)
	var DeltaT_CPC = params.Lamb_2 * 510 / EI_b * (1 - params.Tau_T + params.Tau_T*math.Exp(-1/params.Tau_T))

	// Delta DNL from Carbohydrate Perturbation Constraint ... Equation (39) (43 in 2010 paper)
	var DeltaDNL_CPC = (CI_b+510)*math.Pow(1+DeltaG_CPC/Ginit, params.D)/(math.Pow(params.K_DNL, params.D)+math.Pow(1+DeltaG_CPC/Ginit, params.D)) - DNL_init

	// Delta GNG from Carbohydrate Perturbation Constraint ... Equation (40) (44 in 2010 paper)
	var DeltaGNG_CPC = -params.Gamma_C*(510/CI_b)*params.GNGp_hat + (1-params.S_GNG)*params.Rho_C*(params.M_G/params.M_F)*params.D_F_hat*math.Pow(Finit/params.F_Keys, params.Sl)*(math.Pow(params.K_L, params.S_L))*((params.Lipol_max-params.Lipol_min)*math.Exp(-params.K_L2*(1+510/CI_b))+params.Lipol_min-1)/(math.Pow(params.K_L, params.S_L)+math.Max(0, math.Pow(Finit/params.F_Keys-1, params.S_L)))

	// Delta RMR from Carbohydrate Perturbation Constraint ... Equation (37)
	var DeltaRMR_CPC = params.Kappa_C*params.Eta_G*510/params.Rho_C + params.Eta_G*params.D_G_hat*DeltaG_CPC/params.G_Keys + params.Gamma_FFM_hat*(FFMinit-params.M_B)*(1-params.Sigma)*DeltaT_CPC + (1-params.Epsilon_d)*DeltaDNL_CPC + (1-params.Epsilon_g)*DeltaGNG_CPC

	// DeltaTEE from Equation (34 to 37) (38 to 41 in 2010 paper)
	var DeltaTEE_CPC = params.Alpha_C*510 + activ_b*BWinit*params.Sigma*DeltaT_CPC + DeltaRMR_CPC

	// TEE_hat values for CPC ... Equation (26)
	var TEE_hat_CPC = EI_b + DeltaTEE_CPC - (1-params.Epsilon_d)*(DNL_init+DeltaDNL_CPC) - (2-params.Epsilon_K)*KTG_init - params.GNGp_hat - params.Lipol_b*params.M_G*params.Rho_C - (params.M_G/(params.Rho_F*params.M_F))*FI_b*params.Rho_C - DeltaGNG_CPC + G3P_init

	// Theta from CPC ... Equation (44)
	var Theta_CPC = (CI_b + (1-params.Kappa_C)*510 - (DNL_init + DeltaDNL_CPC)) / TEE_hat_CPC

	// w_hat values for CPC ... Equation (46)
	var w_F_hat_CPC = w_F * math.Pow(Finit/params.F_Keys, params.Sl) * (1 + math.Pow(params.K_L, params.S_L)*((params.Lipol_max-params.Lipol_min)*math.Exp(-1.0 * params.K_L2*(1+510/CI_b))+params.Lipol_min-1)/(math.Pow(params.K_L, params.S_L)+math.Max(0, math.Pow(Finit/params.F_Keys-1, params.S_L))))
	var w_C_hat_CPC = w_C * (Ginit + DeltaG_CPC) / (params.Gmin + Ginit + DeltaG_CPC)
	var w_G_hat_CPC = w_G * (Ginit + DeltaG_CPC) / params.G_Keys

	// S_C from CPC ... Equation (45)
	var S_C = (CI_b / 510) * ((Theta_CPC*(Pinit/params.P_Keys+params.W_P+w_F_hat_CPC+w_G_hat_CPC)-w_F_hat_CPC)/((1-Theta_CPC)*w_C_hat_CPC) - 1) // Updated to 2015
	if S_C < 0 {
		S_C = 0
	}

	// w_F_hat_CPC
	// Theta_CPC
	var Sp_pos = 2.5

	// Constant in Therm ODE
	var Therm_const float64
	if EI < EI_b {
		Therm_const = params.Lamb_1
	} else {
		Therm_const = params.Lamb_2
	}

	//  Constant in Psig ODE
	var Psig_const float64
	if PI < PI_b {
		Psig_const = params.Sp_neg
	} else {
		Psig_const = Sp_pos
	}

	return InitializationHall2015ODESResult{
		EI_b:        EI_b,
		CI_b:        CI_b,
		FI_b:        FI_b,
		PI_b:        PI_b,
		TEFinit:     TEFinit,
		EI:          EI,
		CI:          CI,
		FI:          FI,
		PI:          PI,
		TEF:         TEF,
		DCI:         dCI,
		DFI:         dFI,
		DPI:         dPI,
		BWinit:      BWinit,
		Finit:       Finit,
		Pinit:       Pinit,
		Ginit:       Ginit,
		KICW:        KICW,
		ICS:         ICS,
		ECF_b:       ECF_b,
		BM:          BM,
		ECP:         ECP,
		Activ_tsb:   activ_tsb,
		Activ_b:     activ_b,
		Nu_b:        nu_b,
		Nu:          nu,
		GNGf_init:   GNGf_init,
		W_G:         w_G,
		W_C:         w_C,
		W_F:         w_F,
		S_C:         S_C,
		Sp_pos:      Sp_pos,
		S_A:         S_A,
		K_A:         k_A,
		Ec:          Ec,
		L_PA:        L_PA,
		Na:          params.Na_b,
		Therm_const: Therm_const,
		Psig_const:  Psig_const,
	}, nil
}
