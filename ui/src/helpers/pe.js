// Convert hex to an ascii string.
function hex2a(hexx) {
  var hex = hexx.toString(); // force conversion
  var str = '';
  for (var i = 0;
    (i < hex.length && hex.substr(i, 2) !== '00'); i += 2)
    str += String.fromCharCode(parseInt(hex.substr(i, 2), 16));
  return str;
}

// Reverse will reverse a string.
function reverse(s) {
  return s.split("").reverse().join("");
}

// Convert a dec number to a hex.
function dec2Hex(n) {
  return n.toString(16)
}

// Convert a dec number to a hex string.
function dec2HexString(n) {
  return '0x' + n.toString(16)
}

// prodIdToStr maps product ids to MS internal names.
// list from: https://github.com/kirschju/richheader
function prodIdToStr(prodId) {
  switch (prodId) {
    case 0x0000:
      return "Unknown"
    case 0x0001:
      return "Import0"
    case 0x0002:
      return "Linker510"
    case 0x0003:
      return "Cvtomf510"
    case 0x0004:
      return "Linker600"
    case 0x0005:
      return "Cvtomf600"
    case 0x0006:
      return "Cvtres500"
    case 0x0007:
      return "Utc11_Basic"
    case 0x0008:
      return "Utc11_C"
    case 0x0009:
      return "Utc12_Basic"
    case 0x000a:
      return "Utc12_C"
    case 0x000b:
      return "Utc12_CPP"
    case 0x000c:
      return "AliasObj60"
    case 0x000d:
      return "VisualBasic60"
    case 0x000e:
      return "Masm613"
    case 0x000f:
      return "Masm710"
    case 0x0010:
      return "Linker511"
    case 0x0011:
      return "Cvtomf511"
    case 0x0012:
      return "Masm614"
    case 0x0013:
      return "Linker512"
    case 0x0014:
      return "Cvtomf512"
    case 0x0015:
      return "Utc12_C_Std"
    case 0x0016:
      return "Utc12_CPP_Std"
    case 0x0017:
      return "Utc12_C_Book"
    case 0x0018:
      return "Utc12_CPP_Book"
    case 0x0019:
      return "Implib700"
    case 0x001a:
      return "Cvtomf700"
    case 0x001b:
      return "Utc13_Basic"
    case 0x001c:
      return "Utc13_C"
    case 0x001d:
      return "Utc13_CPP"
    case 0x001e:
      return "Linker610"
    case 0x001f:
      return "Cvtomf610"
    case 0x0020:
      return "Linker601"
    case 0x0021:
      return "Cvtomf601"
    case 0x0022:
      return "Utc12_1_Basic"
    case 0x0023:
      return "Utc12_1_C"
    case 0x0024:
      return "Utc12_1_CPP"
    case 0x0025:
      return "Linker620"
    case 0x0026:
      return "Cvtomf620"
    case 0x0027:
      return "AliasObj70"
    case 0x0028:
      return "Linker621"
    case 0x0029:
      return "Cvtomf621"
    case 0x002a:
      return "Masm615"
    case 0x002b:
      return "Utc13_LTCG_C"
    case 0x002c:
      return "Utc13_LTCG_CPP"
    case 0x002d:
      return "Masm620"
    case 0x002e:
      return "ILAsm100"
    case 0x002f:
      return "Utc12_2_Basic"
    case 0x0030:
      return "Utc12_2_C"
    case 0x0031:
      return "Utc12_2_CPP"
    case 0x0032:
      return "Utc12_2_C_Std"
    case 0x0033:
      return "Utc12_2_CPP_Std"
    case 0x0034:
      return "Utc12_2_C_Book"
    case 0x0035:
      return "Utc12_2_CPP_Book"
    case 0x0036:
      return "Implib622"
    case 0x0037:
      return "Cvtomf622"
    case 0x0038:
      return "Cvtres501"
    case 0x0039:
      return "Utc13_C_Std"
    case 0x003a:
      return "Utc13_CPP_Std"
    case 0x003b:
      return "Cvtpgd1300"
    case 0x003c:
      return "Linker622"
    case 0x003d:
      return "Linker700"
    case 0x003e:
      return "Export622"
    case 0x003f:
      return "Export700"
    case 0x0040:
      return "Masm700"
    case 0x0041:
      return "Utc13_POGO_I_C"
    case 0x0042:
      return "Utc13_POGO_I_CPP"
    case 0x0043:
      return "Utc13_POGO_O_C"
    case 0x0044:
      return "Utc13_POGO_O_CPP"
    case 0x0045:
      return "Cvtres700"
    case 0x0046:
      return "Cvtres710p"
    case 0x0047:
      return "Linker710p"
    case 0x0048:
      return "Cvtomf710p"
    case 0x0049:
      return "Export710p"
    case 0x004a:
      return "Implib710p"
    case 0x004b:
      return "Masm710p"
    case 0x004c:
      return "Utc1310p_C"
    case 0x004d:
      return "Utc1310p_CPP"
    case 0x004e:
      return "Utc1310p_C_Std"
    case 0x004f:
      return "Utc1310p_CPP_Std"
    case 0x0050:
      return "Utc1310p_LTCG_C"
    case 0x0051:
      return "Utc1310p_LTCG_CPP"
    case 0x0052:
      return "Utc1310p_POGO_I_C"
    case 0x0053:
      return "Utc1310p_POGO_I_CPP"
    case 0x0054:
      return "Utc1310p_POGO_O_C"
    case 0x0055:
      return "Utc1310p_POGO_O_CPP"
    case 0x0056:
      return "Linker624"
    case 0x0057:
      return "Cvtomf624"
    case 0x0058:
      return "Export624"
    case 0x0059:
      return "Implib624"
    case 0x005a:
      return "Linker710"
    case 0x005b:
      return "Cvtomf710"
    case 0x005c:
      return "Export710"
    case 0x005d:
      return "Implib710"
    case 0x005e:
      return "Cvtres710"
    case 0x005f:
      return "Utc1310_C"
    case 0x0060:
      return "Utc1310_CPP"
    case 0x0061:
      return "Utc1310_C_Std"
    case 0x0062:
      return "Utc1310_CPP_Std"
    case 0x0063:
      return "Utc1310_LTCG_C"
    case 0x0064:
      return "Utc1310_LTCG_CPP"
    case 0x0065:
      return "Utc1310_POGO_I_C"
    case 0x0066:
      return "Utc1310_POGO_I_CPP"
    case 0x0067:
      return "Utc1310_POGO_O_C"
    case 0x0068:
      return "Utc1310_POGO_O_CPP"
    case 0x0069:
      return "AliasObj710"
    case 0x006a:
      return "AliasObj710p"
    case 0x006b:
      return "Cvtpgd1310"
    case 0x006c:
      return "Cvtpgd1310p"
    case 0x006d:
      return "Utc1400_C"
    case 0x006e:
      return "Utc1400_CPP"
    case 0x006f:
      return "Utc1400_C_Std"
    case 0x0070:
      return "Utc1400_CPP_Std"
    case 0x0071:
      return "Utc1400_LTCG_C"
    case 0x0072:
      return "Utc1400_LTCG_CPP"
    case 0x0073:
      return "Utc1400_POGO_I_C"
    case 0x0074:
      return "Utc1400_POGO_I_CPP"
    case 0x0075:
      return "Utc1400_POGO_O_C"
    case 0x0076:
      return "Utc1400_POGO_O_CPP"
    case 0x0077:
      return "Cvtpgd1400"
    case 0x0078:
      return "Linker800"
    case 0x0079:
      return "Cvtomf800"
    case 0x007a:
      return "Export800"
    case 0x007b:
      return "Implib800"
    case 0x007c:
      return "Cvtres800"
    case 0x007d:
      return "Masm800"
    case 0x007e:
      return "AliasObj800"
    case 0x007f:
      return "PhoenixPrerelease"
    case 0x0080:
      return "Utc1400_CVTCIL_C"
    case 0x0081:
      return "Utc1400_CVTCIL_CPP"
    case 0x0082:
      return "Utc1400_LTCG_MSIL"
    case 0x0083:
      return "Utc1500_C"
    case 0x0084:
      return "Utc1500_CPP"
    case 0x0085:
      return "Utc1500_C_Std"
    case 0x0086:
      return "Utc1500_CPP_Std"
    case 0x0087:
      return "Utc1500_CVTCIL_C"
    case 0x0088:
      return "Utc1500_CVTCIL_CPP"
    case 0x0089:
      return "Utc1500_LTCG_C"
    case 0x008a:
      return "Utc1500_LTCG_CPP"
    case 0x008b:
      return "Utc1500_LTCG_MSIL"
    case 0x008c:
      return "Utc1500_POGO_I_C"
    case 0x008d:
      return "Utc1500_POGO_I_CPP"
    case 0x008e:
      return "Utc1500_POGO_O_C"
    case 0x008f:
      return "Utc1500_POGO_O_CPP"
    case 0x0090:
      return "Cvtpgd1500"
    case 0x0091:
      return "Linker900"
    case 0x0092:
      return "Export900"
    case 0x0093:
      return "Implib900"
    case 0x0094:
      return "Cvtres900"
    case 0x0095:
      return "Masm900"
    case 0x0096:
      return "AliasObj900"
    case 0x0097:
      return "Resource"
    case 0x0098:
      return "AliasObj1000"
    case 0x0099:
      return "Cvtpgd1600"
    case 0x009a:
      return "Cvtres1000"
    case 0x009b:
      return "Export1000"
    case 0x009c:
      return "Implib1000"
    case 0x009d:
      return "Linker1000"
    case 0x009e:
      return "Masm1000"
    case 0x009f:
      return "Phx1600_C"
    case 0x00a0:
      return "Phx1600_CPP"
    case 0x00a1:
      return "Phx1600_CVTCIL_C"
    case 0x00a2:
      return "Phx1600_CVTCIL_CPP"
    case 0x00a3:
      return "Phx1600_LTCG_C"
    case 0x00a4:
      return "Phx1600_LTCG_CPP"
    case 0x00a5:
      return "Phx1600_LTCG_MSIL"
    case 0x00a6:
      return "Phx1600_POGO_I_C"
    case 0x00a7:
      return "Phx1600_POGO_I_CPP"
    case 0x00a8:
      return "Phx1600_POGO_O_C"
    case 0x00a9:
      return "Phx1600_POGO_O_CPP"
    case 0x00aa:
      return "Utc1600_C"
    case 0x00ab:
      return "Utc1600_CPP"
    case 0x00ac:
      return "Utc1600_CVTCIL_C"
    case 0x00ad:
      return "Utc1600_CVTCIL_CPP"
    case 0x00ae:
      return "Utc1600_LTCG_C"
    case 0x00af:
      return "Utc1600_LTCG_CPP"
    case 0x00b0:
      return "Utc1600_LTCG_MSIL"
    case 0x00b1:
      return "Utc1600_POGO_I_C"
    case 0x00b2:
      return "Utc1600_POGO_I_CPP"
    case 0x00b3:
      return "Utc1600_POGO_O_C"
    case 0x00b4:
      return "Utc1600_POGO_O_CPP"
    case 0x00b5:
      return "AliasObj1010"
    case 0x00b6:
      return "Cvtpgd1610"
    case 0x00b7:
      return "Cvtres1010"
    case 0x00b8:
      return "Export1010"
    case 0x00b9:
      return "Implib1010"
    case 0x00ba:
      return "Linker1010"
    case 0x00bb:
      return "Masm1010"
    case 0x00bc:
      return "Utc1610_C"
    case 0x00bd:
      return "Utc1610_CPP"
    case 0x00be:
      return "Utc1610_CVTCIL_C"
    case 0x00bf:
      return "Utc1610_CVTCIL_CPP"
    case 0x00c0:
      return "Utc1610_LTCG_C"
    case 0x00c1:
      return "Utc1610_LTCG_CPP"
    case 0x00c2:
      return "Utc1610_LTCG_MSIL"
    case 0x00c3:
      return "Utc1610_POGO_I_C"
    case 0x00c4:
      return "Utc1610_POGO_I_CPP"
    case 0x00c5:
      return "Utc1610_POGO_O_C"
    case 0x00c6:
      return "Utc1610_POGO_O_CPP"
    case 0x00c7:
      return "AliasObj1100"
    case 0x00c8:
      return "Cvtpgd1700"
    case 0x00c9:
      return "Cvtres1100"
    case 0x00ca:
      return "Export1100"
    case 0x00cb:
      return "Implib1100"
    case 0x00cc:
      return "Linker1100"
    case 0x00cd:
      return "Masm1100"
    case 0x00ce:
      return "Utc1700_C"
    case 0x00cf:
      return "Utc1700_CPP"
    case 0x00d0:
      return "Utc1700_CVTCIL_C"
    case 0x00d1:
      return "Utc1700_CVTCIL_CPP"
    case 0x00d2:
      return "Utc1700_LTCG_C"
    case 0x00d3:
      return "Utc1700_LTCG_CPP"
    case 0x00d4:
      return "Utc1700_LTCG_MSIL"
    case 0x00d5:
      return "Utc1700_POGO_I_C"
    case 0x00d6:
      return "Utc1700_POGO_I_CPP"
    case 0x00d7:
      return "Utc1700_POGO_O_C"
    case 0x00d8:
      return "Utc1700_POGO_O_CPP"
    case 0x00d9:
      return "AliasObj1200"
    case 0x00da:
      return "Cvtpgd1800"
    case 0x00db:
      return "Cvtres1200"
    case 0x00dc:
      return "Export1200"
    case 0x00dd:
      return "Implib1200"
    case 0x00de:
      return "Linker1200"
    case 0x00df:
      return "Masm1200"
    case 0x00e0:
      return "Utc1800_C"
    case 0x00e1:
      return "Utc1800_CPP"
    case 0x00e2:
      return "Utc1800_CVTCIL_C"
    case 0x00e3:
      return "Utc1800_CVTCIL_CPP"
    case 0x00e4:
      return "Utc1800_LTCG_C"
    case 0x00e5:
      return "Utc1800_LTCG_CPP"
    case 0x00e6:
      return "Utc1800_LTCG_MSIL"
    case 0x00e7:
      return "Utc1800_POGO_I_C"
    case 0x00e8:
      return "Utc1800_POGO_I_CPP"
    case 0x00e9:
      return "Utc1800_POGO_O_C"
    case 0x00ea:
      return "Utc1800_POGO_O_CPP"
    case 0x00eb:
      return "AliasObj1210"
    case 0x00ec:
      return "Cvtpgd1810"
    case 0x00ed:
      return "Cvtres1210"
    case 0x00ee:
      return "Export1210"
    case 0x00ef:
      return "Implib1210"
    case 0x00f0:
      return "Linker1210"
    case 0x00f1:
      return "Masm1210"
    case 0x00f2:
      return "Utc1810_C"
    case 0x00f3:
      return "Utc1810_CPP"
    case 0x00f4:
      return "Utc1810_CVTCIL_C"
    case 0x00f5:
      return "Utc1810_CVTCIL_CPP"
    case 0x00f6:
      return "Utc1810_LTCG_C"
    case 0x00f7:
      return "Utc1810_LTCG_CPP"
    case 0x00f8:
      return "Utc1810_LTCG_MSIL"
    case 0x00f9:
      return "Utc1810_POGO_I_C"
    case 0x00fa:
      return "Utc1810_POGO_I_CPP"
    case 0x00fb:
      return "Utc1810_POGO_O_C"
    case 0x00fc:
      return "Utc1810_POGO_O_CPP"
    case 0x00fd:
      return "AliasObj1400"
    case 0x00fe:
      return "Cvtpgd1900"
    case 0x00ff:
      return "Cvtres1400"
    case 0x0100:
      return "Export1400"
    case 0x0101:
      return "Implib1400"
    case 0x0102:
      return "Linker1400"
    case 0x0103:
      return "Masm1400"
    case 0x0104:
      return "Utc1900_C"
    case 0x0105:
      return "Utc1900_CPP"
    case 0x0106:
      return "Utc1900_CVTCIL_C"
    case 0x0107:
      return "Utc1900_CVTCIL_CPP"
    case 0x0108:
      return "Utc1900_LTCG_C"
    case 0x0109:
      return "Utc1900_LTCG_CPP"
    case 0x010a:
      return "Utc1900_LTCG_MSIL"
    case 0x010b:
      return ": 'Utc1900_POGO_I_C"
    case 0x010c:
      return "Utc1900_POGO_I_CPP"
    case 0x010d:
      return "Utc1900_POGO_O_C"
    case 0x010e:
      return "Utc1900_POGO_O_CPP"
  }
  return "?"
}

// prodIdToVsVersion retrieves the Visual Studio version from product id.
// list from: https://github.com/kirschju/richheader
function prodIdToVsVersion(prodId) {
  if (prodId > 0x010e || prodId < 0)
    return ""
  if (prodId >= 0x00fd && prodId < (0x010e + 1))
    return "Visual Studio 2015 14.00"
  if (prodId >= 0x00eb && prodId < 0x00fd)
    return "Visual Studio 2013 12.10"
  if (prodId >= 0x00d9 && prodId < 0x00eb)
    return "Visual Studio 2013 12.00"
  if (prodId >= 0x00c7 && prodId < 0x00d9)
    return "Visual Studio 2012 11.00"
  if (prodId >= 0x00b5 && prodId < 0x00c7)
    return "Visual Studio 2010 10.10"
  if (prodId >= 0x0098 && prodId < 0x00b5)
    return "Visual Studio 2010 10.00"
  if (prodId >= 0x0083 && prodId < 0x0098)
    return "Visual Studio 2008 09.00"
  if (prodId >= 0x006d && prodId < 0x0083)
    return "Visual Studio 2005 08.00"
  if (prodId >= 0x005a && prodId < 0x006d)
    return "Visual Studio 2003 07.10"
  if (prodId === 1)
    return "Visual Studio"
  return ""
}

function machine2String(machine) {
  switch (machine) {
    case 0x0:
      return "Unknown";
    case 0x1d3:
      return "Matsushita AM33";
    case 0x8664:
      return "x64";
    case 0x1c0:
      return "ARM little endian";
    case 0xaa64:
      return "ARM64 little endian";
    case 0x1c4:
      return "ARM Thumb-2 little endian";
    case 0xebc:
      return "EFI byte code";
    case 0x14c:
      return "Intel 386 or later processors and compatible processors";
    case 0x200:
      return "Intel Itanium processor family";
    case 0x9041:
      return "Mitsubishi M32R little endian";
    case 0x266:
      return "MIPS16";
    case 0x366:
      return "MIPS with FPU";
    case 0x466:
      return "MIPS16 with FPU";
    case 0x1f0:
      return "Power PC little endian";
    case 0x1f1:
      return "Power PC with floating point support";
    case 0x166:
      return "MIPS little endian";
    case 0x5032:
      return "RISC-V 32-bit address space";
    case 0x5064:
      return "RISC-V 64-bit address space";
    case 0x5128:
      return "RISC-V 128-bit address space";
    case 0x1a2:
      return "Hitachi SH3";
    case 0x1a3:
      return "Hitachi SH3 DSP";
    case 0x1a6:
      return "Hitachi SH4";
    case 0x1a8:
      return "Hitachi SH5";
    case 0x1c2:
      return "Thumb";
    case 0x169:
      return "MIPS little-endian WCE v2";
  }
  return "?";
}

function unixtime2Human(unixtime) {
  var newDate = new Date();
  newDate.setTime(unixtime * 1000);
  var dateString = newDate.toUTCString();
  return dateString;
}

function fileCharacteristics2String(characteristics) {
  var fileHeaderCharacteristics = {
    0x0001: "Relocation info stripped",
    0x0002: "File is executable",
    0x0004: "Line numbers stripped",
    0x0008: "Local symbols stripped",
    0x0010: "Aggressively trim working set",
    0x0020: "Large Address Aware",
    0x0080: "Bytes of machine word are reversed",
    0x0100: "32 bit word machine",
    0x0200: "Debugging info stripped",
    0x0400: "Removable Run From Swap",
    0x0800: "Net Run From Swap",
    0x1000: "System File",
    0x2000: "File is a DLL",
    0x4000: "Up System Only",
    0x8000: "Bytes of machine word are reversed",
  };
  var values = [];
  for (var k in fileHeaderCharacteristics) {
    if (k & characteristics) {
      values.push(fileHeaderCharacteristics[k]);
    }
  }
  return values;
}

function dllCharacteristics2String(characteristics) {
  var dllCharacteristics = {
    0x0001: "Reserved 0x1",
    0x0002: "Reserved 0x2",
    0x0004: "Reserved 0x4",
    0x0008: "Reserved 0x8",
    0x0020: "High Entropy VA",
    0x0040: "Dynamic Base",
    0x0080: "Code Integrity checks are enforced",
    0x0100: "NX compatible",
    0x0200: "No isolation",
    0x0400: "No SEH",
    0x0800: "Do not bind the image",
    0x1000: "Must execute in an AppContainer",
    0x2000: "WDM driveR",
    0x4000: "Supports Control Flow Guard",
    0x8000: "Terminal Server aware",
  };

  var values = [];

  for (var k in dllCharacteristics) {
    if (k & characteristics) {
      values.push(dllCharacteristics[k]);
    }
  }

  return values;
}

function subsystem2String(subsystem) {
  var subsystemMap = {
    0: "Unknown",
    1: "Native",
    2: "Windows GUI",
    3: "Windows CUI",
    5: "OS/2 character",
    7: "POSIX character",
    8: "Native Win9x driver",
    9: "Windows CE GUI",
    10: "EFI Application",
    11: "EFI Boot Service Driver",
    12: "EFI ROM image",
    13: "EFI ROM image",
    14: "XBOX",
    15: "Windows boot application",
  };

  if ( subsystem in subsystemMap) {
    return subsystemMap[subsystem];
  } 
  
  return "?"
}

export {
  hex2a,
  reverse,
  dec2Hex,
  dec2HexString,
  prodIdToStr,
  prodIdToVsVersion,
  machine2String,
  unixtime2Human,
  fileCharacteristics2String,
  dllCharacteristics2String,
  subsystem2String
}
