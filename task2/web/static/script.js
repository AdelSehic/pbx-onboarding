
function addEvent(message, icon){
    const template = document.getElementById('event');
    const feed = document.getElementById('feed');
    
    const date = new Date();
    const instance = document.importNode(template.content, true);
    instance.querySelector('.date').innerHTML = date.getDate()+'-'+(date.getMonth()+1)+'-'+date.getFullYear()+" "+date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds();
    instance.querySelector('.summary').innerHTML = message;

    instance.querySelector('.label').innerHTML = `<i class="${icon} tiny icon"></i>`;

    feed.insertBefore(instance, feed.firstChild)
}