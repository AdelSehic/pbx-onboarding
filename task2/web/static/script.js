window.onload = fetchData("http://127.0.0.1:9999/init")

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