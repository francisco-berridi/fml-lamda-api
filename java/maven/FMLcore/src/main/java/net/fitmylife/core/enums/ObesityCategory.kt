package net.fitmylife.core.enums

enum class ObesityCategory {
    UnderWeight,
    NormalWeight,
    OverWeight,
    Class1Obesity,
    Class2Obesity,
    Class3Obesity;

    companion object {
        fun getWithBmi(bmi: Double): ObesityCategory {
            return when (bmi) {
                in 0.0..18.50 -> UnderWeight
                in 18.51..25.0 -> NormalWeight
                in 25.0..30.0 -> OverWeight
                in 30.0..35.0 -> Class1Obesity
                in 35.0..40.0 -> Class2Obesity
                else -> Class3Obesity
            }
        }

        fun fromParse(string: String): ObesityCategory? {
            for (value in values()) {
                if (value.name.toUpperCase().contains(string.toUpperCase())) {
                    return value
                }
            }

            return null
        }
    }
}