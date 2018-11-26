function showResult(){
	var first = (page - 1) * 40;
	var end = page * 40;
	var result;
	result = "<section id=\"js-react-search-mid\">\n<div class=\"_2xrGgcY\">";
	for(var i = first; i < pid.length && i < end; i++)
	{
		result += "<div class=\"_7IVJuWZ\"><figure class=\"gmzooM4\" style=\"width: 100px; max-height: 288px;\">";
		
		//output picture box
		/*var pAddress = "https://www.pixiv.net/member_illust.php?mode=medium&illust_id=" + pid[i];
		result += "<div class=\"_1NxDA4N\">";
		result += "<a href=\"" + pAddress + "\" rel=\"noopener\" class=\"bBzsEVG\">";
		result += "<img alt=\"id:" + pid[i] + "\" border=1px class=\"_1QvROXv lazyloaded\" height=150 width=100 src=\"" + pAddress + "\">";//small picture
		result += "</a>";
		result += "</div>";*/
		
		//output text box
		var aAddress = "https://www.pixiv.net/member.php?id=" + aid[i];
		result += "<figcaption class=\"_1-dF98p\"><ul><li class=\"_1Q-G7T5\">";
		result += '<a href="https://www.pixiv.net/member_illust.php?mode=medium&amp;illust_id=' + pid[i] + '" title="' + pname[i] + '">' + pname[i] + '</a>';
		result += "</li>";
		result += "<li>";
		result += '<a href="https://www.pixiv.net/member_illust.php?id=' + aid[i] + '" target="_blank" title="' + aname[i]
		+ '" class="js-click-trackable ui-profile-popup _3mThRAs" data-click-category="recommend 20130415-0531" data-click-action="ClickToMember" data-click-label="" data-user_id="' + aid[i]
		+ '" data-user_name="' + aname[i] + '"><span class="_3UHUppl">' + aname[i] + '</span></a>';
		result += "</li><li>";
		result += "like:" + like[i];
		result += "</li></ul></figcaption>";
		
		result += "</figure></div>";
		document.getElementById("result").innerHTML = result;
	}
	result += "</div></section>";
}

window.addEventListener("load", showResult, false);
