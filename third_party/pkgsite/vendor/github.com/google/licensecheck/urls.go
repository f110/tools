// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package licensecheck

// This file contains the list of known URLs for valid licenses, only.
// Do not add other code.

// TODO: Find a canonical source of this information, or at least a
// disciplined way to develop it. This is cadged from gopkg.in/src-d/go-license-detector.v2.

// The URL text does not include the leading http:// or https://
// or trailing /.
// All entries are lower case.
// Keep this list sorted for easy checking.
var builtinURLs = []License{
	{URL: "creativecommons.org/licenses/by-nc-nd/2.0", ID: "CC-BY-NC-ND-2.0"},
	{URL: "creativecommons.org/licenses/by-nc-nd/2.5", ID: "CC-BY-NC-ND-2.5"},
	{URL: "creativecommons.org/licenses/by-nc-nd/3.0", ID: "CC-BY-NC-ND-3.0"},
	{URL: "creativecommons.org/licenses/by-nc-nd/4.0", ID: "CC-BY-NC-ND-4.0"},
	{URL: "creativecommons.org/licenses/by-nc-sa/1.0", ID: "CC-BY-NC-SA-1.0"},
	{URL: "creativecommons.org/licenses/by-nc-sa/2.0", ID: "CC-BY-NC-SA-2.0"},
	{URL: "creativecommons.org/licenses/by-nc-sa/2.5", ID: "CC-BY-NC-SA-2.5"},
	{URL: "creativecommons.org/licenses/by-nc-sa/3.0", ID: "CC-BY-NC-SA-3.0"},
	{URL: "creativecommons.org/licenses/by-nc-sa/4.0", ID: "CC-BY-NC-SA-4.0"},
	{URL: "creativecommons.org/licenses/by-nc/1.0", ID: "CC-BY-NC-1.0"},
	{URL: "creativecommons.org/licenses/by-nc/2.0", ID: "CC-BY-NC-2.0"},
	{URL: "creativecommons.org/licenses/by-nc/2.5", ID: "CC-BY-NC-2.5"},
	{URL: "creativecommons.org/licenses/by-nc/3.0", ID: "CC-BY-NC-3.0"},
	{URL: "creativecommons.org/licenses/by-nc/4.0", ID: "CC-BY-NC-4.0"},
	{URL: "creativecommons.org/licenses/by-nd-nc/1.0", ID: "CC-BY-NC-ND-1.0"},
	{URL: "creativecommons.org/licenses/by-nd/1.0", ID: "CC-BY-ND-1.0"},
	{URL: "creativecommons.org/licenses/by-nd/2.0", ID: "CC-BY-ND-2.0"},
	{URL: "creativecommons.org/licenses/by-nd/2.5", ID: "CC-BY-ND-2.5"},
	{URL: "creativecommons.org/licenses/by-nd/3.0", ID: "CC-BY-ND-3.0"},
	{URL: "creativecommons.org/licenses/by-nd/4.0", ID: "CC-BY-ND-4.0"},
	{URL: "creativecommons.org/licenses/by-sa/1.0", ID: "CC-BY-SA-1.0"},
	{URL: "creativecommons.org/licenses/by-sa/2.0", ID: "CC-BY-SA-2.0"},
	{URL: "creativecommons.org/licenses/by-sa/2.5", ID: "CC-BY-SA-2.5"},
	{URL: "creativecommons.org/licenses/by-sa/3.0", ID: "CC-BY-SA-3.0"},
	{URL: "creativecommons.org/licenses/by-sa/4.0", ID: "CC-BY-SA-4.0"},
	{URL: "creativecommons.org/licenses/by/1.0", ID: "CC-BY-1.0"},
	{URL: "creativecommons.org/licenses/by/2.0", ID: "CC-BY-2.0"},
	{URL: "creativecommons.org/licenses/by/2.5", ID: "CC-BY-2.5"},
	{URL: "creativecommons.org/licenses/by/3.0", ID: "CC-BY-3.0"},
	{URL: "creativecommons.org/licenses/by/4.0", ID: "CC-BY-4.0"},
	{URL: "creativecommons.org/publicdomain/zero/1.0", ID: "CC0-1.0"},
	{URL: "opensource.org/licenses/apache-1.1", ID: "Apache-1.1"},
	{URL: "opensource.org/licenses/artistic-1.0", ID: "Artistic-1.0"},
	{URL: "opensource.org/licenses/bsdpluspatent", ID: "BSD-2-Clause-Patent"},
	{URL: "opensource.org/licenses/catosl-1.1", ID: "CATOSL-1.1"},
	{URL: "opensource.org/licenses/cpl-1.0", ID: "CPL-1.0"},
	{URL: "opensource.org/licenses/cua-opl-1.0", ID: "CUA-OPL-1.0"},
	{URL: "opensource.org/licenses/ecl-1.0", ID: "ECL-1.0"},
	{URL: "opensource.org/licenses/ecl-2.0", ID: "ECL-2.0"},
	{URL: "opensource.org/licenses/efl-1.0", ID: "EFL-1.0"},
	{URL: "opensource.org/licenses/efl-2.0", ID: "EFL-2.0"},
	{URL: "opensource.org/licenses/entessa", ID: "Entessa"},
	{URL: "opensource.org/licenses/intel", ID: "Intel"},
	{URL: "opensource.org/licenses/lpl-1.0", ID: "LPL-1.0"},
	{URL: "opensource.org/licenses/liliq-p-1.1", ID: "LiLiQ-P-1.1"},
	{URL: "opensource.org/licenses/liliq-r-1.1", ID: "LiLiQ-R-1.1"},
	{URL: "opensource.org/licenses/liliq-rplus-1.1", ID: "LiLiQ-Rplus-1.1"},
	{URL: "opensource.org/licenses/mpl-1.0", ID: "MPL-1.0"},
	{URL: "opensource.org/licenses/mpl-2.0", ID: "MPL-2.0"},
	{URL: "opensource.org/licenses/opl-2.1", ID: "OSET-PL-2.1"},
	{URL: "opensource.org/licenses/osl-1.0", ID: "OSL-1.0"},
	{URL: "opensource.org/licenses/osl-2.1", ID: "OSL-2.1"},
	{URL: "opensource.org/licenses/rpl-1.1", ID: "RPL-1.1"},
	{URL: "opensource.org/licenses/sissl", ID: "SISSL"},
	{URL: "opensource.org/licenses/upl", ID: "UPL-1.0"},
	{URL: "opensource.org/licenses/xnet", ID: "Xnet"},
	{URL: "opensource.org/licenses/zpl-2.0", ID: "ZPL-2.0"},
	{URL: "www.apache.org/licenses/license-2.0", ID: "Apache-2.0"},
	{URL: "www.gnu.org/licenses/agpl.txt", ID: "AGPL-3.0"},
	// {URL: "www.gnu.org/licenses/autoconf-exception-3.0.html", ID: "GPL-3.0-with-autoconf-exception"},
	// {URL: "www.gnu.org/licenses/ecos-license.html", ID: "eCos-2.0"},
	{URL: "www.gnu.org/licenses/fdl-1.3.txt", ID: "GFDL-1.3"},
	// {URL: "www.gnu.org/licenses/gcc-exception-3.1.html", ID: "GPL-3.0-with-GCC-exception"},
	{URL: "www.gnu.org/licenses/gpl-3.0-standalone.html", ID: "GPL-3.0"},
	// {URL: "www.gnu.org/licenses/gpl-faq.html#fontexception", ID: "GPL-2.0-with-font-exception"},
	{URL: "www.gnu.org/licenses/lgpl-3.0-standalone.html", ID: "LGPL-3.0"},
	{URL: "www.gnu.org/licenses/old-licenses/fdl-1.1.txt", ID: "GFDL-1.1"},
	{URL: "www.gnu.org/licenses/old-licenses/gpl-1.0-standalone.html", ID: "GPL-1.0"},
	{URL: "www.gnu.org/licenses/old-licenses/gpl-2.0-standalone.html", ID: "GPL-2.0"},
	{URL: "www.gnu.org/licenses/old-licenses/lgpl-2.0-standalone.html", ID: "LGPL-2.0"},
	{URL: "www.gnu.org/licenses/old-licenses/lgpl-2.1-standalone.html", ID: "LGPL-2.1"},
	{URL: "www.gnu.org/prep/maintain/html_node/license-notices-for-other-files.html", ID: "FSFAP"},
	// {URL: "www.gnu.org/software/classpath/license.html", ID: "GPL-2.0-with-classpath-exception"},
	{URL: "www.opensource.org/licenses/agpl-3.0", ID: "AGPL-3.0"},
	{URL: "www.opensource.org/licenses/apl-1.0", ID: "APL-1.0"},
	{URL: "www.opensource.org/licenses/apache-2.0", ID: "Apache-2.0"},
	{URL: "www.opensource.org/licenses/bsd-2-clause", ID: "BSD-2-Clause"},
	{URL: "www.opensource.org/licenses/bsd-3-clause", ID: "BSD-3-Clause"},
	{URL: "www.opensource.org/licenses/bsl-1.0", ID: "BSL-1.0"},
	{URL: "www.opensource.org/licenses/cnri-python", ID: "CNRI-Python"},
	{URL: "www.opensource.org/licenses/cpal-1.0", ID: "CPAL-1.0"},
	{URL: "www.opensource.org/licenses/epl-1.0", ID: "EPL-1.0"},
	{URL: "www.opensource.org/licenses/epl-2.0", ID: "EPL-2.0"},
	{URL: "www.opensource.org/licenses/eudatagrid", ID: "EUDatagrid"},
	{URL: "www.opensource.org/licenses/eupl-1.1", ID: "EUPL-1.1"},
	{URL: "www.opensource.org/licenses/fair", ID: "Fair"},
	{URL: "www.opensource.org/licenses/frameworx-1.0", ID: "Frameworx-1.0"},
	{URL: "www.opensource.org/licenses/gpl-2.0", ID: "GPL-2.0"},
	{URL: "www.opensource.org/licenses/gpl-3.0", ID: "GPL-3.0"},
	{URL: "www.opensource.org/licenses/hpnd", ID: "HPND"},
	{URL: "www.opensource.org/licenses/ipa", ID: "IPA"},
	{URL: "www.opensource.org/licenses/ipl-1.0", ID: "IPL-1.0"},
	{URL: "www.opensource.org/licenses/isc", ID: "ISC"},
	{URL: "www.opensource.org/licenses/lgpl-2.1", ID: "LGPL-2.1"},
	{URL: "www.opensource.org/licenses/lgpl-3.0", ID: "LGPL-3.0"},
	{URL: "www.opensource.org/licenses/lpl-1.02", ID: "LPL-1.02"},
	{URL: "www.opensource.org/licenses/lppl-1.3c", ID: "LPPL-1.3c"},
	{URL: "www.opensource.org/licenses/mit", ID: "MIT"},
	{URL: "www.opensource.org/licenses/mpl-1.1", ID: "MPL-1.1"},
	{URL: "www.opensource.org/licenses/ms-pl", ID: "MS-PL"},
	{URL: "www.opensource.org/licenses/ms-rl", ID: "MS-RL"},
	{URL: "www.opensource.org/licenses/miros", ID: "MirOS"},
	{URL: "www.opensource.org/licenses/motosoto", ID: "Motosoto"},
	{URL: "www.opensource.org/licenses/multics", ID: "Multics"},
	{URL: "www.opensource.org/licenses/nasa-1.3", ID: "NASA-1.3"},
	{URL: "www.opensource.org/licenses/ncsa", ID: "NCSA"},
	{URL: "www.opensource.org/licenses/ngpl", ID: "NGPL"},
	{URL: "www.opensource.org/licenses/nosl3.0", ID: "NPOSL-3.0"},
	{URL: "www.opensource.org/licenses/ntp", ID: "NTP"},
	{URL: "www.opensource.org/licenses/naumen", ID: "Naumen"},
	{URL: "www.opensource.org/licenses/oclc-2.0", ID: "OCLC-2.0"},
	{URL: "www.opensource.org/licenses/ofl-1.1", ID: "OFL-1.1"},
	{URL: "www.opensource.org/licenses/ogtsl", ID: "OGTSL"},
	{URL: "www.opensource.org/licenses/osl-3.0", ID: "OSL-3.0"},
	{URL: "www.opensource.org/licenses/php-3.0", ID: "PHP-3.0"},
	{URL: "www.opensource.org/licenses/postgresql", ID: "PostgreSQL"},
	{URL: "www.opensource.org/licenses/python-2.0", ID: "Python-2.0"},
	{URL: "www.opensource.org/licenses/qpl-1.0", ID: "QPL-1.0"},
	{URL: "www.opensource.org/licenses/rpl-1.5", ID: "RPL-1.5"},
	{URL: "www.opensource.org/licenses/rpsl-1.0", ID: "RPSL-1.0"},
	{URL: "www.opensource.org/licenses/rscpl", ID: "RSCPL"},
	{URL: "www.opensource.org/licenses/spl-1.0", ID: "SPL-1.0"},
	{URL: "www.opensource.org/licenses/simpl-2.0", ID: "SimPL-2.0"},
	{URL: "www.opensource.org/licenses/sleepycat", ID: "Sleepycat"},
	{URL: "www.opensource.org/licenses/vsl-1.0", ID: "VSL-1.0"},
	{URL: "www.opensource.org/licenses/w3c", ID: "W3C"},
	// {URL: "www.opensource.org/licenses/wxwindows", ID: "wxWindows"},
	{URL: "www.opensource.org/licenses/watcom-1.0", ID: "Watcom-1.0"},
	{URL: "www.opensource.org/licenses/zlib", ID: "Zlib"},
	{URL: "www.opensource.org/licenses/afl-3.0", ID: "AFL-3.0"},
	{URL: "www.opensource.org/licenses/artistic-license-2.0", ID: "Artistic-2.0"},
	{URL: "www.opensource.org/licenses/attribution", ID: "AAL"},
	{URL: "www.opensource.org/licenses/cddl1", ID: "CDDL-1.0"},
	{URL: "www.opensource.org/licenses/nokia", ID: "Nokia"},
}