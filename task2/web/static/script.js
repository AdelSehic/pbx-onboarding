let socket = null;

document.addEventListener("DOMContentLoaded", function () {
    socket = new WebSocket("ws://10.1.0.109:9999/");
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
        switch (jmsg.type) {
            case "setup":
                setup(jmsg.data)
                break;
            case "devcountupdate":
                document.getElementById("devs").innerHTML = jmsg.data
                break;
            case "brcountupdate":
                document.getElementById("chans").innerHTML = jmsg.data[1]
                addEvent(jmsg.data[0], "exchange")
                break;
            case "activedevs":
                document.getElementById("devs").innerHTML = jmsg.data
                break;
            case "devstatechange":
                const elem = document.getElementById(jmsg.data[0]);
                elem.innerHTML = `<a>${jmsg.data[0]}</a> ${jmsg.data[1]}`;
                addEvent(`<a> ${jmsg.data[0]} </a> is now <a> ${jmsg.data[1]} </a>`, "phone")
                break;
            case "succauth":
                addEvent(`Successful authentication by <a>${jmsg.data[0]}</a> from <a>${jmsg.data[1]}</a>`, "key")
                break;
            default:
                console.log("ERROR: unrecognized request")
        }
    }
})

function setup(data) {
    const template = document.getElementById('device')
    const devlist = document.getElementById("devlist")

    document.getElementById("devs").innerHTML = data.activedev
    document.getElementById("chans").innerHTML = data.bridgecount

    for (let index = 0; index < data.devicelist.length; index++) {
        const element = data.devicelist[index];
        let instance = document.importNode(template.content, true);
        const newItem = document.createElement("div");
        newItem.className = "ui segment";
        newItem.id = element.name
        newItem.innerHTML = `<a>${element.name}</a> ${element.status}`;
        devlist.appendChild(newItem);
    }
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