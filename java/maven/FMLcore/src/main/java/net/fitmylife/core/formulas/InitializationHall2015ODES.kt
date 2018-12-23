package net.fitmylife.core.formulas

import net.fitmylife.core.enums.Gender
import net.fitmylife.core.formulas.params.*

class InitializationHall2015ODES {

    fun run(input: InitializationHall2015ODESInput): InitializationHall2015ODESResult {

        val F_imbal  = 0
        val P_imbal  = 0
        val G_imbal  = 0

        // Body composition
        val Ginit  = 500.0
        val Finit = 1000 * input.initialFMkg
        val BWinit = 1000 * input.initialBWkg
        val FFMinit = BWinit - Finit
        val BM = BMfracBW * BWinit // Bone mineral mass in grams at 4% of initial BW

        val ECF_b = when (input.gender) {
            Gender.Male -> 1000 * (0.191*input.initialBWkg + 9.570*input.heightCMS/100 + 0.025*input.age - 12.424)
            Gender.Female -> 1000 * (0.167*input.initialBWkg + 5.987*input.heightCMS/100 - 4.027)
        }


        val ECP = 0.732*BM + 0.01087*ECF_b      // Extracellular protein, Reference: Wang AJCN 2003 ECP =1877.6
        val CMinit = FFMinit - BM - ECP - ECF_b // initial Cell Mass
        val Pinit = PfracCM * CMinit
        val KICW = ICWfracCM *CMinit - Pinit* Hp - Ginit* Hg // constant ICW amount of intracellular water mass
        val ICS = (1- PfracCM - ICWfracCM)*CMinit - Ginit           // Intracellular solids including salts, nucleic acids, minerals, etc.

        // Baseline intake
        val EI_b = input.pal * input.rmr
        val CI_b = input.propC * EI_b
        val FI_b = input.propF * EI_b
        val PI_b = input.propP * EI_b

        // New intake and delta intake
        val CI = input.NewC
        val PI = input.NewP
        val FI = input.NewF
        val EI = CI + FI + PI
        val dCI = CI - CI_b
        val dPI = PI - PI_b
        val dFI = FI - FI_b

        // Thermic Effect of Feeding ... Equation (5)
        val TEFinit = Alpha_C *CI_b + Alpha_F *FI_b + Alpha_P *PI_b
        val TEF = Alpha_C *CI + Alpha_F *FI + Alpha_P *PI

        // Activity parameters
        var activ_b = input.prop_pi * (EI_b - input.rmr - TEFinit - F_imbal - P_imbal - G_imbal) / BWinit // activ_b is defined to satisfy Equation (28)
        var nu_b = (1 - input.prop_pi) * (EI_b - input.rmr - TEFinit - F_imbal - P_imbal - G_imbal) / BWinit
        if (activ_b < 0) {
            activ_b = 0.0
            nu_b = 0.0
        }
        val activ_tsb = activ_b
        val nu = nu_b + input.extraExCal/BWinit

        // Initial DNL ... Equation (24)
        val DNL_init = CI_b / (1 + Math.pow(K_DNL, D))

        // Initial D_F ... Equation (11)
        val Lipol_init = Math.pow(Finit/ F_Keys, Sl) // Lipol_init = D_F_init/D_F_hat, where D_F_hat = Lipol_b*M_F
        val D_F_init = D_F_hat * Lipol_init                // at t=0 L_diet=1 and L_PA=0

        // initial D_P ... Equation (18)
        val Proteol_init = Pinit / P_Keys     // Proteol_init = D_P_init/D_P_hat
        val D_P_init = D_P_hat * Proteol_init // Proteol_b*AA_mass = D_P_hat

        // G3P_init ... Equation (21) Assumes D_F_init = Synth_F_init
        val G3P_init = Rho_C * D_F_init * M_G / M_F
        // val G3P_init = Rho_C*Lipol_b*Lipol_init*M_G

        //  Equation (14)
        val L_PA = input.psi * ((activ_tsb+nu)/(activ_b+nu_b) - 1)

        // Initial KTG ... Equation (15)
        //parameter2015.rho_K * D_F_init * (3 * parameter2015.M_FA / parameter2015.M_F) * parameter2015.A_K * (1 / (parameter2015.K_K + 1)) * (Math.exp(-parameter2015.k_G)) * Math.exp(-parameter2015.k_P);
        val KTG_init = Rho_K * D_F_init * (3 * M_FA / M_F) * A_K * (1 / (K_K + 1)) * (Math.exp(-1* K_G)) * Math.exp(-1* K_P)

        // Initial KetOx ... Equation (17) & assumes that KU_excretion at baseline is zero (end of page 13)
        val KetOx_init = KTG_init

        // Initial GNGp ... Equation (23)
        val GNGp_init = Math.max(0.0, GNGp_hat *Proteol_init) // Proteol_init = D_P_init/D_P_hat

        // Initial GNGf ... Equation (22)
        val GNGf_init = FI_b* Rho_C * M_G /(Rho_F * M_F) + Rho_C *D_F_init*(M_G / M_F)

        // Initial D_G ... Equation (19)
        val D_G_init = D_G_hat * Ginit / G_Keys // Glycogenol_b*G_mass = D_G_hat

        // Nutrient Balance Parameter Constraints ... Equation (28) ... at baseline gamma_FFM = gamma_FFM_hat
        val Ec = EI_b - (TEFinit + (activ_b+nu_b)*BWinit + Gamma_B * M_B + Gamma_FFM_hat *(FFMinit- M_B) + Gamma_F *Finit + (1- Epsilon_d)*DNL_init + (1- Epsilon_g)*(GNGp_init+GNGf_init) + Pi_K *(1- Epsilon_K)*KTG_init + (Eta_P + Epsilon_P)*D_P_init + Eta_F *D_F_init + Eta_G *D_G_init + (1+ Eta_F / Rho_F)*F_imbal + (1+ Eta_P / Rho_P)*P_imbal + (1+ Eta_G / Rho_C)*G_imbal + Eta_N *(PI_b-P_imbal)/(6.25* Rho_P))

        // Nutrient Balance Parameter Constraints ... Equation (31) (35) in 2010 paper
        val Omega = EI_b - (1- Epsilon_d)*DNL_init - (1- Epsilon_K)*KTG_init - KetOx_init - GNGf_init - GNGp_init + G3P_init - G_imbal - F_imbal - P_imbal

        // Nutrient Balance constraints ... Equation (30) (34 in 2010 paper)
        //           ((1 - (parameter2015.rho_C * parameter2015.M_G) / (parameter2015.rho_F * parameter2015.M_F)) * FI_b + parameter2015.epsilon_d * DNL_init - (1 - parameter2015.epsilon_K) * KTG_init - KetOx_init - F_imbal) / Omega;
        val zeta_F = ((1-(Rho_C * M_G)/(Rho_F * M_F))*FI_b + Epsilon_d *DNL_init - (1- Epsilon_K)*KTG_init - KetOx_init - F_imbal) / Omega
        val zeta_C = (CI_b - DNL_init - G_imbal) / Omega
        val zeta_P = (PI_b - GNGp_init - P_imbal) / Omega

        // Macronutrient Imbalance ... Equation (33) (37 in 2010 paper)
        val w_G = (zeta_C / zeta_P) * (Pinit/ P_Keys + W_P) / (1 + W_CG * G_Keys /(Gmin +Ginit))
        val w_F = zeta_F / (1 - zeta_F) * (1 + zeta_C/zeta_P) * (Pinit/ P_Keys + W_P) * Math.pow(F_Keys /Finit, Sl)

        // zeta_F 0.26597860528660755
        // zeta_C 0.6398900272971827
        // zeta_P 0.09413136741620974
        // Pinit 8067.762095
        // Finit 27590
        // Omega 2162.7110766638543

        val w_C = w_G * W_CG

        //////// Physical inactivity constraints (PIC)

        // DeltaTEE from equation (59) (63 in 2010 paper)
        val DeltaTEE_PIC = -(activ_b+nu_b)*BWinit - Eta_N * N_bal

        val DeltaG_PIC = DeltaTEE_PIC / (2 * Rho_C)

        // Xi ... Equation (63) (67 in 2010 paper)
        val Inact_P = (PI_b - GNGp_init - 6.25* Rho_P * N_bal) / (EI_b + DeltaTEE_PIC - (1- Epsilon_d)*DNL_init - (2- Epsilon_K)*KTG_init - GNGp_init - GNGf_init + G3P_init)

        // w_hat from PIC ... Equation (61)
        val w_G_hat_PIC = w_G * (1 + DeltaG_PIC/Ginit)
        val w_C_hat_PIC = w_C * ((Ginit + DeltaG_PIC) / (Gmin + Ginit + DeltaG_PIC))
        val w_F_hat_PIC = w_F * (1 - input.psi) * Math.pow(Finit/ F_Keys, Sl)

        // S_A from PIC ... Equation (62) (66 in 2010 paper)
        val S_A = Math.max(1.0, (P_Keys /Pinit)*((Inact_P/(1-Inact_P))*(w_C_hat_PIC+w_G_hat_PIC+w_F_hat_PIC)- W_P)) //Inact_P is capital Xi
        val k_A = Math.log((S_A - Activ_min) / (1.0 - Activ_min))

        //////// Carbohydrate Perturbation Constraint (CPC)

        val DeltaG_CPC = Kappa_C * 510 / Rho_C

        // Delta Thermogenesis from Carbohydrate Perturbation Constraint ... Equation (38)
        val DeltaT_CPC = Lamb_2 * 510 / EI_b * (1 - Tau_T + Tau_T *Math.exp(-1/ Tau_T))

        // Delta DNL from Carbohydrate Perturbation Constraint ... Equation (39) (43 in 2010 paper)
        val DeltaDNL_CPC = (CI_b+510)*Math.pow(1+DeltaG_CPC/Ginit, D)/(Math.pow(K_DNL, D)+Math.pow(1+DeltaG_CPC/Ginit, D)) - DNL_init

        // Delta GNG from Carbohydrate Perturbation Constraint ... Equation (40) (44 in 2010 paper)
        val DeltaGNG_CPC = -Gamma_C *(510/CI_b)* GNGp_hat + (1- S_GNG)* Rho_C *(M_G / M_F)* D_F_hat *Math.pow(Finit/ F_Keys, Sl)*(Math.pow(K_L, S_L))*((Lipol_max - Lipol_min)*Math.exp(-K_L2 *(1+510/CI_b))+ Lipol_min -1)/(Math.pow(K_L, S_L)+Math.max(0.0, Math.pow(Finit/ F_Keys -1, S_L)))

        // Delta calculateRMR from Carbohydrate Perturbation Constraint ... Equation (37)
        val DeltaRMR_CPC = Kappa_C * Eta_G *510/ Rho_C + Eta_G * D_G_hat *DeltaG_CPC/ G_Keys + Gamma_FFM_hat *(FFMinit- M_B)*(1- Sigma)*DeltaT_CPC + (1- Epsilon_d)*DeltaDNL_CPC + (1- Epsilon_g)*DeltaGNG_CPC

        // DeltaTEE from Equation (34 to 37) (38 to 41 in 2010 paper)
        val DeltaTEE_CPC = Alpha_C *510 + activ_b*BWinit* Sigma *DeltaT_CPC + DeltaRMR_CPC

        // TEE_hat values for CPC ... Equation (26)
        val TEE_hat_CPC = EI_b + DeltaTEE_CPC - (1- Epsilon_d)*(DNL_init+DeltaDNL_CPC) - (2- Epsilon_K)*KTG_init - GNGp_hat - Lipol_b * M_G * Rho_C - (M_G /(Rho_F * M_F))*FI_b* Rho_C - DeltaGNG_CPC + G3P_init

        // Theta from CPC ... Equation (44)
        val Theta_CPC = (CI_b + (1- Kappa_C)*510 - (DNL_init + DeltaDNL_CPC)) / TEE_hat_CPC

        // w_hat values for CPC ... Equation (46)
        val w_F_hat_CPC = w_F * Math.pow(Finit/ F_Keys, Sl) * (1 + Math.pow(K_L, S_L)*((Lipol_max - Lipol_min)*Math.exp(-1.0 * K_L2 *(1+510/CI_b))+ Lipol_min -1)/(Math.pow(K_L, S_L)+Math.max(0.0, Math.pow(Finit/ F_Keys -1, S_L))))
        val w_C_hat_CPC = w_C * (Ginit + DeltaG_CPC) / (Gmin + Ginit + DeltaG_CPC)
        val w_G_hat_CPC = w_G * (Ginit + DeltaG_CPC) / G_Keys

        // S_C from CPC ... Equation (45)
        var S_C = (CI_b / 510) * ((Theta_CPC*(Pinit/ P_Keys + W_P +w_F_hat_CPC+w_G_hat_CPC)-w_F_hat_CPC)/((1-Theta_CPC)*w_C_hat_CPC) - 1) // Updated to 2015
        if (S_C < 0) {
            S_C = 0.0
        }

        // w_F_hat_CPC
        // Theta_CPC
        val Sp_pos = 2.5

        // Constant in Therm ODE
        val Therm_const = if (EI < EI_b) {
            Lamb_1
        } else {
            Lamb_2
        }

        //  Constant in Psig ODE
        val Psig_const = if (PI < PI_b) {
            Sp_neg
        } else {
            Sp_pos
        }

        return InitializationHall2015ODESResult(
            EI_b=       EI_b,
            CI_b=        CI_b,
            FI_b=        FI_b,
            PI_b=        PI_b,
            TEFinit=     TEFinit,
            EI=          EI,
            CI=          CI,
            FI=          FI,
            PI=          PI,
            TEF=         TEF,
            DCI=         dCI,
            DFI=         dFI,
            DPI=         dPI,
            BWinit=     BWinit,
            Finit=       Finit,
            Pinit=       Pinit,
            Ginit=       Ginit,
            KICW=        KICW,
            ICS=         ICS,
            ECF_b=       ECF_b,
            BM=          BM,
            ECP=         ECP,
            Activ_tsb=   activ_tsb,
            Activ_b=     activ_b,
            Nu_b=        nu_b,
            Nu=          nu,
            GNGf_init=   GNGf_init,
            W_G=         w_G,
            W_C=         w_C,
            W_F=         w_F,
            S_C=         S_C,
            Sp_pos=      Sp_pos,
            S_A=         S_A,
            K_A=         k_A,
            Ec=          Ec,
            L_PA=        L_PA,
            Na=          Na_b,
            Therm_const= Therm_const,
            Psig_const=  Psig_const
        )
    }
}

data class InitializationHall2015ODESInput(
    val gender: Gender,
    val age: Int,
    val heightCMS: Double,
    val pal: Double,
    val rmr: Double,
    val initialFMkg: Double,
    val initialBWkg: Double,
    val propC: Double,
    val propF: Double,
    val propP: Double,
    val prop_pi: Double,
    val NewC: Double,
    val NewF: Double,
    val NewP: Double,
    val psi: Double,
    val extraExCal: Double
)


data class InitializationHall2015ODESResult(
    val EI_b: Double,
    val CI_b: Double,
    val FI_b: Double,
    val PI_b: Double,
    val TEFinit: Double,
    val EI: Double,
    val CI: Double,
    val FI: Double,
    val PI: Double,
    val TEF: Double,
    val DCI: Double,
    val DFI: Double,
    val DPI: Double,
    val BWinit: Double,
    val Finit: Double,
    val Pinit: Double,
    val Ginit: Double,
    val KICW: Double,
    val ICS: Double,
    val ECF_b: Double,
    val BM: Double,
    val ECP: Double,
    val Activ_tsb: Double,
    val Activ_b: Double,
    val Nu_b: Double,
    val Nu: Double,
    val GNGf_init: Double,
    val W_G: Double,
    val W_C: Double,
    val W_F: Double,
    val S_C: Double,
    val Sp_pos: Double,
    val S_A: Double,
    val K_A: Double,
    val Ec: Double,
    val L_PA: Double,
    val Na: Double,
    val Therm_const: Double,
    val Psig_const: Double
)