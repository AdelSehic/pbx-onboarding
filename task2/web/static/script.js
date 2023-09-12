let socket = null;
document.addEventListener("DOMContentLoaded", function() {
    socket = new WebSocket("ws://127.0.0.1:9999/");
    socket.onopen = (data) => {
        console.log("Websocket connection established")
    }
    socket.onclose = () => {
        console.log("Websocket connection closed")
    }
    socket.onerror = error => {
        console.log(error)
    }
    socket.onmessage = msg => {
        let jmsg = JSON.parse(msg.data)
        switch (jmsg.type){
            case "setup":
                setup(jmsg.data)
                break;
            default:
                console.log("ERROR: unrecognized request")
        }
    }
})

function setup(data){
    document.getElementById("devs").innerHTML = data.devicecount
    document.getElementById("chans").innerHTML = data.bridgecount
}

function fetchData(url) {
    fetch(url, {
        headers: {
            'Accept': 'application/json'
        }
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            const devCount = data.DeviceCount;
            const brCount = data.BridgeCount;
            // const devices = data.DeviceList;

            document.getElementById("devs").innerHTML = devCount
            document.getElementById("chans").innerHTML = brCount
        })
        .catch(error => {
            console.log('Error: ', error)
        })
}

function addEvent(message, icon) {
    const template = document.getElementById('event');
    const feed = document.getElementById('feed');

    const date = new Date();
    const instance = document.importNode(template.content, true);
    instance.querySelector('.date').innerHTML = date.getDate() + '-' + (date.getMonth() + 1) + '-' + date.getFullYear() + " " + date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds();
    instance.querySelector('.summary').innerHTML = message;

    instance.querySelector('.label').innerHTML = `<i class="${icon} tiny icon"></i>`;

    feed.insertBefore(instance, feed.firstChild)
}