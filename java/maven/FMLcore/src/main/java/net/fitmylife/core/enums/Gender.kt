package net.fitmylife.core.enums

enum class Gender(
        val Xi_BW: Double,
        val BodyFatFactor1: Double,
        val BodyFatFactor2: Double,
        val RMRFactor: Int,
        val CoeffGender: Double
) {
    Female(0.17, 39.96, 102.01, -161, 0.413),
    Male(0.19, 37.31, 103.94, 5, 0.563)
}