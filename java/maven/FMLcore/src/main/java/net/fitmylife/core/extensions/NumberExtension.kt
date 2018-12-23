package net.fitmylife.core.extensions

val Number.inch: Double get() = inchesToCms(this)
val Number.cm: Double get() = this.toDouble()
val Number.lb: Double get() = lbsToKgs(this)
val Number.kg: Double get() = this.toDouble()
val Number.feet: Double get() = feetToCms(this)
fun Number.feet(inches: Number): Double = this.feet  + inches.inch


fun inchesToCms(any:Number): Double {
    return any.toDouble() * 2.54
}

fun lbsToKgs(any:Number): Double {
    return any.toDouble() * 0.453592
}

fun feetToCms(any: Number): Double {
    return any.inch * 12
}