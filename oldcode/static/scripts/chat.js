
function getCookie(name) {
    const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
    if (match) {
        return match[2];
    }
    return null;
}

let newMessageBannerVisible = false;
const userid = getCookie('userid');
const token = getCookie('token');

var button = document.getElementById("inputRename");
button.onkeydown = (event) => {
    if (event.key == "Enter") {
        if (input.value.length == 0) {
            return
        }
        var jsonMessage = JSON.stringify({
            token: token,
            userid: userid,
            channelid: 1,
            message: input.value,
        });
        socket.send(jsonMessage);
        input.value = "";
    }
}


if (!userid || !token) {
    window.location.href = "/login";
}

var input = document.getElementById("inputRename");
var socket = new WebSocket("ws://localhost:8080/server");

window.onbeforeunload = function() {
    if (socket) {
        socket.close();
    }
};
socket.onerror = function(error) {
    console.log("WebSocket Error: ", error);
};
socket.onclose = function(event) {
    if (event.wasClean) {
        console.log(`Closed cleanly, code=${event.code}, reason=${event.reason}`);
    } else {
        console.error(`Connection closed unexpectedly: ${event.code}`);
    }
};

document.addEventListener("DOMContentLoaded", () => {
    function hideBannner() {
        const banner = document.getElementById("newMessageBanner");
        if (banner) {
            banner.style.display = "none";
            newMessageBannerVisible = false;
        }
    }

    function showBannner() {
        const banner = document.getElementById("newMessageBanner");
        if (banner) {
            banner.style.display = "block";
            newMessageBannerVisible = true;
        }
    }
    const output = document.getElementById("output");
    var lastScrollTop = 0;
    var scrolledToBottom = null;
    const atBottom = () => {
        if (output.scrollTop < lastScrollTop) {
            return false;
        }
        lastScrollTop = output.scrollTop + output.offsetHeight <= 0 ? 0 : output.scrollTop;
        var val = output.scrollTop + output.offsetHeight
        if (val >= output.scrollHeight) {
            return true;
        }
        return false;
    }
    output.onscroll = (e) => {
        scrolledToBottom = atBottom();
        if (scrolledToBottom) {
            hideBannner();
        }
    };
    socket.onopen = function() {

        const statusDiv = document.createElement("div");
        statusDiv.className = "status";
        statusDiv.textContent = "Status: Connected";
        output.appendChild(statusDiv);
    };
    if (!output) {
        console.error("Output element nont found");
        return
    }
    socket.onmessage = function(e) {
        const startedOnBottom = atBottom();
        let data = JSON.parse(e.data);
        const messageContainer = document.createElement("div");
        messageContainer.className = "message-container";
        const usernameDiv = document.createElement("div");
        usernameDiv.className = "username";
        usernameDiv.textContent = data.username;
        const messageDiv = document.createElement("div");
        messageDiv.className = "message";
        messageDiv.textContent = data.message;

        const messageTimestampDiv = document.createElement("div");
        messageTimestampDiv.className = "message-timestamp";
        messageTimestampDiv.textContent = data.date;
        messageContainer.appendChild(messageTimestampDiv);
        messageContainer.appendChild(usernameDiv);
        messageContainer.appendChild(messageDiv);
        output.appendChild(messageContainer);
        newestMessageDiv = messageContainer;
        if (startedOnBottom || data.userid == userid) {
            newestMessageDiv.scrollIntoView({ behavior: "smooth" });
            scrolledToBottom = true;
        } else if (!atBottom()) {
            showBannner();
        }
    };
});
