package net.fitmylife.core.extensions

fun Double.decimals(number: Int): Double = String.format("%.${number}f", this).toDouble()