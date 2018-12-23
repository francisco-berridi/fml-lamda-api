package formulas

import (
	"enums/Gender"
	"formulas/analytics/params"
	"math"
)

func Hall2015ODESonly(gender Gender.Enum, userDependparamaters InitializationHall2015ODESResult, S_F, G, Psig, P, L_diet, Fa, ECF_F, ECF_S, Therm float64) (dG, dPsig, dP, dL_diet, dFa, dECF_F, dECF_S, dTherm float64) {

	// Compute current body composition
	var CM = userDependparamaters.KICW + P*(1+params.Hp) + G*(1+params.Hg) + userDependparamaters.ICS
	var FFM = userDependparamaters.BM + userDependparamaters.ECP + userDependparamaters.ECF_b + ECF_F + ECF_S + CM
	var BW = Fa + FFM

	// DNL ... Equation (24)
	var DNL = userDependparamaters.CI * math.Pow(G/userDependparamaters.Ginit, params.D) / (math.Pow(G/userDependparamaters.Ginit, params.D) + math.Pow(params.K_DNL, params.D))

	// D_F ... Equation (11)
	var Lipol = math.Pow(Fa/params.F_Keys, params.Sl) * (L_diet + userDependparamaters.L_PA)
	var D_F = params.D_F_hat * Lipol

	// KTG ... Equation (15)
	var Ketogen = D_F * (3 * params.M_FA / params.M_F) * (params.A_K * ((L_diet + userDependparamaters.L_PA) / (params.K_K + L_diet + userDependparamaters.L_PA)) * math.Exp(-params.K_G*G/userDependparamaters.Ginit) * math.Exp(-params.K_P*userDependparamaters.PI/userDependparamaters.PI_b))
	var KTG = params.Rho_K * Ketogen

	// KU_excr ... Equation (16)
	var Kurine float64 = 0
	if Ketogen >= params.KTGthresh {
		Kurine = (params.KUmax*Ketogen - params.KUmax*params.KTGthresh) / (params.KTGmax - params.KTGthresh)
	}
	var KU_excr = params.Rho_K * Kurine

	// KetOx ... Equation (17)
	var KetOx = KTG - KU_excr

	// D_P ... Equation (18)
	var Proteol = P/params.P_Keys + params.Chi*(userDependparamaters.PI-userDependparamaters.PI_b)/userDependparamaters.PI_b
	var D_P = params.D_P_hat * Proteol

	// GNGf ... Equation (22)
	var GNGf = userDependparamaters.FI*(params.M_G*params.Rho_C)/(params.Rho_F*params.M_F) + D_F*params.Rho_C*(params.M_G/params.M_F)

	// GNGp ... Equation (23)
	// GNGp_glycogen_effect = max(c(0,1-Glyc_GNGeffect*tanh(G/G_Keys-1)))   ## Since Glyc_GNGeffect=0 --> GNGp_glycogen_effect=1
	// This term appears in the code but not in the papers
	var GNGp = params.GNGp_hat*(Proteol-params.Gamma_C*userDependparamaters.DCI/userDependparamaters.CI_b+(params.Gamma_P+params.Chi)*userDependparamaters.DPI/userDependparamaters.PI_b) - params.S_GNG*(GNGf-userDependparamaters.GNGf_init)
	if GNGp < 0 {
		GNGp = 0
	}

	// D_G ... Equation (19)
	var D_G = params.D_G_hat * G / params.G_Keys

	// Fat Free Mass Metabolic rate ... Equation (10)
	var gamma_FFM = params.Gamma_FFM_hat * (1 + (1-params.Sigma)*Therm)

	// PAE, Physical Activity Expenditure ... Equation (7)
	var PAE = userDependparamaters.Activ_tsb*(1+params.Sigma*Therm)*BW + userDependparamaters.Nu*BW

	// Carb, Fat and Protein weights ... Equation (27)
	var f_C = userDependparamaters.W_G*math.Pow(D_G/params.D_G_hat, params.Sg) + math.Max(0, userDependparamaters.W_C*(1+userDependparamaters.S_C*userDependparamaters.DCI/userDependparamaters.CI_b))*G/(params.Gmin+G)
	var f_F = (userDependparamaters.W_F*(D_F/params.D_F_hat) + S_F*userDependparamaters.DFI/userDependparamaters.FI_b)
	var f_P = math.Max(0, params.W_P*(1+Psig)) + Proteol*((userDependparamaters.S_A-params.Activ_min)*math.Exp(-userDependparamaters.K_A*(userDependparamaters.Activ_tsb+userDependparamaters.Nu)/(userDependparamaters.Activ_b+userDependparamaters.Nu_b))+params.Activ_min)
	var Z = f_C + f_F + f_P // normalize fractions so they add up to one: nomalizedf_P = f_P/Z, nomalizedf_C = f_c/Z, nomalizedf_F = f_F/ Z
	f_C = f_C / Z
	f_F = f_F / Z
	f_P = f_P / Z

	// The following temrs (K_1, K_2, K_3, K_4 and K_5) are defined and used to obtain formulas
	// for TEE and G3P (see PDF with algebra derivations and Clay's white paper)
	var factor_eta_N = params.Eta_N / (params.Rho_P * 6.25)
	var factor_Ms = params.Rho_C * params.M_G / params.M_F
	var K_1 = userDependparamaters.TEF + PAE + userDependparamaters.Ec + params.Gamma_B*params.M_B +
		gamma_FFM*((FFM-params.M_B)-(ECF_F+ECF_S)-(G-params.G_Keys)*(1+params.Hg)) + params.Gamma_F*Fa + (1-params.Epsilon_d)*DNL + (1-params.Epsilon_g)*(GNGp+GNGf) +
		params.Pi_K*(1-params.Epsilon_K)*KTG + (params.Eta_P+params.Epsilon_P)*D_P + params.Eta_F*D_F + params.Eta_G*D_G
	var K_2 = -(1-params.Epsilon_d)*DNL - (1-params.Epsilon_K)*KTG - KetOx - GNGf - GNGp
	var K_3 = 1 + (params.Eta_P/params.Rho_P-factor_eta_N)*f_P + (params.Eta_G/params.Rho_C)*f_C + (params.Eta_F/(params.Rho_F-factor_Ms))*f_F
	var K_4 = K_1 + factor_eta_N*(GNGp+f_P*K_2) + (params.Eta_P/params.Rho_P)*(userDependparamaters.PI-GNGp-f_P*K_2) +
		(params.Eta_F/(params.Rho_F-factor_Ms))*((1-factor_Ms/params.Rho_F)*userDependparamaters.FI+params.Epsilon_d*DNL-KU_excr-(1-params.Epsilon_K)*KTG-KetOx-f_F*K_2) +
		(params.Eta_G/params.Rho_C)*(userDependparamaters.CI-DNL-f_C*K_2)
	var K_5 = (1-factor_Ms/params.Rho_F)*userDependparamaters.FI + params.Epsilon_d*DNL - (1-params.Epsilon_K)*KTG - KU_excr - KetOx - f_F*(K_2+K_4/K_3)

	// Equation in clay's white paper (last line in page 17)
	var G3P = (factor_Ms*D_F + factor_Ms*K_5/(params.Rho_F-factor_Ms)) / (1 + f_F*factor_Ms/((params.Rho_F-factor_Ms)*K_3))

	// Equation (1.49) in clay's white paper
	var TEE = (K_4 + factor_eta_N*(f_P*G3P) -
		(params.Eta_G/params.Rho_C)*(f_C*G3P) -
		(params.Eta_F/(params.Rho_F-factor_Ms))*(f_F*G3P) -
		(params.Eta_P/params.Rho_P)*(f_P*G3P)) / K_3

	// TEE_hat, Remaining Energy Expenditure ... Equation (26)
	// Equation in clay's white paper (last line in page 15)
	var TEE_hat = TEE + K_2 + G3P

	// Oxidation Rates ... Equation (25)
	var CarbOx = GNGf + GNGp - G3P + f_C*TEE_hat
	var FatOx = f_F * TEE_hat
	var ProtOx = f_P * TEE_hat

	// #1 ODE
	dFa = ((1-params.Rho_C*params.M_G/(params.Rho_F*params.M_F))*userDependparamaters.FI + params.Epsilon_d*DNL - KU_excr -
		(1-params.Epsilon_K)*KTG - KetOx - FatOx) / (params.Rho_F - params.Rho_C*(params.M_G/params.M_F))

	// #2 ODE
	dP = (userDependparamaters.PI - GNGp - ProtOx) / params.Rho_P

	// #3 ODE
	dG = (userDependparamaters.CI - DNL + GNGp + GNGf - G3P - CarbOx) / params.Rho_C

	// #4 ODE
	dECF_F = (1 / params.Na_conc) * (userDependparamaters.Na - params.Na_b - params.Xi_Na*(ECF_F+ECF_S) - params.Xi_CI*(1-userDependparamaters.CI/userDependparamaters.CI_b))

	// #5 ODE
	dECF_S = (gender.Xi_BW()*(BW-userDependparamaters.BWinit) - ECF_S) / params.Tau_BW

	// #6 ODE
	var L_diet_target = 1 + (math.Pow(params.K_L, params.S_L))*((params.Lipol_max-params.Lipol_min)*math.Exp(-1.0*params.K_L2*userDependparamaters.CI/userDependparamaters.CI_b)+
		params.Lipol_min-1)/(math.Pow(params.K_L, params.S_L)+math.Max(0, math.Pow(userDependparamaters.Finit/params.F_Keys-1, params.S_L)))
	dL_diet = (L_diet_target - L_diet) / params.Tau_L

	// #7 ODE
	dTherm = (userDependparamaters.Therm_const*(userDependparamaters.EI-userDependparamaters.EI_b)/userDependparamaters.EI_b - Therm) / params.Tau_T

	// #8 ODE
	dPsig = (userDependparamaters.Psig_const*userDependparamaters.DPI/userDependparamaters.PI_b - Psig) / params.Tau_PI

	//return dG, dPsig, dP, dL_diet, dFa, dECF_F, dECF_S, dTherm
	return

}
