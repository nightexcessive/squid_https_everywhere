package https_everywhere

import (
	"bytes"
	"testing"
)

const testFile = `
<!--
	Other Google rulesets:

		- 2mdn.net.xml
		- Admeld.xml
		- ChannelIntelligence.com.xml
		- Doubleclick.net.xml
		- FeedBurner.xml
		- Google.org.xml
		- GoogleAPIs.xml
		- Google_App_Engine.xml
		- GoogleImages.xml
		- GoogleShopping.xml
		- Ingress.xml
		- Meebo.xml
		- Orkut.xml
		- Postini.xml


	Nonfunctional domains:

		- feedproxy.google.com			(404, valid cert)
		- partnerpage.google.com *
		- safebrowsing.clients.google.com	(404, mismatched)
		- (www.)googlesyndicatedsearch.com	(404; mismatched, CN: google.com)
		- buttons.googlesyndication.com *

	* 404, valid cert


	Nonfunctional google.com paths:

		- analytics	(redirects to http)
		- imgres
		- gadgets *
		- hangouts	(404)
		- u/		(404)

	* Redirects to http


	Problematic domains:

		- www.goo.gl		(404; mismatched, CN: *.google.com)

		- google.com subdomains:

			- books		(googlebooks/, images/, & intl/ 404, but works when rewritten to www)
			- cbks0 ****
			- gg		($ 404s)
			- knoll *
			- scholar **
			- trends *

		- news.google.cctld **
		- scholar.google.cctld **
		- *-opensocial.googleusercontent.com ***

	**** $ 404s
	* 404, valid cert
	** Redirects to http, valid cert
	*** Breaks followers widget - https://trac.torproject.org/projects/tor/ticket/7294


	Partially covered domains:

		- google.cctld subdomains:

			- scholar	(→ www)

		- google.com subdomains:

			- (www.)
			- cbks0		($ 404s)
			- gg		($ 404s)
			- news		(→ www)
			- scholar	(→ www)

		- *.googleusercontent.com	(*-opensocial excluded)


	Fully covered domains:

		- lh[3-6].ggpht.com
		- (www.)goo.gl		(www → ^)

		- google.com subdomains:

			- accounts
			- adwords
			- apis
			- appengine
			- books		(→ encrypted)
			- calendar
			- checkout
			- chrome
			- clients[12]
			- code
			- *.corp
			- developers
			- dl
			- docs
			- docs\d
			- drive
			- encrypted
			- encrypted-tbn[123]
			- feedburner
			- fiber
			- finance
			- glass
			- groups
			- health
			- helpouts
			- history
			- hostedtalkgadget
			- id
			- investor
			- knol
			- knoll		(→ knol)
			- lh\d
			- mail
			- chatenabled.mail
			- pack
			- picasaweb
			- pki
			- play
			- plus
			- plusone
			- productforums
			- profiles
			- safebrowsing-cache
			- cert-test.sandbox
			- plus.sandbox
			- sb-ssl
			- script
			- security
			- servicessites
			- sites
			- spreadsheets
			- spreadsheets\d
			- support
			- talk
			- talkgadget
			- tbn2			(→ encrypted-tbn2)
			- tools
			- translate
			- trends		(→ www)

		- partner.googleadservices.com
		- (www.)googlecode.com
		- *.googlecode.com	(per-project subdomains)
		- googlesource.com
		- *.googlesource.com
		- pagead2.googlesyndication.com
		- tpc.googlesyndication.com
		- mail-attachment.googleusercontent.com
		- webcache.googleusercontent.com


	XXX: Needs more testing

-->
<ruleset name="Google Services">

	<target host="*.ggpht.com" />
	<target host="gmail.com" />
	<target host="www.gmail.com" />
	<target host="goo.gl" />
	<target host="www.goo.gl" />
	<target host="google.*" />
	<target host="accounts.google.*" />
	<target host="adwords.google.*" />
	<target host="finance.google.*" />
	<target host="groups.google.*" />
	<target host="it.google.*" />
	<target host="news.google.*" />
		<exclusion pattern="^http://(?:news\.)?google\.com/(?:archivesearch|newspapers)" />
	<target host="picasaweb.google.*" />
	<target host="scholar.google.*" />
	<target host="translate.google.*" />
	<target host="www.google.*" />
	<target host="*.google.ca" />
	<target host="google.co.*" />
	<target host="accounts.google.co.*" />
	<target host="adwords.google.co.*" />
	<target host="finance.google.co.*" />
	<target host="groups.google.co.*" />
	<target host="id.google.co.*" />
	<target host="news.google.co.*" />
	<target host="picasaweb.google.co.*" />
	<target host="scholar.google.co.*" />
	<target host="translate.google.co.*" />
	<target host="www.google.co.*" />
	<target host="google.com" />
	<target host="*.google.com" />
		<exclusion pattern="^http://(?:www\.)?google\.com/analytics/*(?:/[^/]+)?(?:\?.*)?$" />
		<!--exclusion pattern="^http://books\.google\.com/(?!books/(\w+\.js|css/|javascript/)|favicon\.ico|googlebooks/|images/|intl/)" /-->
		<exclusion pattern="^http://cbks0\.google\.com/(?:$|\?)" />
		<exclusion pattern="^http://gg\.google\.com/(?!csi(?:$|\?))" />
	<target host="google.com.*" />
	<target host="accounts.google.com.*" />
	<target host="adwords.google.com.*" />
	<target host="groups.google.com.*" />
	<target host="id.google.com.*" />
	<target host="news.google.com.*" />
	<target host="picasaweb.google.com.*" />
	<target host="scholar.google.com.*" />
	<target host="translate.google.com.*" />
	<target host="www.google.com.*" />
	<target host="partner.googleadservices.com" />
	<target host="googlecode.com" />
	<target host="*.googlecode.com" />
	<target host="googlemail.com" />
	<target host="www.googlemail.com" />
	<target host="googlesource.com" />
	<target host="*.googlesource.com" />
	<target host="*.googlesyndication.com" />
	<target host="www.googletagservices.com" />
	<target host="googleusercontent.com" />
	<target host="*.googleusercontent.com" />
		<!--
			Necessary for the Followers widget:

				 https://trac.torproject.org/projects/tor/ticket/7294
											-->
		<exclusion pattern="http://[^@:\./]+-opensocial\.googleusercontent\.com" />


	<!--	Can we secure any of these wildcard cookies safely?
									-->
	<!--securecookie host="^\.google\.com$" name="^(hl|I4SUserLocale|NID|PREF|S)$" /-->
	<!--securecookie host="^\.google\.[\w.]{2,6}$" name="^(hl|I4SUserLocale|NID|PREF|S|S_awfe)$" /-->
	<securecookie host="^(?:accounts|adwords|\.code|login\.corp|developers|docs|\d\.docs|fiber|mail|picasaweb|plus|\.?productforums|support)\.google\.[\w.]{2,6}$" name=".+" />
	<securecookie host="^www\.google\.com$" name="^GoogleAccountsLocale_session$" />
	<securecookie host="^mail-attachment\.googleusercontent\.com$" name=".+" />
	<securecookie host="^gmail\.com$" name=".+" />
	<securecookie host="^www\.gmail\.com$" name=".+" />
	<securecookie host="^googlemail\.com$" name=".+" />
	<securecookie host="^www\.googlemail\.com$" name=".+" />


	<!--    - lh 3-6 exist
		- All appear identical
		- Identical to lh\d.googleusercontent.com
					-->
	<rule from="^http://lh(\d)\.ggpht\.com/"
		to="https://lh$1.ggpht.com/" />

	<rule from="^http://lh(\d)\.google\.ca/"
		to="https://lh$1.google.ca/" />


	<rule from="^http://(www\.)?g(oogle)?mail\.com/"
		to="https://$1g$2mail.com/" />

	<rule from="^http://(?:www\.)?goo\.gl/"
		to="https://goo.gl/" />


	<!--	Redirects to http when rewritten to www:
							-->
	<rule from="^http://books\.google\.com/"
		to="https://encrypted.google.com/" />

	<!--	Paths that work on all in google.*
							-->
	<rule from="^http://(?:www\.)?google\.((?:com?\.)?\w{2,3})/(accounts|adplanner|ads|adsense|adwords|analytics|bookmarks|chrome|contacts|coop|cse|css|culturalinstitute|doodles|favicon\.ico|finance|goodtoknow|googleblogs|green|hostednews|images|intl|js|landing|logos|mapmaker|newproducts|news|nexus|patents|policies|prdhp|profiles|products|reader|s2|settings|shopping|support|tools|transparencyreport|trends|urchin|webmasters)(?=$|[?/])"
		 to="https://www.google.$1/$2" />

	<!--	Paths that 404 on .ccltd, but work on .com:
								-->
	<rule from="^http://(?:www\.)?google\.(?:com?\.)?\w{2,3}/(calendar|dictionary|doubleclick|help|ideas|postini|powermeter|url)"
		 to="https://www.google.com/$1" />

	<rule from="^http://(?:www\.)?google\.(?:com?\.)?\w{2,3}/custom"
		 to="https://www.google.com/cse" />

	<!--	Paths that only exist/work on .com
							-->
	<rule from="^http://(?:www\.)?google\.com/(\+|appsstatus|books|buzz|extern_js|glass|googlebooks|ig|insights|moderator|phone|safebrowsing|videotargetting|webfonts)($|[?/])"
		to="https://www.google.com/$1$2" />

	<!--	Subdomains that work on all in google.*
							-->
	<rule from="^http://(accounts|adwords|finance|groups|id|picasaweb|translate)\.google\.((?:com?\.)?\w{2,3})/"
		to="https://$1.google.$2/" />

	<!--	Subdomains that only exist/work on .com
							-->
	<rule from="^http://(apis|appengine|books|calendar|cbks0|chat|checkout|chrome|clients[12]|code|[\w-]+\.corp|developers|dl|docs\d?|drive|encrypted|encrypted-tbn[123]|feedburner|fiber|fonts|gg|glass||health|helpouts|history|(?:hosted)?talkgadget|investor|lh\d|(?:chatenabled\.)?mail|pack|pki|play|plus(?:\.sandbox)?|plusone|productforums|profiles|safebrowsing-cache|cert-test\.sandbox|sb-ssl|script|security|servicessites|sites|spreadsheets\d?|support|talk|tools)\.google\.com/"
		to="https://$1.google.com/" />

	<exclusion pattern="^http://clients[0-9]\.google\.com/ocsp"/>

	<rule from="^http://scholar\.google\.((?:com?\.)?\w{2,3})/intl/"
		to="https://www.google.$1/intl/" />

	<rule from="^http://(?:encrypted-)?tbn2\.google\.com/"
		to="https://encrypted-tbn2.google.com/" />


	<rule from="^http://knoll?\.google\.com/"
		to="https://knol.google.com/" />


	<rule from="^http://news\.google\.(?:com?\.)?\w{2,3}/(?:$|news|newshp)"
		to="https://www.google.com/news" />

	<rule from="^http://trends\.google\.com/"
		 to="https://www.google.com/trends" />


	<rule from="^http://([^/:@]+\.)?googlecode\.com/"
		 to="https://$1googlecode.com/" />

	<rule from="^http://([^\./]\.)?googlesource\.com/"
		to="https://$1googlesource.com/" />


	<rule from="^http://partner\.googleadservices\.com/"
		 to="https://partner.googleadservices.com/" />

	<rule from="^http://(pagead2|tpc)\.googlesyndication\.com/"
		 to="https://$1.googlesyndication.com/" />

	<!--	!www doesn't exist.
					-->
	<rule from="^http://www\.googletagservices\.com/tag/js/"
		to="https://www.googletagservices.com/tag/js/" />


	<rule from="^http://([^@:\./]+)\.googleusercontent\.com/"
		to="https://$1.googleusercontent.com/" />
	

</ruleset>`

func testRuleset(t *testing.T, r interface {
	Apply(original string, host string) (bool, string, error)
}, original string, expected string) {
	applied, newUrl, err := r.Apply(original, "")
	if err != nil {
		t.Errorf("transforming %s, got error: %s", original, err)
		return
	}
	if len(expected) <= 0 {
		if applied {
			t.Errorf("transforming %s, expected no application, got %s", original, newUrl)
		}
		return
	}

	if !applied {
		t.Errorf("transforming %s, expected %s, got no application", original, expected)
		return
	}

	if newUrl != expected {
		t.Errorf("transforming %s, expected %s, got %s", original, expected, newUrl)
		return
	}
}

func TestParseSingle(t *testing.T) {
	t.Parallel()

	testFileReader := bytes.NewReader([]byte(testFile))
	ruleset, err := ParseRuleFile(testFileReader)
	if err != nil {
		t.Fatalf("Error parsing ruleset: %s", err)
		return
	}

	testRuleset(t, ruleset, "http://news.google.com/news/rtc?ncl=dUgdV9mUeyEyqIM5jKPAAkeoCXAiM&authuser=0&ned=us&topic=h&siidp=e56f61287b9a51478c8d1d287dabd193c114", "https://www.google.com/news/rtc?ncl=dUgdV9mUeyEyqIM5jKPAAkeoCXAiM&authuser=0&ned=us&topic=h&siidp=e56f61287b9a51478c8d1d287dabd193c114")
	testRuleset(t, ruleset, "http://lh5.googleusercontent.com/-Muu7kIym7Mo/Uv1BoZia10I/AAAAAAAAB0Q/MMh8F1DVumQ/w630-h709-no/camera_20140213160225672.jpg", "https://lh5.googleusercontent.com/-Muu7kIym7Mo/Uv1BoZia10I/AAAAAAAAB0Q/MMh8F1DVumQ/w630-h709-no/camera_20140213160225672.jpg")
}
