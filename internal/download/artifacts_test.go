package download

import "testing"

const (
	shaA = "5c3065c261fb3ba82032fee6803f247198fec221"
	shaB = "8cbec36a01aedf659430869cddaee6ac1c949d08"
	shaC = "0d23e0c9b9fc382a8be4b875cbde71775ccf89d3"
)

// The current docs.fivem.net page: absolute links, build list in __NEXT_DATA__,
// recommended build repeated inside the latest list.
const docsPage = `<html><body>
<a href="https://runtime.fivem.net/artifacts/fivem/build_server_windows/master/38435-` + shaA + `/server.7z">Download</a>
<script id="__NEXT_DATA__" type="application/json">{"props":{"pageProps":{"legacy":{
"recommended":{
 "windows":[{"displayName":"server.7z","subtitle":"build 38435","downloadURL":"https://runtime.fivem.net/artifacts/fivem/build_server_windows/master/38435-` + shaA + `/server.7z"}],
 "linux":[{"displayName":"fx.tar.xz","subtitle":"build 38435","downloadURL":"https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/38435-` + shaA + `/fx.tar.xz"}]},
"latest":{
 "windows":[
  {"downloadURL":"https://runtime.fivem.net/artifacts/fivem/build_server_windows/master/38781-` + shaB + `/server.7z"},
  {"downloadURL":"https://runtime.fivem.net/artifacts/fivem/build_server_windows/master/38758-` + shaC + `/server.7z"},
  {"downloadURL":"https://runtime.fivem.net/artifacts/fivem/build_server_windows/master/38435-` + shaA + `/server.7z"}],
 "linux":[
  {"downloadURL":"https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/38781-` + shaB + `/fx.tar.xz"}]}
}}}}</script></body></html>`

// The old runtime.fivem.net directory listing: relative links, recommended in text.
const legacyListing = `<html><body>
<a href="./38435-` + shaA + `/server.7z">LATEST RECOMMENDED (38435)</a>
<a href="./38781-` + shaB + `/server.7z">38781</a>
<a href="./38435-` + shaA + `/server.7z">38435</a>
</body></html>`

func TestParseBuildsDocsPageWindows(t *testing.T) {
	builds, err := parseBuilds(docsPage, "windows", "server.7z")
	if err != nil {
		t.Fatal(err)
	}

	want := []int{38435, 38781, 38758}
	if len(builds) != len(want) {
		t.Fatalf("got %d builds, want %d: %+v", len(builds), len(want), builds)
	}
	for i, n := range want {
		if builds[i].Number != n {
			t.Errorf("builds[%d] = %d, want %d", i, builds[i].Number, n)
		}
	}
	if !builds[0].Recommended || builds[1].Recommended || builds[2].Recommended {
		t.Errorf("only build 38435 should be recommended: %+v", builds)
	}
	if builds[0].Hash != "38435-"+shaA {
		t.Errorf("hash = %q", builds[0].Hash)
	}
}

func TestParseBuildsDocsPageLinux(t *testing.T) {
	builds, err := parseBuilds(docsPage, "linux", "fx.tar.xz")
	if err != nil {
		t.Fatal(err)
	}
	if len(builds) != 2 || builds[0].Number != 38435 || builds[1].Number != 38781 {
		t.Fatalf("unexpected linux builds: %+v", builds)
	}
	if !builds[0].Recommended {
		t.Errorf("38435 should be recommended")
	}
}

func TestParseBuildsLegacyListingFallback(t *testing.T) {
	builds, err := parseBuilds(legacyListing, "windows", "server.7z")
	if err != nil {
		t.Fatal(err)
	}
	if len(builds) != 2 || builds[0].Number != 38435 || builds[1].Number != 38781 {
		t.Fatalf("unexpected builds: %+v", builds)
	}
	if !builds[0].Recommended {
		t.Errorf("38435 should be recommended")
	}
}

func TestParseBuildsNone(t *testing.T) {
	if _, err := parseBuilds("<html></html>", "windows", "server.7z"); err == nil {
		t.Fatal("expected error for a page with no builds")
	}
}
