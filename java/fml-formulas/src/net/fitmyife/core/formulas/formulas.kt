package net.fitmyife.core.formulas

fun randomInt(min:Int, max:Int): Int = (Math.random()*(max - min) + min).toInt()