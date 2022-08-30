/*
 * Live2D Widget
 * https://github.com/stevenjoezhang/live2d-widget
 */

function loadWidget(config) {
    let { waifuPath, apiPath, cdnPath } = config;
    let useCDN = false, modelList;
    if (typeof cdnPath === "string") {
        useCDN = true;
        if (!cdnPath.endsWith("/")) cdnPath += "/";
    } else if (typeof apiPath === "string") {
        if (!apiPath.endsWith("/")) apiPath += "/";
    } else {
        console.error("Invalid initWidget argument!");
        return;
    }
    localStorage.removeItem("waifu-display");
    sessionStorage.removeItem("waifu-text");
    document.body.insertAdjacentHTML("beforeend", `<div id="waifu" class="waifu">
			<div id="waifu-tips" class="waifu-tips"></div>
			<input id="sendMsg" placeholder="在这输入信息哦!" type="text" class="waifu-input"/>
			<canvas id="live2d" width="800" height="800"></canvas>
			<div id="waifu-tool" class="waifu-tool">
				<span class="fa fa-lg fa-comment"></span>
				<span class="fa fa-lg fa-paper-plane"></span>
				<span class="fa fa-lg fa-user-circle"></span>
				<span class="fa fa-lg fa-street-view"></span>
				<span class="fa fa-lg fa-camera-retro"></span>
				<span class="fa fa-lg fa-info-circle"></span>
				<span class="fa fa-lg fa-times"></span>
			</div>
		</div>`);
    // https://stackoverflow.com/questions/24148403/trigger-css-transition-on-appended-element
    setTimeout(() => {
        document.getElementById("waifu").style.bottom = 0;
    }, 0);

    function randomSelection(obj) {
        return Array.isArray(obj) ? obj[Math.floor(Math.random() * obj.length)] : obj;
    }
    // 检测用户活动状态，并在空闲时显示消息
    let messageTimer;

    (function registerEventListener() {
        document.querySelector("#waifu-tool .fa-comment").addEventListener("click", showHitokoto);
        document.querySelector("#waifu-tool .fa-paper-plane").addEventListener("click", () => {
            if (window.Asteroids) {
                if (!window.ASTEROIDSPLAYERS) window.ASTEROIDSPLAYERS = [];
                window.ASTEROIDSPLAYERS.push(new Asteroids());
            } else {
                const script = document.createElement("script");
                script.src = "https://cdn.jsdelivr.net/gh/stevenjoezhang/asteroids/asteroids.js";
                document.head.appendChild(script);
            }
        });
        document.querySelector("#waifu-tool .fa-user-circle").addEventListener("click", loadOtherModel);
        document.querySelector("#waifu-tool .fa-street-view").addEventListener("click", loadRandModel);
        document.querySelector("#waifu-tool .fa-camera-retro").addEventListener("click", () => {
            showMessage("照好了嘛，是不是很可爱呢？", 6000, 9);
            Live2D.captureName = "photo.png";
            Live2D.captureFrame = true;
        });
        document.querySelector("#waifu-tool .fa-info-circle").addEventListener("click", () => {
            open("https://github.com/taoshihan1991/go-fly");
        });
        document.querySelector("#waifu-tool .fa-times").addEventListener("click", () => {
            localStorage.setItem("waifu-display", Date.now());
            document.getElementById("waifu").style.bottom = "-500px";
            setTimeout(() => {
                document.getElementById("waifu").style.display = "none";
                document.getElementById("waifu-toggle").classList.add("waifu-toggle-active");
            }, 3000);
        });
        document.querySelector("#waifu-tips,#sendMsg").addEventListener("click", function(e){
            var obj=document.getElementById('layui-layer-iframe19911116');
            if(!obj){
                GOFLY.showPanel();
            }
        });
        document.querySelector("#sendMsg").addEventListener("keypress", function(e){
            var obj=document.getElementById('layui-layer-iframe19911116');
            if(e.keyCode != "13"){
                return;
            }
            if(!obj){
                GOFLY.showPanel();
            }
            var msg={}
            msg.type='send_message';
            msg.content=document.getElementById("sendMsg").value;
            document.getElementById("sendMsg").value="";
            obj.contentWindow.postMessage(msg, "*");
        });

    })();

    function showHitokoto() {
    }

    function showMessage(text, timeout, priority) {
        if (!text || (sessionStorage.getItem("waifu-text") && sessionStorage.getItem("waifu-text") > priority)) return;
        if (messageTimer) {
            clearTimeout(messageTimer);
            messageTimer = null;
        }
        text = randomSelection(text);
        sessionStorage.setItem("waifu-text", priority);
        const tips = document.getElementById("waifu-tips");
        tips.innerHTML = text;
        tips.classList.add("waifu-tips-active");
        messageTimer = setTimeout(() => {
            sessionStorage.removeItem("waifu-text");
            tips.classList.remove("waifu-tips-active");
        }, timeout);
    }

    (function initModel() {
        let modelId = localStorage.getItem("modelId"),
            modelTexturesId = localStorage.getItem("modelTexturesId");
        if (modelId === null) {
            // 首次访问加载 指定模型 的 指定材质
            modelId = 2; // 模型 ID
            modelTexturesId = 53; // 材质 ID
        }
        loadModel(modelId, modelTexturesId);
    })();

    async function loadModelList() {
        const response = await fetch(`${cdnPath}model_list.json`);
        modelList = await response.json();
    }

    async function loadModel(modelId, modelTexturesId, message) {
        localStorage.setItem("modelId", modelId);
        localStorage.setItem("modelTexturesId", modelTexturesId);
        showMessage(message, 4000, 10);
        if (useCDN) {
            if (!modelList) await loadModelList();
            const target = randomSelection(modelList.models[modelId]);
            loadlive2d("live2d", `${cdnPath}model/${target}/index.json`);
        } else {
            loadlive2d("live2d", `${apiPath}get/?id=${modelId}-${modelTexturesId}`);
            console.log(`Live2D 模型 ${modelId}-${modelTexturesId} 加载完成`);
        }
    }

    async function loadRandModel() {
        const modelId = localStorage.getItem("modelId"),
            modelTexturesId = localStorage.getItem("modelTexturesId");
        if (useCDN) {
            if (!modelList) await loadModelList();
            const target = randomSelection(modelList.models[modelId]);
            loadlive2d("live2d", `${cdnPath}model/${target}/index.json`);
            showMessage("我的新衣服好看嘛？", 4000, 10);
        } else {
            // 可选 "rand"(随机), "switch"(顺序)
            fetch(`${apiPath}rand_textures/?id=${modelId}-${modelTexturesId}`)
                .then(response => response.json())
                .then(result => {
                    if (result.textures.id === 1 && (modelTexturesId === 1 || modelTexturesId === 0)) showMessage("我还没有其他衣服呢！", 4000, 10);
                    else loadModel(modelId, result.textures.id, "我的新衣服好看嘛？");
                });
        }
    }

    async function loadOtherModel() {
        let modelId = localStorage.getItem("modelId");
        if (useCDN) {
            if (!modelList) await loadModelList();
            const index = (++modelId >= modelList.models.length) ? 0 : modelId;
            loadModel(index, 0, modelList.messages[index]);
        } else {
            fetch(`${apiPath}switch/?id=${modelId}`)
                .then(response => response.json())
                .then(result => {
                    loadModel(result.model.id, 0, result.model.message);
                });
        }
    }
}

function initWidget(config, apiPath) {
    if (typeof config === "string") {
        config = {
            waifuPath: config,
            apiPath
        };
    }
    document.body.insertAdjacentHTML("beforeend", `<div class="waifu-toggle" id="waifu-toggle">
			<span>看板娘</span>
		</div>`);
    const toggle = document.getElementById("waifu-toggle");
    toggle.addEventListener("click", () => {
        toggle.classList.remove("waifu-toggle-active");
        if (toggle.getAttribute("first-time")) {
            loadWidget(config);
            toggle.removeAttribute("first-time");
        } else {
            localStorage.removeItem("waifu-display");
            document.getElementById("waifu").style.display = "";
            setTimeout(() => {
                document.getElementById("waifu").style.bottom = 0;
            }, 0);
        }
    });
    if (localStorage.getItem("waifu-display") && Date.now() - localStorage.getItem("waifu-display") <= 86400000) {
        toggle.setAttribute("first-time", true);
        setTimeout(() => {
            toggle.classList.add("waifu-toggle-active");
        }, 0);
    } else {
        loadWidget(config);
    }
}
