let socket = null;

const statusColor = {
    UNKNOWN: 'gray',
    NOT_INUSE: 'green',
    INUSE: 'blue',
    BUSY: 'red',
    INVALID: 'purple',
    UNAVAILABLE: 'orange',
    RINGING: 'yellow',
    RINGINUSE: 'orange',
    ONHOLD: 'lightblue',
};

document.addEventListener("DOMContentLoaded", function () {
    socket = new WebSocket("ws://localhost:9999/");
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
                var state = document.getElementById(`devstate_${jmsg.data[0]}`)
                state.innerHTML = jmsg.data[1];
                var classes = Array.from(state.classList);
                classes.pop();
                classes.push(statusColor[jmsg.data[1]])
                state.className = classes.join(' ');
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
    const devlist = document.getElementById('devlist')

    document.getElementById('devs').innerHTML = data.activedev
    document.getElementById('chans').innerHTML = data.bridgecount

    for (let index = 0; index < data.devicelist.length; index++) {
        const element = data.devicelist[index];

        instance = document.importNode(template.content, true);

        var devName = instance.querySelector('.ui.left.aligned.segment');
        var devStat = instance.querySelector('.ui.inverted.secondary.segment');
        devName.id = 'devname_' + element.name;
        devStat.id = 'devstate_' + element.name;

        devName.innerHTML = element.name;
        devStat.innerHTML = element.status;

        devStat.classList.add(statusColor[element.status])

        devlist.appendChild(instance);
    }
}

function addEvent(message, icon) {
    const template = document.getElementById('event');
    const feed = document.getElementById('feed');

    const date = new Date();
    const instance = document.importNode(template.content, true);
    instance.querySelector('.date').innerHTML = date.getDate() + '-' + (date.getMonth() + 1) + '-' + date.getFullYear() + " " + date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds();
    instance.querySelector('.summary').innerHTML = message;

    instance.querySelector('.label').innerHTML = `<i class="${icon} tiny icon inverted"></i>`;

    feed.insertBefore(instance, feed.firstChild)
}