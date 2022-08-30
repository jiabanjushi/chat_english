new Vue({
    el: '#app',
    delimiters:["<{","}>"],
    data: {
        window:window,
        server:getWsBaseUrl()+"/ws_visitor",
        socket:null,
        socketClosed:false,
        msgList:[],
        messageContent:"",
        visitor:{},
        chatTitle:"",
        flyLang:GOFLY_LANG[LANG],
        textareaFocused:false,
        sendDisabled:false,
        face:[],
        heartsFlow:null,
    },
    methods: {
        //初始化websocket
        initConn:function() {
            this.socket = new ReconnectingWebSocket(this.server+"?visitor_id="+this.visitor.visitor_id+"&to_id="+this.visitor.to_id+"&room_id="+this.visitor.to_id);//创建Socket实例
            this.socket.debug = true;
            this.socket.timeoutInterval = 10000;//连接超时时间
            this.socket.reconnectInterval=5000;//重连间隔时间
            this.socket.maxReconnectInterval = 600000;//最大重连间隔时间
            this.socket.maxReconnectAttempts = 10;//最大重连尝试次数
            this.socket.onmessage = this.OnMessage;
            this.socket.onopen = this.OnOpen;
            this.socket.onerror = this.OnError;
            this.socket.onclose = this.OnClose;
        },
        OnOpen:function() {
            console.log("ws:onopen");
            this.chatTitle=GOFLY_LANG[LANG]['connectok'];
            this.socketClosed=false;
        },
        OnMessage:function(e) {
            console.log("ws:onmessage");
            const redata = JSON.parse(e.data);
            if (redata.type == "userOnline") {
                let msg = redata.data
                if(!msg){
                    return;
                }
                this.showTitle(msg.username);
                this.$notify({
                    dangerouslyUseHTMLString: true,
                    message: '<span style="line-height: 40px"><img style="float: left" class="el-avatar el-avatar--circle" src="'+msg.avator+'"/>'+msg.username+'</span>',
                    position: 'top-left',
                });
            }
            if (redata.type == "notice") {
                let msg = redata.data
                if(!msg){
                    return;
                }
                this.chatTitle=msg
                this.showTitle(msg);
                this.scrollBottom();
            }
            if (redata.type == "close") {
                this.chatTitle=GOFLY_LANG[LANG]['closemes'];
                this.showTitle(this.chatTitle);
                this.scrollBottom();
                this.socket.close();
                this.socketClosed=true;
            }
            if (redata.type == "force_close") {
                this.chatTitle=GOFLY_LANG[LANG]['forceclosemes'];
                this.showTitle(this.chatTitle);
                this.scrollBottom();
                this.socket.close();
                this.socketClosed=true;
            }
            if (redata.type == "auto_close") {
                this.chatTitle=GOFLY_LANG[LANG]['autoclosemes'];
                this.showTitle(this.chatTitle);
                this.scrollBottom();
                this.socket.close();
            }
            if (redata.type == "message") {
                let msg = redata.data
                if(msg.id==this.visitor.visitor_id){
                    return;
                }
                let content = {}
                content.avator = msg.avator;
                content.name = msg.name;
                content.content =replaceContent(msg.content);
                content.is_kefu = false;
                content.time = msg.time;
                this.msgList.push(content);
                this.scrollBottom();
            }
            window.parent.postMessage(redata,"*");
        },
        OnClose:function() {
            console.log("ws:onclose");
            this.socketClosed=true;
        },
        OnError:function() {
            console.log("ws:onerror");
            this.socketClosed=true;
        },
        //初始化用户信息
        initUserInfo:function(){
            var entId=getQuery("ent_id");
            var roomId=getQuery("room_id");
            if(entId==""||roomId==""){
                this.$message({
                    message: "error:ent_id,room_id",
                    type: 'error'
                });
                return;
            }
            var visitorId="";
            var toId=roomId;
            var obj=getLocalStorage("visitor_"+entId);

            if(obj){
                visitorId=obj.visitor_id;
                toId=obj.to_id;
            }
            var params={visitor_id:visitorId,to_id:toId,ent_id:entId};
            var  _this=this;
            this.sendAjax("/room/login","post",params,function(data){
                setLocalStorage("visitor_"+entId,data);
                _this.visitor=data;
                _this.initConn();
                _this.ping();
            })
        },
        //软键盘问题
        textareaFocus:function(){
            this.scrollBottom()
            if(/Android|webOS|iPhone|iPad|BlackBerry/i.test(navigator.userAgent)) {
                $(".chatContext").css("margin-bottom","0");
                $(".chatBoxSend").css("position","static");
                this.textareaFocused=true;
            }
        },
        textareaBlur:function(){
            if(this.textareaFocused&&/Android|webOS|iPhone|iPad|BlackBerry/i.test(navigator.userAgent)) {
                var chatBoxSendObj=$(".chatBoxSend");
                var chatContextObj=$(".chatContext");
                if(this.textareaFocused&&chatBoxSendObj.css("position")!="fixed"){
                    chatContextObj.css("margin-bottom","105px");
                    chatBoxSendObj.css("position","fixed");
                    this.textareaFocused=false;
                }

            }
        },
        //滚动到底部
        scrollBottom:function(){
            var _this=this;
            this.$nextTick(function(){
                $('.chatVisitorPage').scrollTop($(".chatVisitorPage")[0].scrollHeight);
            });
        },
        chatToUser:function() {
            if(this.socketClosed){
                this.initConn();
                // this.$message({
                //     message: GOFLY_LANG[LANG]['socketclose'],
                //     type: 'warning'
                // });
                // return;
            }
            var messageContent=this.messageContent.trim("\r\n");
            messageContent=messageContent.replace("\n","");
            messageContent=messageContent.replace("\r\n","");
            if(messageContent==""||messageContent=="\r\n"){
                this.messageContent="";
                return;
            }
            this.messageContent=messageContent;
            this.sendDisabled=true;
            let _this=this;
            let content = {}
            content.avator=_this.visitor.avator;
            content.content = replaceContent(_this.messageContent);
            content.name = _this.visitor.name;
            content.is_kefu = true;
            content.time = "";
            content.show_time=false;
            _this.msgList.push(content);
            _this.scrollBottom();

            let mes = {};
            mes.type = "visitor";
            mes.content = this.messageContent;
            mes.from_id = this.visitor.visitor_id;
            mes.to_id = this.visitor.to_id;
            mes.content = this.messageContent;
            $.post("/room/message",mes,function(res){
                _this.sendDisabled=false;
                if(res.code!=200){
                    _this.msgList.pop();
                    _this.$message({
                        message: res.msg,
                        type: 'error'
                    });
                    return;
                }
                _this.messageContent = "";
                _this.sendDisabled=false;
            });
        },
        //发送ajax
        sendAjax:function(url,method,params,callback){
            let _this=this;
            $.ajax({
                type: method,
                url: url,
                data:params,
                error:function(res){
                    var data=JSON.parse(res.responseText);
                    console.log(data);
                    if(data.code!=200){
                        _this.$message({
                            message: data.msg,
                            type: 'error'
                        });
                    }
                },
                success: function(data) {
                    if(data.code!=200){
                        _this.$message({
                            message: data.msg,
                            type: 'error'
                        });
                    }else if(data.result!=null){
                        callback(data.result);
                    }else{
                        callback(data);
                    }
                }
            });
        },
        showTitle:function(title){
            $(".chatBox").append("<div class=\"chatTime\"><span>"+title+"</span></div>");
            this.scrollBottom();
        },
        //表情点击事件
        faceIconClick:function(index){
            $('.faceBox').hide();
            this.messageContent+="face"+this.face[index].name;
        },
        initCss:function(){
            var _this=this;
            $(function () {
                //展示表情
                var faces=placeFace();
                $.each(faceTitles, function (index, item) {
                    _this.face.push({"name":item,"path":faces[item]});
                });
                $(".visitorFaceBtn").click(function(e){
                    var status=$('.faceBox').css("display");
                    if(status=="block"){
                        $('.faceBox').hide();
                    }else{
                        $('.faceBox').show();
                    }
                    return false;
                });
                $('body').click(function(){
                    $('.faceBox').hide();
                });
                _this.heartsFlow = new HeartsFlow({
                    canvasEl: '#hearts-canvas',
                    amount: 20 });
            });

        },
        //心跳
        ping:function(){
            let _this=this;
            let mes = {}
            mes.type = "ping";
            mes.data = "visitor:"+_this.visitor.visitor_id;
            setInterval(function () {
                if(_this.socket!=null&&!this.socketClosed){
                    _this.socket.send(JSON.stringify(mes));
                }
            },55000);
            setInterval(function () {
                _this.likeHeart();
            },3000);
        },
        likeHeart:function(){
            this.heartsFlow.startAnimation();
        },
    },
    mounted:function() {

    },
    created: function () {
        this.initCss();
        this.initUserInfo();
    }
})
