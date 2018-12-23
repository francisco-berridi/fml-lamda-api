package net.fitmylife.core.enums

enum class Goal {
    WeightLoss,
    GetHealthy,
    Endurance,
    Strength,
    Aesthetics,
    Circuit;

    companion object {
        fun fromParse(string: String): Goal? {
            if (string == "weightManagement") {
                return WeightLoss
            }

            for (value in values()) {
                if (string.toUpperCase().contains(value.name.toUpperCase())) {
                    return value
                }
            }

            return null
        }
    }
}