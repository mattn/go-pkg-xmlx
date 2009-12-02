package xmlx

import "fmt"
import "utf8"
import "regexp"
import "strconv"

var reg_entity = regexp.MustCompile("^&#[0-9]+;$");

// Converts a single numerical html entity to a regular Go utf-token.
//    ex: "&#9827;" -> "♣"
func HtmlToUTF8(entity string) string {
	// Make sure we have a valid entity: &#123;
	ok := reg_entity.MatchString(entity);
	if !ok { return "" }

	// Convert entity to number
	num, err := strconv.Atoi(entity[2:len(entity)-1]);
	if err != nil { return "" }

	var arr [3]byte;
	size := utf8.EncodeRune(num, &arr);
	if size == 0 { return "" }

	return string(&arr);
}

// Converts a single Go utf-token to it's an Html entity.
//   ex: "♣" -> "&#9827;"
func UTF8ToHtml(token string) string {
	rune, size := utf8.DecodeRuneInString(token);
	if size == 0 { return "" }
	return fmt.Sprintf("&#%d;", rune);
}

/*
	http://www.w3.org/TR/html4/sgml/entities.html

	Portions © International Organization for Standardization 1986
	Permission to copy in any form is granted for use with
	conforming SGML systems and applications as defined in
	ISO 8879, provided this notice is included in all copies.

	Fills the supplied map with html entities mapped to their Go utf8
	equivalents. This map can be assigned to xml.Parser.Entity
	It will be used to map non-standard xml entities to a proper value.
	If the parser encounters any unknown entities, it will throw a syntax
	error and abort the parsing. Hence the ability to supply this map.
 */
func loadNonStandardEntities(em *map[string]string) {
	(*em)["pi"] = "\u03c0";
	(*em)["nabla"] = "\u2207";
	(*em)["isin"] = "\u2208";
	(*em)["loz"] = "\u25ca";
	(*em)["prop"] = "\u221d";
	(*em)["para"] = "\u00b6";
	(*em)["Aring"] = "\u00c5";
	(*em)["euro"] = "\u20ac";
	(*em)["sup3"] = "\u00b3";
	(*em)["sup2"] = "\u00b2";
	(*em)["sup1"] = "\u00b9";
	(*em)["prod"] = "\u220f";
	(*em)["gamma"] = "\u03b3";
	(*em)["perp"] = "\u22a5";
	(*em)["lfloor"] = "\u230a";
	(*em)["fnof"] = "\u0192";
	(*em)["frasl"] = "\u2044";
	(*em)["rlm"] = "\u200f";
	(*em)["omega"] = "\u03c9";
	(*em)["part"] = "\u2202";
	(*em)["euml"] = "\u00eb";
	(*em)["Kappa"] = "\u039a";
	(*em)["nbsp"] = "\u00a0";
	(*em)["Eacute"] = "\u00c9";
	(*em)["brvbar"] = "\u00a6";
	(*em)["otimes"] = "\u2297";
	(*em)["ndash"] = "\u2013";
	(*em)["thinsp"] = "\u2009";
	(*em)["nu"] = "\u03bd";
	(*em)["Upsilon"] = "\u03a5";
	(*em)["upsih"] = "\u03d2";
	(*em)["raquo"] = "\u00bb";
	(*em)["yacute"] = "\u00fd";
	(*em)["delta"] = "\u03b4";
	(*em)["eth"] = "\u00f0";
	(*em)["supe"] = "\u2287";
	(*em)["ne"] = "\u2260";
	(*em)["ni"] = "\u220b";
	(*em)["eta"] = "\u03b7";
	(*em)["uArr"] = "\u21d1";
	(*em)["image"] = "\u2111";
	(*em)["asymp"] = "\u2248";
	(*em)["oacute"] = "\u00f3";
	(*em)["rarr"] = "\u2192";
	(*em)["emsp"] = "\u2003";
	(*em)["acirc"] = "\u00e2";
	(*em)["shy"] = "\u00ad";
	(*em)["yuml"] = "\u00ff";
	(*em)["acute"] = "\u00b4";
	(*em)["int"] = "\u222b";
	(*em)["ccedil"] = "\u00e7";
	(*em)["Acirc"] = "\u00c2";
	(*em)["Ograve"] = "\u00d2";
	(*em)["times"] = "\u00d7";
	(*em)["weierp"] = "\u2118";
	(*em)["Tau"] = "\u03a4";
	(*em)["omicron"] = "\u03bf";
	(*em)["lt"] = "\u003c";
	(*em)["Mu"] = "\u039c";
	(*em)["Ucirc"] = "\u00db";
	(*em)["sub"] = "\u2282";
	(*em)["le"] = "\u2264";
	(*em)["sum"] = "\u2211";
	(*em)["sup"] = "\u2283";
	(*em)["lrm"] = "\u200e";
	(*em)["frac34"] = "\u00be";
	(*em)["Iota"] = "\u0399";
	(*em)["Ugrave"] = "\u00d9";
	(*em)["THORN"] = "\u00de";
	(*em)["rsaquo"] = "\u203a";
	(*em)["not"] = "\u00ac";
	(*em)["sigma"] = "\u03c3";
	(*em)["iuml"] = "\u00ef";
	(*em)["epsilon"] = "\u03b5";
	(*em)["spades"] = "\u2660";
	(*em)["theta"] = "\u03b8";
	(*em)["divide"] = "\u00f7";
	(*em)["Atilde"] = "\u00c3";
	(*em)["uacute"] = "\u00fa";
	(*em)["Rho"] = "\u03a1";
	(*em)["trade"] = "\u2122";
	(*em)["chi"] = "\u03c7";
	(*em)["agrave"] = "\u00e0";
	(*em)["or"] = "\u2228";
	(*em)["circ"] = "\u02c6";
	(*em)["middot"] = "\u00b7";
	(*em)["plusmn"] = "\u00b1";
	(*em)["aring"] = "\u00e5";
	(*em)["lsquo"] = "\u2018";
	(*em)["Yacute"] = "\u00dd";
	(*em)["oline"] = "\u203e";
	(*em)["copy"] = "\u00a9";
	(*em)["icirc"] = "\u00ee";
	(*em)["lowast"] = "\u2217";
	(*em)["Oacute"] = "\u00d3";
	(*em)["aacute"] = "\u00e1";
	(*em)["oplus"] = "\u2295";
	(*em)["crarr"] = "\u21b5";
	(*em)["thetasym"] = "\u03d1";
	(*em)["Beta"] = "\u0392";
	(*em)["laquo"] = "\u00ab";
	(*em)["rang"] = "\u232a";
	(*em)["tilde"] = "\u02dc";
	(*em)["Uuml"] = "\u00dc";
	(*em)["zwj"] = "\u200d";
	(*em)["mu"] = "\u03bc";
	(*em)["Ccedil"] = "\u00c7";
	(*em)["infin"] = "\u221e";
	(*em)["ouml"] = "\u00f6";
	(*em)["rfloor"] = "\u230b";
	(*em)["pound"] = "\u00a3";
	(*em)["szlig"] = "\u00df";
	(*em)["thorn"] = "\u00fe";
	(*em)["forall"] = "\u2200";
	(*em)["piv"] = "\u03d6";
	(*em)["rdquo"] = "\u201d";
	(*em)["frac12"] = "\u00bd";
	(*em)["frac14"] = "\u00bc";
	(*em)["Ocirc"] = "\u00d4";
	(*em)["Ecirc"] = "\u00ca";
	(*em)["kappa"] = "\u03ba";
	(*em)["Euml"] = "\u00cb";
	(*em)["minus"] = "\u2212";
	(*em)["cong"] = "\u2245";
	(*em)["hellip"] = "\u2026";
	(*em)["equiv"] = "\u2261";
	(*em)["cent"] = "\u00a2";
	(*em)["Uacute"] = "\u00da";
	(*em)["darr"] = "\u2193";
	(*em)["Eta"] = "\u0397";
	(*em)["sbquo"] = "\u201a";
	(*em)["rArr"] = "\u21d2";
	(*em)["igrave"] = "\u00ec";
	(*em)["uml"] = "\u00a8";
	(*em)["lambda"] = "\u03bb";
	(*em)["oelig"] = "\u0153";
	(*em)["harr"] = "\u2194";
	(*em)["ang"] = "\u2220";
	(*em)["clubs"] = "\u2663";
	(*em)["and"] = "\u2227";
	(*em)["permil"] = "\u2030";
	(*em)["larr"] = "\u2190";
	(*em)["Yuml"] = "\u0178";
	(*em)["cup"] = "\u222a";
	(*em)["Xi"] = "\u039e";
	(*em)["Alpha"] = "\u0391";
	(*em)["phi"] = "\u03c6";
	(*em)["ucirc"] = "\u00fb";
	(*em)["oslash"] = "\u00f8";
	(*em)["rsquo"] = "\u2019";
	(*em)["AElig"] = "\u00c6";
	(*em)["mdash"] = "\u2014";
	(*em)["psi"] = "\u03c8";
	(*em)["eacute"] = "\u00e9";
	(*em)["otilde"] = "\u00f5";
	(*em)["yen"] = "\u00a5";
	(*em)["gt"] = "\u003e";
	(*em)["Iuml"] = "\u00cf";
	(*em)["Prime"] = "\u2033";
	(*em)["Chi"] = "\u03a7";
	(*em)["ge"] = "\u2265";
	(*em)["reg"] = "\u00ae";
	(*em)["hearts"] = "\u2665";
	(*em)["auml"] = "\u00e4";
	(*em)["Agrave"] = "\u00c0";
	(*em)["sect"] = "\u00a7";
	(*em)["sube"] = "\u2286";
	(*em)["sigmaf"] = "\u03c2";
	(*em)["Gamma"] = "\u0393";
	(*em)["amp"] = "\u0026";
	(*em)["ensp"] = "\u2002";
	(*em)["ETH"] = "\u00d0";
	(*em)["Igrave"] = "\u00cc";
	(*em)["Omega"] = "\u03a9";
	(*em)["Lambda"] = "\u039b";
	(*em)["Omicron"] = "\u039f";
	(*em)["there4"] = "\u2234";
	(*em)["ntilde"] = "\u00f1";
	(*em)["xi"] = "\u03be";
	(*em)["dagger"] = "\u2020";
	(*em)["egrave"] = "\u00e8";
	(*em)["Delta"] = "\u0394";
	(*em)["OElig"] = "\u0152";
	(*em)["diams"] = "\u2666";
	(*em)["ldquo"] = "\u201c";
	(*em)["radic"] = "\u221a";
	(*em)["Oslash"] = "\u00d8";
	(*em)["Ouml"] = "\u00d6";
	(*em)["lceil"] = "\u2308";
	(*em)["uarr"] = "\u2191";
	(*em)["atilde"] = "\u00e3";
	(*em)["iquest"] = "\u00bf";
	(*em)["lsaquo"] = "\u2039";
	(*em)["Epsilon"] = "\u0395";
	(*em)["iacute"] = "\u00ed";
	(*em)["cap"] = "\u2229";
	(*em)["deg"] = "\u00b0";
	(*em)["Otilde"] = "\u00d5";
	(*em)["zeta"] = "\u03b6";
	(*em)["ocirc"] = "\u00f4";
	(*em)["scaron"] = "\u0161";
	(*em)["ecirc"] = "\u00ea";
	(*em)["ordm"] = "\u00ba";
	(*em)["tau"] = "\u03c4";
	(*em)["Auml"] = "\u00c4";
	(*em)["dArr"] = "\u21d3";
	(*em)["ordf"] = "\u00aa";
	(*em)["alefsym"] = "\u2135";
	(*em)["notin"] = "\u2209";
	(*em)["Pi"] = "\u03a0";
	(*em)["sdot"] = "\u22c5";
	(*em)["upsilon"] = "\u03c5";
	(*em)["iota"] = "\u03b9";
	(*em)["hArr"] = "\u21d4";
	(*em)["Sigma"] = "\u03a3";
	(*em)["lang"] = "\u2329";
	(*em)["curren"] = "\u00a4";
	(*em)["Theta"] = "\u0398";
	(*em)["lArr"] = "\u21d0";
	(*em)["Phi"] = "\u03a6";
	(*em)["Nu"] = "\u039d";
	(*em)["rho"] = "\u03c1";
	(*em)["alpha"] = "\u03b1";
	(*em)["iexcl"] = "\u00a1";
	(*em)["micro"] = "\u00b5";
	(*em)["cedil"] = "\u00b8";
	(*em)["Ntilde"] = "\u00d1";
	(*em)["Psi"] = "\u03a8";
	(*em)["Dagger"] = "\u2021";
	(*em)["Egrave"] = "\u00c8";
	(*em)["Icirc"] = "\u00ce";
	(*em)["nsub"] = "\u2284";
	(*em)["bdquo"] = "\u201e";
	(*em)["empty"] = "\u2205";
	(*em)["aelig"] = "\u00e6";
	(*em)["ograve"] = "\u00f2";
	(*em)["macr"] = "\u00af";
	(*em)["Zeta"] = "\u0396";
	(*em)["beta"] = "\u03b2";
	(*em)["sim"] = "\u223c";
	(*em)["uuml"] = "\u00fc";
	(*em)["Aacute"] = "\u00c1";
	(*em)["Iacute"] = "\u00cd";
	(*em)["exist"] = "\u2203";
	(*em)["prime"] = "\u2032";
	(*em)["rceil"] = "\u2309";
	(*em)["real"] = "\u211c";
	(*em)["zwnj"] = "\u200c";
	(*em)["bull"] = "\u2022";
	(*em)["quot"] = "\u0022";
	(*em)["Scaron"] = "\u0160";
	(*em)["ugrave"] = "\u00f9";
}



