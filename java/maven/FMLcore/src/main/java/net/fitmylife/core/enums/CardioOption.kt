package net.fitmylife.core.enums

enum class CardioOption {
    RowingMachine,
    Elliptical,
    Treadmill,
    Stairs,
    Running,
    Walking,
    IndoorSpinCycle,
    Swimming,
    StairMaster,
    IndoorSpinBike;

    companion object {
        fun fromParse(string: String): CardioOption? {
            for (value in CardioOption.values()) {
                if (value.name.toUpperCase() == string.toUpperCase()) {
                    return value;
                }
            }

            return null
        }
    }
}