let socket

// Opens a new websocket connection to backend
// for managing gas price subscriptions
const createWebsocketConnection = _ => {

    if(socket && socket.readyState === socket.OPEN) {
        return socket
    }

    socket = new WebSocket(`ws://localhost:7000/v1/subscribe`)

    socket.onopen = _ => {
        console.log('[ `gasz` ] Connection Opened')
    }
    
    socket.onclose = _ => {
        console.log('[ `gasz` ] Connection Closed')
    }
    
    socket.onerror = _ => {
        console.log('[ `gasz` ] Error in connection')
        socket.close()
    }

    return socket

}

this.addEventListener('activate', _ => {

    

    // Handling case when message being received from server
    socket.onmessage = e => {
        // data received from server
        const msg = JSON.parse(e.data)
        
        // -- Staring to handle subscription/ unsubsciption messages
        if ('code' in msg){
            if(msg['code'] !== 1) {
                if (msg['message'] === 'Already Subscribed') { } else { }
            } else {
                if (msg['message'].includes('Subscribed')) { } else { }
            }

            return
        }
        // -- upto this point

        if (Notification.permission === 'granted') {
            const notify = new Notification('Gasz ⚡️', {body: `Body`, icon: 'gasz.png'})

            notify.onclick = _ => {
                notify.close()
            }
        }
    }

})

this.addEventListener('message', m => {
    console.log(`Received from client : '${m.data}'`)

    this.registration.showNotification('Gasz ⚡️', {body: `${m.data}`, icon: 'gasz.png'})

    m.source.postMessage(m.data)
})
