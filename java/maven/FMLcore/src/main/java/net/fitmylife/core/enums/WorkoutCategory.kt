package net.fitmylife.core.enums

enum class WorkoutCategory {
    Rest,
    Cardio,
    Strength,
    Circuit,
    FullBodyStrength,
    FullBodyCircuit,
    LowerStrength,
    UpperStrength,
    PullStrength,
    PushStrength,
    UpperCircuit,
    LowerCircuit,
    Other,
    OtherStrength;


    companion object {
        fun fromParse(string: String): WorkoutCategory {
            for (value in WorkoutCategory.values()) {
                if (value.name.toUpperCase() == string.toUpperCase()) {
                    return value
                }
            }

            return Other
        }
    }
}