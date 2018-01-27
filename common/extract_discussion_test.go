package common

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractDiscussion(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(commentTreeHTML))
	if err != nil {
		t.Fatal(err)
	}
	threads := ExtractDiscussion(doc.Selection)
	if expected, actual := 2, len(threads); actual != expected {
		t.Fatalf("Expected len(threads)=%v but actual=%v", expected, actual)
	}
	if expected, actual := 4, threads.Len(); actual != expected {
		t.Fatalf("Expected threads.Len()=%v but actual=%v", expected, actual)
	}
	if expected, actual := 2, threads[0].N; actual != expected {
		t.Fatalf("Expected threads[0].N=%v but actual=%v", expected, actual)
	}
	if expected, actual := 2, threads[1].N; actual != expected {
		t.Fatalf("Expected threads[0].N=%v but actual=%v", expected, actual)
	}
}

func TestExtractComment(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(singleCommentHTML))
	if err != nil {
		t.Fatal(err)
	}
	comment := extractComment(doc.Selection)
	if comment == nil {
		t.Fatal("Expected comment but got nil")
	}
	if expected, actual := 80, comment.Width; actual != expected {
		t.Fatalf("Expected comment.Width=%v but actual=%v", expected, actual)
	}
	// t.Logf("c=%# v", *comment)
}

const singleCommentHTML = `
<table>
        <tr class='athing comtr ' id='18919099'><td>
            <table border='0'>  <tr>    <td class='ind'><img src="s.gif" height="1" width="80"></td><td valign="top" class="votelinks"><center><a id='up_18919099' onclick='return vote(event, this, "up")' href='vote?id=18919099&amp;how=up&amp;auth=5c4efd1e01050a032cb0bf3d6635e696d4293fb8&amp;goto=item%3Fid%3D18914411#18919099'><div class='votearrow' title='upvote'></div></a></center></td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=nraynaud" class="hnuser">nraynaud</a> <span class="age"><a href="item?id=18919099">6 days ago</a></span> <span id="unv_18919099"></span><span class="par"></span> <a class="togg" n="52" href="javascript:void(0)" onclick="return toggle(event, 18919099)"></a>          <span class='storyon'></span>
                  </span></div><br><div class="comment">
                  <span class="commtext c00">This is not the voter’s fault. In a functioning democracy, there is a system to approve a text by popular vote instead of the usual parlement. So the question at such a vote should always be: “do you approve the proposed text?”, then the text becomes law.<p>By organizing an abstract opinion poll, the government  just prepared a shit show, there no way to make any significant chunk of the population happy, everybody had a different view of brexit.</span>
              <div class='reply'>        <p><font size="1">
                      <u><a href="reply?id=18919099&amp;goto=item%3Fid%3D18914411%2318919099">reply</a></u>
                  </font>
      </div></div></td></tr>
      </table></td></tr>
</table>
`

const commentTreeHTML = `
<html op="item"><head><meta name="referrer" content="origin"><meta name="viewport" content="width=device-width, initial-scale=1.0"><link rel="stylesheet" type="text/css" href="news.css?tR69YSVmkcPOlJWQ6QcH">
            <link rel="shortcut icon" href="favicon.ico">
        <title>1) https:&#x2F;&#x2F;www.ctvnews.ca&#x2F;business&#x2F;iran-signs-5b-gas-deal-with-france-s-total-ch... | Hacker News</title></head><body><center><table id="hnmain" border="0" cellpadding="0" cellspacing="0" width="85%" bgcolor="#f6f6ef">
        <tr><td bgcolor="#cccccc"><table border="0" cellpadding="0" cellspacing="0" width="100%" style="padding:2px"><tr><td style="width:18px;padding-right:4px"><a href="https://news.ycombinator.com"><img src="y18.gif" width="18" height="18" style="border:1px white solid;"></a></td>
                  <td style="line-height:12pt; height:10px;"><span class="pagetop"><b class="hnname"><a href="news">Hacker News</a></b>
              <a href="newest">new</a> | <a href="threads?id=jaytaylor">threads</a> | <a href="newcomments">comments</a> | <a href="ask">ask</a> | <a href="show">show</a> | <a href="jobs">jobs</a> | <a href="submit">submit</a>            </span></td><td style="text-align:right;padding-right:4px;"><span class="pagetop">
                              <a id='me' href="user?id=jaytaylor">jaytaylor</a>                (3868) |
                <a id='logout' href="logout?auth=a977c3ce2e5cb1b40f916fae46db1179c3f9a668&amp;goto=item%3Fid%3D18927109">logout</a>                          </span></td>
              </tr></table></td></tr>
<tr id="pagespace" title="1) https:&#x2F;&#x2F;www.ctvnews.ca&#x2F;business&#x2F;iran-signs-5b-gas-deal-with-france-s-total-ch..." style="height:10px"></tr><tr><td><table class="fatitem" border="0">
    <tr class='athing' id='18927109'>    <td class='ind'></td><td valign="top" class="votelinks"><center><a id='up_18927109' onclick='return vote(event, this, "up")' href='vote?id=18927109&amp;how=up&amp;auth=0efdd19736e2ffb2b985bd0bf6475808794ec1bb&amp;goto=item%3Fid%3D18927109#18927109'><div class='votearrow' title='upvote'></div></a></center></td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=candiodari" class="hnuser">candiodari</a> <span class="age"><a href="item?id=18927109">12 days ago</a></span> <span id="unv_18927109"></span><span class="par"> | <a href="item?id=18926647">parent</a></span> | <a href="flag?id=18927109&amp;auth=0efdd19736e2ffb2b985bd0bf6475808794ec1bb&amp;goto=item%3Fid%3D18927109">flag</a> | <a href="fave?id=18927109&amp;auth=0efdd19736e2ffb2b985bd0bf6475808794ec1bb">favorite</a>          <span class='storyon'> | on: <a href="item?id=18914411">Brexit Deal Fails in Parliament</a></span>
                  </span></div><br><div class="comment">
                  <span class="commtext c00">1) <a href="https:&#x2F;&#x2F;www.ctvnews.ca&#x2F;business&#x2F;iran-signs-5b-gas-deal-with-france-s-total-chinese-oil-firm-1.3486786" rel="nofollow">https:&#x2F;&#x2F;www.ctvnews.ca&#x2F;business&#x2F;iran-signs-5b-gas-deal-with-...</a><p>Now you could say there&#x27;s some excuse here. The oil in question technically never enters the EU, so it &quot;doesn&#x27;t have to&quot; follow EU legislation (which prohibits dealing commercially with Iran). I&#x27;m sure if necessary something like that will be pointed out. Maybe it&#x27;s also &quot;not&quot; Total selling the oil, but obviously they get part of the profit. And I&#x27;m sure there&#x27;s some reason &quot;Iran is not involved at all&quot;.<p>Also worth pointing out: the US has in the meantime successfully forced Total, over the LOUD protest of both French and German governments to abandon this deal.<p>2) <a href="https:&#x2F;&#x2F;www.expatica.com&#x2F;ch&#x2F;moving&#x2F;visas&#x2F;guide-for-eu-efta-citizens-and-relatives-moving-to-switzerland-443220&#x2F;#romanians" rel="nofollow">https:&#x2F;&#x2F;www.expatica.com&#x2F;ch&#x2F;moving&#x2F;visas&#x2F;guide-for-eu-efta-c...</a><p>You can see the rules about employing Romanian and Bulgarian citizens (which are EU citizens) in Switzerland are <i>very</i> different from employing, say, French or German citizens. Note that this is one of the things the UK now tells Britain is non-negotiable. Somehow in practice another (labour related) trade deal the EU ... has negotiated it.<p>Of course, there are historical reasons for this, but of course that&#x27;s true for all trade deals.</span>
              <div class='reply'></div></div></td></tr>
          <tr style="height:10px"></tr><tr><td colspan="2"></td><td>
          <form method="post" action="comment"><input type="hidden" name="parent" value="18927109"><input type="hidden" name="goto" value="item?id=18927109"><input type="hidden" name="hmac" value="fddd74e9277db12228d6fd668b3be8e445130303"><textarea name="text" rows="6" cols="60"></textarea>
                <br><br><input type="submit" value="reply"></form>
      </td></tr>
  </table><br><br>
  <table border="0" class='comment-tree'>
            <tr class='athing comtr ' id='18929547'><td>
            <table border='0'>  <tr>    <td class='ind'><img src="s.gif" height="1" width="0"></td><td valign="top" class="votelinks"><center><a id='up_18929547' onclick='return vote(event, this, "up")' href='vote?id=18929547&amp;how=up&amp;auth=9df5c9348c54a7598bb0b1b81bd569a5be158ea6&amp;goto=item%3Fid%3D18927109#18929547'><div class='votearrow' title='upvote'></div></a></center></td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=purple_ducks" class="hnuser">purple_ducks</a> <span class="age"><a href="item?id=18929547">12 days ago</a></span> <span id="unv_18929547"></span><span class="par"></span> <a class="togg" n="2" href="javascript:void(0)" onclick="return toggle(event, 18929547)"></a>          <span class='storyon'></span>
                  </span></div><br><div class="comment">
                  <span class="commtext c00">re: 1) That article is about a French multinational oil company entering a commercial agreement with Iran. This doesn&#x27;t relate to an EU member state doing a &quot;side&quot; trade treaty with a non-EU member.<p>Any other EU company could have done the same.<p>re: 2) Seems like that provision was in the agreement Switzerland made with the EU:<p>from: <a href="https:&#x2F;&#x2F;www.swissinfo.ch&#x2F;eng&#x2F;business&#x2F;freedom-of-movement_switzerland-prolongs-immigration-limits-for-bulgarians-and-romanians&#x2F;44057938" rel="nofollow">https:&#x2F;&#x2F;www.swissinfo.ch&#x2F;eng&#x2F;business&#x2F;freedom-of-movement_sw...</a><p>&gt; Switzerland first activated such a safeguard clause – a controversial instrument of its complex dealings with the EU – in 2012, to limit the number of citizens arriving from certain new member countries who joined the EU in 2004..</span>
              <div class='reply'>        <p><font size="1">
                      <u><a href="reply?id=18929547&amp;goto=item%3Fid%3D18927109%2318929547">reply</a></u>
                  </font>
      </div></div></td></tr>
      </table></td></tr>
        <tr class='athing comtr ' id='18929689'><td>
            <table border='0'>  <tr>    <td class='ind'><img src="s.gif" height="1" width="40"></td><td valign="top" class="votelinks"><center><a id='up_18929689' onclick='return vote(event, this, "up")' href='vote?id=18929689&amp;how=up&amp;auth=dabb5b40f0c8dc21c63a04979bea664c90f69228&amp;goto=item%3Fid%3D18927109#18929689'><div class='votearrow' title='upvote'></div></a></center></td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=inflagranti" class="hnuser">inflagranti</a> <span class="age"><a href="item?id=18929689">12 days ago</a></span> <span id="unv_18929689"></span><span class="par"></span> <a class="togg" n="1" href="javascript:void(0)" onclick="return toggle(event, 18929689)"></a>          <span class='storyon'></span>
                  </span></div><br><div class="comment">
                  <span class="commtext c00">From my understanding the Switzerland negotiated clause re Romania and other late joiners was timebound, which is why it was acceptable.</span>
              <div class='reply'>        <p><font size="1">
                      <u><a href="reply?id=18929689&amp;goto=item%3Fid%3D18927109%2318929689">reply</a></u>
                  </font>
      </div></div></td></tr>
      </table></td></tr>
                <tr class='athing comtr ' id='18928080'><td>
            <table border='0'>  <tr>    <td class='ind'><img src="s.gif" height="1" width="0"></td><td valign="top" class="votelinks"><center><a id='up_18928080' onclick='return vote(event, this, "up")' href='vote?id=18928080&amp;how=up&amp;auth=4d895ed10c673e86a55e4123154dd674fba7acdb&amp;goto=item%3Fid%3D18927109#18928080'><div class='votearrow' title='upvote'></div></a></center></td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=lightgreen" class="hnuser">lightgreen</a> <span class="age"><a href="item?id=18928080">12 days ago</a></span> <span id="unv_18928080"></span><span class="par"></span> <a class="togg" n="2" href="javascript:void(0)" onclick="return toggle(event, 18928080)"></a>          <span class='storyon'></span>
                  </span></div><br><div class="comment">
                  <span class="commtext c00">2) employment is not trade. So neither of your examples prove that EU members can negotiate side trade deal.</span>
              <div class='reply'>        <p><font size="1">
                      <u><a href="reply?id=18928080&amp;goto=item%3Fid%3D18927109%2318928080">reply</a></u>
                  </font>
      </div></div></td></tr>
      </table></td></tr>
        <tr class='athing comtr ' id='18929385'><td>
            <table border='0'>  <tr>    <td class='ind'><img src="s.gif" height="1" width="40"></td><td valign="top" class="votelinks"><center><a id='up_18929385' onclick='return vote(event, this, "up")' href='vote?id=18929385&amp;how=up&amp;auth=6d7bd62e5bbe2cb3760b29fabe118636e241a1f6&amp;goto=item%3Fid%3D18927109#18929385'><div class='votearrow' title='upvote'></div></a></center></td><td class="default"><div style="margin-top:2px; margin-bottom:-10px;"><span class="comhead">
          <a href="user?id=candiodari" class="hnuser">candiodari</a> <span class="age"><a href="item?id=18929385">12 days ago</a></span> <span id="unv_18929385"></span><span class="par"></span> <a class="togg" n="1" href="javascript:void(0)" onclick="return toggle(event, 18929385)"></a>          <span class='storyon'></span>
                  </span></div><br><div class="comment">
                  <span class="commtext c00">You might want to read exactly on what point does the fight between the EU and UK is mostly fought ... Yep: labour and movement of people is one of the important points ...</span>
              <div class='reply'>        <p><font size="1">
                      <u><a href="reply?id=18929385&amp;goto=item%3Fid%3D18927109%2318929385">reply</a></u>
                  </font>
      </div></div></td></tr>
      </table></td></tr>
              </table>
      <br><br>
  </td></tr>
<tr><td><img src="s.gif" height="10" width="0"><table width="100%" cellspacing="0" cellpadding="1"><tr><td bgcolor="#cccccc"></td></tr></table><br><center><a href="https://www.ycombinator.com/apply/">
        Applications are open for YC Summer 2019
      </a></center><br><center><span class="yclinks"><a href="newsguidelines.html">Guidelines</a>
        | <a href="newsfaq.html">FAQ</a>
        | <a href="mailto:hn@ycombinator.com">Support</a>
        | <a href="https://github.com/HackerNews/API">API</a>
        | <a href="security.html">Security</a>
        | <a href="lists">Lists</a>
        | <a href="bookmarklet.html" rel="nofollow">Bookmarklet</a>
        | <a href="http://www.ycombinator.com/legal/">Legal</a>
        | <a href="http://www.ycombinator.com/apply/">Apply to YC</a>
        | <a href="mailto:hn@ycombinator.com">Contact</a></span><br><br><form method="get" action="//hn.algolia.com/">Search:
          <input type="text" name="q" value="" size="17" autocorrect="off" spellcheck="false" autocapitalize="off" autocomplete="false"></form>
            </center></td></tr>
      </table></center></body><script type='text/javascript' src='hn.js?tR69YSVmkcPOlJWQ6QcH'></script>
  </html>
`
