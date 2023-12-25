/*
 * Match any PE file as defined by MZ and PE signatures at required locations.
 */

rule IsPeFile {
	strings:
		$mz = "MZ"

	condition:
		$mz at 0 and uint32(uint32(0x3C)) == 0x4550
}
