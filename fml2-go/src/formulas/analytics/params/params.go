package params

import "math"

const (
	Alpha_F              = 0.025
	Alpha_C              = 0.075
	Alpha_P              = 0.25
	Rho_C                = 3.7
	Rho_F                = 9.4
	Rho_P                = 4.7
	Rho_K                = 4.45
	Digest_C             = 0.95
	Digest_F             = 0.96
	Digest_P             = 0.90
	PfracCM              = 0.25
	ICWfracCM            = 0.7
	BMfracBW             = 0.04
	Hp                   = 1.6
	Hg                   = 2.7
	Na_conc              = 3.22
	Na_b         float64 = 4000
	Xi_CI        float64 = 4000
	Xi_Na        float64 = 3
	Tau_BW       float64 = 200
	ATP_kcal     float64 = 19
	ATP_TG_synth float64 = 8
	AA_mass      float64 = 110
	ATP_G_synth  float64 = 2
	ATP_Nexcr    float64 = 4
	G_mass       float64 = 180
	M_G          float64 = 92
	M_F          float64 = 860
	M_FA                 = 274.595745 // Defined by M_FA = (rho_F*M_F - rho_C*M_G)/(3*rho_F) as Nancy observed on Oct 162017
	N_mass       float64 = 14
	Eta_P                = 1.34
	Epsilon_P            = 0.54
	Epsilon_d            = 0.95
	Epsilon_g            = 0.8
	Epsilon_K            = 0.77
	Pi_K                 = 0.9
	Sl                   = 0.333333333333333333
	K_L          float64 = 4
	S_L          float64 = 2
	Tau_L                = 3.32192809489
	Lipol_min            = 0.87
	Lipol_max    float64 = 3
	Lipol_b              = 0.16
	Psi                  = 0.4
	K_DNL        float64 = 2
	D            float64 = 4
	A_K                  = 0.80
	K_K          float64 = 1
	KTGthresh    float64 = 70
	KUmax        float64 = 20
	KTGmax       float64 = 400
	M_B          float64 = 1400
	Slope_liver          = 0.01736
	Slope_skmusc         = 0.5934
	Slope_kid            = 0.003786
	Slope_hrt            = 0.00288
	Slope_res            = 0.3747
	Gamma_B              = 0.24
	Gamma_liver          = 0.2
	Gamma_skmusc         = 0.013
	Gamma_kid            = 0.44
	Gamma_hrt            = 0.44
	Gamma_res            = 0.012
	Gamma_F              = 0.0045
	Tau_T        float64 = 7
	Tau_PI               = 1.1
	Sg           float64 = 1
	Gmin         float64 = 10
	Chi          float64 = 0
	D_F_hat      float64 = 150
	GNGp_hat     float64 = 300
	S_GNG                = 0.5
	Gamma_C              = 0.46
	Gamma_P              = 0.31
	Glycogenol_b float64 = 1
	Kappa_C              = 0.6
	N_bal        float64 = -2
	Activ_min            = 0
	BW_Keys      float64 = 67533
	F_Keys       float64 = 9050
	G_Keys       float64 = 500
	Lamb_1               = 0.8
	Lamb_2               = 0.1
	Sigma                = 0.6
	W_P                  = 1.2
	W_CG         float64 = 1
	Sp_neg               = 0.85
	Proteol_b            = 2.73
	K_P                  = math.Ln2
	K_G                  = math.Ln2

	// Calculated
	// -------------------------------------------------------------------------------
	D_G_hat       = Glycogenol_b * G_mass
	D_P_hat       = Proteol_b * AA_mass
	FFM_Keys      = BW_Keys - F_Keys
	BM_Keys       = BMfracBW * BW_Keys
	ECF_Keys      = 0.7 * 0.235 * BW_Keys
	ECP_Keys      = 0.732*BM_Keys + 0.01087*ECF_Keys
	P_Keys        = PfracCM * (FFM_Keys - BM_Keys - ECP_Keys - ECF_Keys)
	Gamma_FFM_hat = Gamma_liver*Slope_liver + Gamma_skmusc*Slope_skmusc + Gamma_kid*Slope_kid + Gamma_hrt*Slope_hrt + Gamma_res*Slope_res
	Eta_F         = ATP_TG_synth * ATP_kcal / M_F   // Fat deposition cost in kcal/g
	Eta_G         = ATP_G_synth * ATP_kcal / G_mass // Glycogen deposition cost in kcal/g
	Eta_N         = ATP_Nexcr * ATP_kcal / N_mass
	K_L2 float64  = 2.796342808247888 // result of math.Log((Lipol_max - Lipol_min) / (1 - Lipol_min))  // Equation (13)
)