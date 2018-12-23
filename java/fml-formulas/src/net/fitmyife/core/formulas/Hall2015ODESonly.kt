package net.fitmyife.core.formulas

import net.fitmyife.core.formulas.params.*
import net.fitmylife.core.enums.Gender

class Hall2015ODESonly {
    fun run(input: Hall2015ODESonlyInput): Hall2015ODESonlyResult {

        val S_F = input.S_F
        val G = input.G
        val Psig = input.Psig
        val P = input.P
        val L_diet = input.L_diet
        val Fa = input.Fa
        val ECF_F = input.ECF_F
        val ECF_S = input.ECF_S
        val Therm = input.Therm
        val userDependantParameters = input.userDependantParameters

        // Compute current body composition
        val CM = userDependantParameters.KICW + P*(1+Hp) + G*(1+Hg) + userDependantParameters.ICS
        val FFM = userDependantParameters.BM + userDependantParameters.ECP + userDependantParameters.ECF_b + ECF_F + ECF_S + CM
        val BW = Fa + FFM

        // DNL ... Equation (24)
        val DNL = userDependantParameters.CI * Math.pow(G/userDependantParameters.Ginit, D) / (Math.pow(G/userDependantParameters.Ginit, D) + Math.pow(K_DNL, D))

        // D_F ... Equation (11)
        val Lipol = Math.pow(Fa/F_Keys, Sl) * (L_diet + userDependantParameters.L_PA)
        val D_F = D_F_hat * Lipol

        // KTG ... Equation (15)
        val Ketogen = D_F * (3 * M_FA / M_F) * (A_K * ((L_diet + userDependantParameters.L_PA) / (K_K + L_diet + userDependantParameters.L_PA)) * Math.exp(-K_G*G/userDependantParameters.Ginit) * Math.exp(-K_P*userDependantParameters.PI/userDependantParameters.PI_b))
        val KTG = Rho_K * Ketogen

        // KU_excr ... Equation (16)
        var Kurine = 0.0
        if (Ketogen >= KTGthresh) {
            Kurine = (KUmax*Ketogen - KUmax*KTGthresh) / (KTGmax - KTGthresh)
        }
        
        val KU_excr = Rho_K * Kurine

        // KetOx ... Equation (17)
        val KetOx = KTG - KU_excr

        // D_P ... Equation (18)
        val Proteol = P/P_Keys + Chi*(userDependantParameters.PI-userDependantParameters.PI_b)/userDependantParameters.PI_b
        val D_P = D_P_hat * Proteol

        // GNGf ... Equation (22)
        val GNGf = userDependantParameters.FI*(M_G*Rho_C)/(Rho_F*M_F) + D_F*Rho_C*(M_G/M_F)

        // GNGp ... Equation (23)
        // GNGp_glycogen_effect = max(c(0,1-Glyc_GNGeffect*tanh(G/G_Keys-1)))   ## Since Glyc_GNGeffect=0 --> GNGp_glycogen_effect=1
        // This term appears in the code but not in the papers
        var GNGp = GNGp_hat*(Proteol-Gamma_C*userDependantParameters.DCI/userDependantParameters.CI_b+(Gamma_P+Chi)*userDependantParameters.DPI/userDependantParameters.PI_b) - S_GNG*(GNGf-userDependantParameters.GNGf_init)
        if (GNGp < 0) {
            GNGp = 0.0
        }

        // D_G ... Equation (19)
        val D_G = D_G_hat * G / G_Keys

        // Fat Free Mass Metabolic rate ... Equation (10)
        val gamma_FFM = Gamma_FFM_hat * (1 + (1-Sigma)*Therm)

        // PAE, Physical Activity Expenditure ... Equation (7)
        val PAE = userDependantParameters.Activ_tsb*(1+Sigma*Therm)*BW + userDependantParameters.Nu*BW

        // Carb, Fat and Protein weights ... Equation (27)
        var f_C = userDependantParameters.W_G*Math.pow(D_G/D_G_hat, Sg) + Math.max(0.0, userDependantParameters.W_C*(1+userDependantParameters.S_C*userDependantParameters.DCI/userDependantParameters.CI_b))*G/(Gmin+G)
        var f_F = (userDependantParameters.W_F*(D_F/D_F_hat) + S_F*userDependantParameters.DFI/userDependantParameters.FI_b)
        var f_P = Math.max(0.0, W_P*(1+Psig)) + Proteol*((userDependantParameters.S_A-Activ_min)*Math.exp(-userDependantParameters.K_A*(userDependantParameters.Activ_tsb+userDependantParameters.Nu)/(userDependantParameters.Activ_b+userDependantParameters.Nu_b))+Activ_min)
        val Z = f_C + f_F + f_P // normalize fractions so they add up to one: nomalizedf_P = f_P/Z, nomalizedf_C = f_c/Z, nomalizedf_F = f_F/ Z
        f_C = f_C / Z
        f_F = f_F / Z
        f_P = f_P / Z

        // The following temrs (K_1, K_2, K_3, K_4 and K_5) are defined and used to obtain formulas
        // for TEE and G3P (see PDF with algebra derivations and Clay's white paper)
        val factor_eta_N = Eta_N / (Rho_P * 6.25)
        val factor_Ms = Rho_C * M_G / M_F
        val K_1 = userDependantParameters.TEF + PAE + userDependantParameters.Ec + Gamma_B*M_B +
                gamma_FFM*((FFM-M_B)-(ECF_F+ECF_S)-(G-G_Keys)*(1+Hg)) + Gamma_F*Fa + (1-Epsilon_d)*DNL + (1-Epsilon_g)*(GNGp+GNGf) +
                Pi_K*(1-Epsilon_K)*KTG + (Eta_P+Epsilon_P)*D_P + Eta_F*D_F + Eta_G*D_G
        val K_2 = -(1-Epsilon_d)*DNL - (1-Epsilon_K)*KTG - KetOx - GNGf - GNGp
        val K_3 = 1 + (Eta_P/Rho_P-factor_eta_N)*f_P + (Eta_G/Rho_C)*f_C + (Eta_F/(Rho_F-factor_Ms))*f_F
        val K_4 = K_1 + factor_eta_N*(GNGp+f_P*K_2) + (Eta_P/Rho_P)*(userDependantParameters.PI-GNGp-f_P*K_2) +
                (Eta_F/(Rho_F-factor_Ms))*((1-factor_Ms/Rho_F)*userDependantParameters.FI+Epsilon_d*DNL-KU_excr-(1-Epsilon_K)*KTG-KetOx-f_F*K_2) +
                (Eta_G/Rho_C)*(userDependantParameters.CI-DNL-f_C*K_2)
        val K_5 = (1-factor_Ms/Rho_F)*userDependantParameters.FI + Epsilon_d*DNL - (1-Epsilon_K)*KTG - KU_excr - KetOx - f_F*(K_2+K_4/K_3)

        // Equation in clay's white paper (last line in page 17)
        val G3P = (factor_Ms*D_F + factor_Ms*K_5/(Rho_F-factor_Ms)) / (1 + f_F*factor_Ms/((Rho_F-factor_Ms)*K_3))

        // Equation (1.49) in clay's white paper
        val TEE = (K_4 + factor_eta_N*(f_P*G3P) -
                (Eta_G/Rho_C)*(f_C*G3P) -
                (Eta_F/(Rho_F-factor_Ms))*(f_F*G3P) -
                (Eta_P/Rho_P)*(f_P*G3P)) / K_3

        // TEE_hat, Remaining Energy Expenditure ... Equation (26)
        // Equation in clay's white paper (last line in page 15)
        val TEE_hat = TEE + K_2 + G3P

        // Oxidation Rates ... Equation (25)
        val CarbOx = GNGf + GNGp - G3P + f_C*TEE_hat
        val FatOx = f_F * TEE_hat
        val ProtOx = f_P * TEE_hat

        // #1 ODE
        val dFa = ((1-Rho_C*M_G/(Rho_F*M_F))*userDependantParameters.FI + Epsilon_d*DNL - KU_excr -
                (1-Epsilon_K)*KTG - KetOx - FatOx) / (Rho_F - Rho_C*(M_G/M_F))

        // #2 ODE
        val dP = (userDependantParameters.PI - GNGp - ProtOx) / Rho_P

        // #3 ODE
        val dG = (userDependantParameters.CI - DNL + GNGp + GNGf - G3P - CarbOx) / Rho_C

        // #4 ODE
        val dECF_F = (1 / Na_conc) * (userDependantParameters.Na - Na_b - Xi_Na*(ECF_F+ECF_S) - Xi_CI*(1-userDependantParameters.CI/userDependantParameters.CI_b))

        // #5 ODE
        val dECF_S = (input.gender.Xi_BW*(BW-userDependantParameters.BWinit) - ECF_S) / Tau_BW

        // #6 ODE
        val L_diet_target = 1 + (Math.pow(K_L, S_L))*((Lipol_max-Lipol_min)*Math.exp(-1.0*K_L2*userDependantParameters.CI/userDependantParameters.CI_b)+
                Lipol_min-1)/(Math.pow(K_L, S_L)+Math.max(0.0, Math.pow(userDependantParameters.Finit/F_Keys-1, S_L)))
        val dL_diet = (L_diet_target - L_diet) / Tau_L

        // #7 ODE
        val dTherm = (userDependantParameters.Therm_const*(userDependantParameters.EI-userDependantParameters.EI_b)/userDependantParameters.EI_b - Therm) / Tau_T

        // #8 ODE
        val dPsig = (userDependantParameters.Psig_const*userDependantParameters.DPI/userDependantParameters.PI_b - Psig) / Tau_PI

        return Hall2015ODESonlyResult(
            dG = dG,
            dPsig= dPsig,
            dP= dP,
            dL_diet= dL_diet,
            dFa= dFa,
            dECF_F= dECF_F,
            dECF_S= dECF_S,
            dTherm = dTherm
        )
    }
}

data class Hall2015ODESonlyInput(
    val userDependantParameters: InitializationHall2015ODESResult,
    val gender: Gender,
    val S_F: Double,
    val G: Double,
    val Psig: Double,
    val P: Double,
    val L_diet: Double,
    val Fa: Double,
    val ECF_F: Double,
    val ECF_S: Double,
    val Therm: Double
)

data class Hall2015ODESonlyResult(
    val dG: Double,
    val dPsig: Double,
    val dP: Double,
    val dL_diet: Double,
    val dFa: Double,
    val dECF_F: Double,
    val dECF_S: Double,
    val dTherm: Double
)