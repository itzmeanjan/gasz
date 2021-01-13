let socket

// Opens a new websocket connection to backend
// for managing gas price subscriptions
const createWebsocketConnection = _ => {

    if (socket && socket.readyState === socket.OPEN) {
        return
    }

    socket = new WebSocket(`ws://localhost:7000/v1/subscribe`)

    // websocket connection is open now
    socket.onopen = _ => {
        console.log('[ `gasz` ] Connection Opened')
    }

    // connection with backend got closed
    socket.onclose = _ => {
        console.log('[ `gasz` ] Connection Closed')
    }

    // due to some error encountered, closing connection with backend
    socket.onerror = _ => {
        console.log('[ `gasz` ] Error in connection')
        socket.close()
    }

    // Handling case when message being received from server
    socket.onmessage = e => {
        // data received from server
        const msg = JSON.parse(e.data)

        // -- Staring to handle subscription/ unsubsciption messages
        if ('code' in msg) {
            if (msg['code'] !== 1) {
                if (msg['message'] === 'Already Subscribed') {

                    self.clients.matchAll().then(clients => {
                        clients.forEach(client => client.postMessage({ msg: 'Hello from SW' }));
                    })

                } else {
                    self.clients.matchAll().then(clients => {
                        clients.forEach(client => client.postMessage({ msg: 'Hello from SW' }));
                    })
                }
            } else {
                if (msg['message'].includes('Subscribed')) {
                    self.clients.matchAll().then(clients => {
                        clients.forEach(client => client.postMessage({ msg: 'Hello from SW' }));
                    })
                } else {
                    self.clients.matchAll().then(clients => {
                        clients.forEach(client => client.postMessage({ msg: 'Hello from SW' }));
                    })
                }
            }

            return
        }
        // -- upto this point

        this.registration.showNotification('Gasz ⚡️', { body: `${m.data}`, icon: 'gasz.png' })
    }

}

this.addEventListener('activate', _ => {

    // Checking whether already connected via websocket or not
    //
    // if not, new connection will be created
    createWebsocketConnection()

})

this.addEventListener('message', m => {
    createWebsocketConnection()

    socket.send(JSON.stringify(m.data))
})
