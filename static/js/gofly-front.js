var GOFLY={
    GOFLY_URL:"",
    GOFLY_KEFU_ID:"",
    GOFLY_ENT:"",
    GOFLY_LANG:"en",
    GOFLY_EXTRA: {},
    GOFLY_AUTO_OPEN:true,//是否自动打开
    GOFLY_SHOW_TYPES:1,//展示样式，1：普通右下角，2：圆形icon
    GOFLY_AUTO_SHOW:false,
    GOFLY_WITHOUT_BOX:false,
};
GOFLY.launchButtonFlag=false;
GOFLY.titleTimer=0;
GOFLY.titleNum=0;
GOFLY.noticeTimer=null;
GOFLY.originTitle=document.title;
GOFLY.chatPageTitle="GOFLY";
GOFLY.kefuName="";
GOFLY.kefuAvator="";
GOFLY.offLine=false;
GOFLY.TEXT={
    "cn":{
        "online_notice":"和我们在线交谈",
        "offline_notice":"离线请留言",
    },
    "en":{
        "online_notice":"we are online,chat with us",
        "offline_notice":"we are offline",
    },
};
GOFLY.init=function(config){
    var _this=this;
    if(typeof config=="undefined"){
        return;
    }

    if (typeof config.GOFLY_URL!="undefined"){
        this.GOFLY_URL=config.GOFLY_URL.replace(/([\w\W]+)\/$/,"$1");
    }
    this.dynamicLoadCss(this.GOFLY_URL+"/static/css/gofly-front.css?v="+Date.now());
    this.dynamicLoadCss(this.GOFLY_URL+"/static/css/layui/css/layui.css?v="+Date.now());
    if (typeof config.GOFLY_KEFU_ID!="undefined"){
        this.GOFLY_KEFU_ID=config.GOFLY_KEFU_ID;
    }
    if (typeof config.GOFLY_ENT!="undefined"){
        this.GOFLY_ENT=config.GOFLY_ENT;
    }
    if (typeof config.GOFLY_EXTRA!="undefined"){
        this.GOFLY_EXTRA=config.GOFLY_EXTRA;
    }
    if (typeof config.GOFLY_AUTO_OPEN!="undefined"){
        this.GOFLY_AUTO_OPEN=config.GOFLY_AUTO_OPEN;
    }
    if (typeof config.GOFLY_SHOW_TYPES!="undefined"){
        this.GOFLY_SHOW_TYPES=config.GOFLY_SHOW_TYPES;
    }
    if (typeof config.GOFLY_WITHOUT_BOX!="undefined"){
        this.GOFLY_WITHOUT_BOX=config.GOFLY_WITHOUT_BOX;
    }
    var refer=document.referrer?document.referrer:"无";
    this.GOFLY_EXTRA.refer=refer;
    this.GOFLY_EXTRA.host=document.location.href;
    this.GOFLY_EXTRA=JSON.stringify(_this.GOFLY_EXTRA);

    this.dynamicLoadJs(this.GOFLY_URL+"/static/js/functions.js?v=1",function(){
        if (typeof config.GOFLY_LANG!="undefined"){
            _this.GOFLY_LANG=config.GOFLY_LANG;
        }else{
            _this.GOFLY_LANG=checkLang();
        }
        _this.GOFLY_EXTRA=utf8ToB64(_this.GOFLY_EXTRA);
    });
    if (typeof $!="function"){
        this.dynamicLoadJs("https://cdn.staticfile.org/jquery/3.6.0/jquery.min.js",function () {
            _this.dynamicLoadJs("https://cdn.staticfile.org/layer/3.4.0/layer.min.js",function () {
                _this.jsCallBack();
            });
        });
    }else{
        this.dynamicLoadJs("https://cdn.staticfile.org/layer/3.4.0/layer.min.js",function () {
            _this.jsCallBack();
        });
    }
    _this.addEventlisten();
}
GOFLY.jsCallBack=function(){
    var _this=this;
    if(_this.isMobile()){
        _this.GOFLY_SHOW_TYPES=2;
        _this.GOFLY_AUTO_OPEN=false;
    }
    _this.getNotice(function(welcomeInfo){
        //展示的样式
        switch(_this.GOFLY_SHOW_TYPES){
            case 1:
                _this.showPcTips(welcomeInfo);
                break;
            case 2:
                _this.showCircleIcon(welcomeInfo);
                break;
            case 3:
                _this.showLineTips(welcomeInfo);
                break;
            default:

        }
        _this.addClickEvent();
    });
}
//pc端的样式
GOFLY.showPcTips=function(welcomeInfo){
    var _this=this;
    _this.kefuAvator=welcomeInfo.avatar;
    _this.kefuName=welcomeInfo.username;
    _this.chatPageTitle=welcomeInfo.username;
    _this.offLine=welcomeInfo.all_offline;
    if(welcomeInfo.all_offline){
        _this.GOFLY_AUTO_OPEN=false;
        var userInfo="<img style='margin-top: 5px;' class='flyAvatar' src='"+_this.GOFLY_URL+welcomeInfo.avatar+"'/> <span class='flyUsername'>"+GOFLY.TEXT[_this.GOFLY_LANG].offline_notice+"</span>"
    }else{
        var userInfo="<img style='margin-top: 5px;' class='flyAvatar' src='"+_this.GOFLY_URL+welcomeInfo.avatar+"'/> <span class='flyUsername'>"+GOFLY.TEXT[_this.GOFLY_LANG].online_notice+"</span>"
    }
    //自动展开
    if(_this.GOFLY_AUTO_OPEN&&_this.isIE()<=0){
        setTimeout(function () {
            _this.showKefu();
        },5000);
    }
    var html="<div class='launchButtonBox'>" +
        '<div id="launchButton" class="launchButton">' +
        '<div id="launchIcon" class="launchIcon">1</div> ' +
        '<div class="launchButtonText">'+userInfo+'</div></div>' +
        '<div id="launchButtonNotice" class="launchButtonNotice"></div>' +
        '</div>';

    $('body').append(html);
    if(_this.GOFLY_AUTO_OPEN){
        return;
    }
    if (!welcomeInfo.welcome) return;
    var msgs = welcomeInfo.welcome;
    var len=msgs.length;
    var i=0;
    if(len<=0) {
        return;
    }
    var delaySecond=0;
    for(let i in msgs){
        var msg=msgs[i];
        if(msg.delay_second){
            delaySecond+=msg.delay_second;
        }else{
            delaySecond+=4;
        }
        var timer =  setTimeout(function (msg) {
            msg.content = replaceSpecialTag(msg.content);



            var welcomeHtml="<div class='flyUser'><img class='flyAvatar' src='"+_this.GOFLY_URL+msg.avator+"'/> <span class='flyUsername'>"+msg.name+"</span>" +
                "<span id='launchNoticeClose' class='flyClose'>×</span>" +
                "</div>";
            welcomeHtml+="<div id='launchNoticeContent'>"+replaceSpecialTag(msg.content,_this.GOFLY_URL)+"</div>";

            var obj=$("#launchButtonNotice");
            if(obj){
                obj.html(welcomeHtml).fadeIn();
                setTimeout(function (obj) {
                    obj.fadeOut();
                },3000,obj);
            }

            i++;
            $("#launchIcon").text(i).fadeIn();
        },1000*delaySecond,msg);
    }
}
//pc端的第二种样式
GOFLY.showLineTips=function(welcomeInfo){
    var _this=this;
    _this.kefuAvator=welcomeInfo.avatar;
    _this.kefuName=welcomeInfo.username;
    _this.chatPageTitle=welcomeInfo.username;
    _this.offLine=welcomeInfo.all_offline;
    //自动展开
    if(_this.GOFLY_AUTO_OPEN&&_this.isIE()<=0){
        setTimeout(function () {
            _this.showKefu();
        },10000);
    }
   str=`
<div class="lineBox">
    <div class="lineItem" onclick="javascript:GOFLY.showPanel();">
        <i class="layui-icon">&#xe606;</i>
    </div>
    <div class="lineItem">
        <i class="layui-icon">&#xe677;</i>
        <div class="lineTip lineWechat">
            <img class="lineWechat" src="h"/>
        </div>
    </div>
    <div class="lineItem">
        <i class="layui-icon">&#xe676;</i>
        <div class="lineTip">
            QQ:630892807
        </div>
    </div>
    <div class="lineItem">
        <i class="layui-icon">&#xe626;</i>
        <div class="lineTip">
            QQ:630892807
        </div>
    </div>
    <div class="lineItem lineTop" id="launchTopButton">
        <i class="layui-icon">&#xe604;</i>
    </div>
</div>`
    $('body').append(str);
}
//圆形样式
GOFLY.showCircleIcon=function(welcomeInfo){
    this.offLine=welcomeInfo.all_offline;
    if(welcomeInfo.all_offline){
        var imgUrl=GOFLY.GOFLY_URL+"/static/images/iconchat.png";
        var tipText=GOFLY.TEXT[this.GOFLY_LANG].offline_notice
        var imgHtml="<img class='flySimpleDefaultImg' src='"+imgUrl+"'/> ";
    }else{
        var imgUrl=GOFLY.GOFLY_URL+welcomeInfo.avatar;
        var tipText=GOFLY.TEXT[this.GOFLY_LANG].online_notice
        var imgHtml="<img class='flySimpleUserImg' src='"+imgUrl+"'/> ";
    }
    var _this=this;
    _this.kefuAvator=welcomeInfo.avatar;
    _this.kefuName=welcomeInfo.username;
    _this.chatPageTitle=welcomeInfo.username;
    var html="<div class='flySimpleIconBox'>" +
        "<div class='flySimpleIcon'>" +
        imgHtml+
        '</div>' +
        "<div class='flySimpleIconTip'> " +tipText+
        "<div class='flyClose'>×</div>" +
        "</div>" +
        '</div>';
    $('body').append(html);
    setTimeout(function () {
        $(".flySimpleIconTip").fadeIn();
        setTimeout(function () {
            $(".flySimpleIconTip").fadeOut();
        },5000);
    },5000);
    window.addEventListener('message',function(e){
        var msg=e.data;
        if(msg.type=="message"){
            $(".flySimpleIconTip").html(replaceSpecialTag(msg.data.content,_this.GOFLY_URL)+"<div class='flyClose'>×</div>").show();
            setTimeout(function () {
                $(".flySimpleIconTip").fadeOut();
            },5000);
        }
        if(msg.type=="force_close"){
            layer.close(layer.index);
        }
    });

    //自动展开
    if(_this.GOFLY_AUTO_OPEN&&_this.isIE()<=0){
        setTimeout(function () {
            _this.showKefu();
        },5000);
    }
}
GOFLY.addClickEvent=function(){
    var _this=this;
    $("#launchButton").on("click",function() {
        if(_this.launchButtonFlag){
            return;
        }
        _this.GOFLY_AUTO_SHOW=true;
        _this.showKefu();
        $("#launchIcon").text(0).hide();
    });

    $("body").on("click","#launchNoticeClose",function() {
        $("#launchButtonNotice").fadeOut();
    });
    $("body").on("click",".flySimpleIconTip",function() {
        $(".flySimpleIconTip").fadeOut();
    });
    $("body").on("click","#launchTopButton",function() {
        $('body,html').scrollTop(0);
    });
    $("body").on("mouseover mouseout",".lineItem",function(event) {
        if(event.type == "mouseover"){
            //鼠标悬浮
            $(this).find(".lineTip").show();
        }else if(event.type == "mouseout"){
            //鼠标离开
            $(".lineTip").hide();
        }
    });
    $("body").click(function () {
        clearTimeout(_this.titleTimer);
        document.title = _this.originTitle;
        //剪贴板
        try{
            var selecter = window.getSelection().toString();
            if (selecter != null && selecter.trim() != ""){
                var str=selecter.trim().substr(0,20);
                _this.postMessageToIframe(str);
            }
        } catch (err){
            var selecter = document.selection.createRange();
            var s = selecter.text;
            if (s != null && s.trim() != ""){
                var str=s.trim().substr(0,20);
                _this.postMessageToIframe(str);
            }
        }
    });
    var ms= 1000*2;
    var lastClick = Date.now() - ms;
    $("a,div,p,li").mouseover(function(){
        if (Date.now() - lastClick >= ms) {
            lastClick = Date.now();
            var msg=$(this).text().trim().substr(0,20);
            _this.postMessageToIframe(msg);
        }
    });
    $("body").on("click",".flySimpleIcon",function() {
        _this.showPanel();
    });
}

GOFLY.postMessageToIframe=function(str){
    var msg={}
    msg.type='inputing_message';
    msg.content=str;
    this.postToIframe(msg);
}
GOFLY.postToIframe=function(messageObj){
    var obj=document.getElementById('layui-layer-iframe19911116');
    if(!obj||!messageObj){
        return;
    }
    document.getElementById('layui-layer-iframe19911116').contentWindow.postMessage(messageObj, "*");
}
GOFLY.addEventlisten=function(){
    var _this=this;
    window.addEventListener('message',function(e){
        var msg=e.data;
        if(msg.type=="message"){
            clearInterval(_this.noticeTimer);
            var width=$(window).width();
            if(width>768){
                _this.flashTitle();//标题闪烁
            }
            if (_this.launchButtonFlag){
                return;
            }
            var welcomeHtml="<div class='flyUser'><img class='flyAvatar' src='"+_this.GOFLY_URL+msg.data.avator+"'/> <span class='flyUsername'>"+msg.data.name+"</span>" +
                "<span id='launchNoticeClose' class='flyClose'>×</span>" +
                "</div>";
            welcomeHtml+="<div id='launchNoticeContent'>"+replaceSpecialTag(msg.data.content,_this.GOFLY_URL)+"</div>";
            var obj=$("#launchButtonNotice");
            if(obj){
                obj.html(welcomeHtml).fadeIn();
                setTimeout(function (obj) {
                    obj.fadeOut();
                },3000,obj);
            }
            var news=$("#launchIcon").text();
            $("#launchIcon").text(++news).show();
        }
        if(msg.type=="focus"){
            clearTimeout(_this.titleTimer);
            _this.titleTimer=0;
            document.title = _this.originTitle;
        }
        if(msg.type=="force_close"){
            layer.close(layer.index);
        }
    });
    window.onfocus = function () {
        clearTimeout(_this.titleTimer);
        _this.titleTimer=0;
        document.title = _this.originTitle;
        _this.postToIframe({type:"focus"});
    };
}
GOFLY.dynamicLoadCss=function(url){
    var head = document.getElementsByTagName('head')[0];
    var link = document.createElement('link');
    link.type='text/css';
    link.rel = 'stylesheet';
    link.href = url;
    head.appendChild(link);
}
GOFLY.dynamicLoadJs=function(url, callback){
    var head = document.getElementsByTagName('head')[0];
    var script = document.createElement('script');
    script.type = 'text/javascript';
    script.src = url;
    if(typeof(callback)=='function'){
        script.onload = script.onreadystatechange = function () {
            if (!this.readyState || this.readyState === "loaded" || this.readyState === "complete"){
                callback();
                script.onload = script.onreadystatechange = null;
            }
        };
    }
    head.appendChild(script);
}

//获取基础信息
GOFLY.getNotice=function(callback){
    var _this=this;
    $.get(this.GOFLY_URL+"/notice?ent_id="+this.GOFLY_ENT,function(res) {
        if(!res) return;
        var welcomeInfo=res.result;
        callback(welcomeInfo);
    });
}
GOFLY.isIE=function(){
    var userAgent = navigator.userAgent; //取得浏览器的userAgent字符串
    var isIE = userAgent.indexOf("compatible") > -1 && userAgent.indexOf("MSIE") > -1; //判断是否IE<11浏览器
    var isEdge = userAgent.indexOf("Edge") > -1 && !isIE; //判断是否IE的Edge浏览器
    var isIE11 = userAgent.indexOf('Trident') > -1 && userAgent.indexOf("rv:11.0") > -1;
    if(isIE) {
        var reIE = new RegExp("MSIE (\\d+\\.\\d+);");
        reIE.test(userAgent);
        var fIEVersion = parseFloat(RegExp["$1"]);
        if(fIEVersion == 7) {
            return 7;
        } else if(fIEVersion == 8) {
            return 8;
        } else if(fIEVersion == 9) {
            return 9;
        } else if(fIEVersion == 10) {
            return 10;
        } else {
            return 6;//IE版本<=7
        }
    } else if(isEdge) {
        return 'edge';//edge
    } else if(isIE11) {
        return 11; //IE11
    }else{
        return -1;//不是ie浏览器
    }
}
GOFLY.showPanel=function (){
    var width=$(window).width();
    this.GOFLY_AUTO_SHOW=true;
    if(this.isIE()>0){
        this.windowOpen();
        return;
    }
    if(width<768){
        this.layerOpen("100%","100%");
        return;
    }
    this.layerOpen("435px","550px");
    this.launchButtonFlag=true;
}
GOFLY.showKefu=function (){
    if (this.launchButtonFlag) return;
    var width=$(window).width();
    if(this.isIE()>0){
        this.windowOpen();
        return;
    }
    if(width<768){
        this.layerOpen("100%","100%");
        return;
    }
    this.layerOpen("450px","580px");
    this.launchButtonFlag=true;
    $(".launchButtonBox").hide();
}
GOFLY.layerOpen=function (width,height){
    if (this.launchButtonFlag) return;
    var layBox=$("#layui-layer19911116");
    if(layBox.css("display")=="none"){
        layBox.show();
        return;
    }
    var onlineStatus="<i></i>";
    if(this.offLine){
        onlineStatus="<i class='offline'></i>";
    }
    var title=`
    <div class="kfBar">
        <div class="kfBarAvator">
            <img src="`+this.GOFLY_URL+this.kefuAvator+`" class="kfBarLogo">
            `+onlineStatus+`
        </div>
        <div class="kfBarText">
            <div class="kfName">`+this.kefuName+`</div>
            <div class="kfDesc">独立部署客服系统请加微信</div>
         </div>
    </div>
    `;
    var _this=this;
    layer.index="19911115";
    layer.open({
        type: 2,
        title: title,
        skin:"kfLayer",
        closeBtn: 1, //不显示关闭按钮
        shade: 0,
        area: [width, height],
        offset: 'rb', //右下角弹出
        anim: 2,
        content: [this.GOFLY_URL+'/chatIndex?kefu_id='+this.GOFLY_KEFU_ID+'&ent_id='+this.GOFLY_ENT+'&lang='+this.GOFLY_LANG+'&refer='+window.document.title+'&url='+document.location.href+'&extra='+this.GOFLY_EXTRA , 'yes'], //iframe的url，no代表不显示滚动条
        success:function(){
            var layBox=$("#layui-layer19911116");
            _this.launchButtonFlag=true;
            if(!_this.GOFLY_WITHOUT_BOX&&_this.GOFLY_AUTO_SHOW&&layBox.css("display")=="none"){
                layBox.show();
            }
            $("#layui-layer-iframe19911116 .chatEntTitle").hide();
        },
        end: function(){
            _this.launchButtonFlag=false;
            $(".launchButtonBox").show();
        },
        cancel: function(index, layero){
            $("#layui-layer19911116").hide();
            _this.launchButtonFlag=false;
            $(".launchButtonBox").show();
            return false;
        }
    });
}
GOFLY.windowOpen=function (){
   window.open(this.GOFLY_URL+'/chatIndex?kefu_id='+this.GOFLY_KEFU_ID+'&lang='+this.GOFLY_LANG+'&refer='+window.document.title+'&ent_id='+this.GOFLY_ENT+'&extra='+this.GOFLY_EXTRA);
}
GOFLY.flashTitle=function () {
    if(this.titleTimer!=0){
        return;
    }
    this.titleTimer = setInterval("GOFLY.flashTitleFunc()", 500);
}
GOFLY.flashTitleFunc=function(){
    this.titleNum++;
    if (this.titleNum >=3) {
        this.titleNum = 1;
    }
    if (this.titleNum == 1) {
        document.title = '【】' + this.originTitle;
    }
    if (this.titleNum == 2) {
        document.title = '【new message】' + this.originTitle;
    }
}
/**
 * 判断是否是手机访问
 * @returns {boolean}
 */
GOFLY.isMobile=function () {
    if( /Mobile|Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) ) {
        return true;
    }
    return false;
}

