let socket
let subscriptions = {}

// Opens a new websocket connection to backend
// for managing gas price subscriptions
const createWebsocketConnection = _ => {

    const connAlreadyOpen = '[ `gasz` ] Connection Already Open'
    const connOpen = '[ `gasz` ] Connection Opened'
    const connClosed = '[ `gasz` ] Connection Closed'
    const connError = '[ `gasz` ] Error in connection'

    return new Promise((res, rej) => {

        if (socket && socket.readyState === socket.OPEN) {
            return res(connAlreadyOpen)
        }

        socket = new WebSocket(`ws://localhost:7000/v1/subscribe`)

        // websocket connection is open now
        socket.onopen = _ => {
            return res(connOpen)
        }

        // connection with backend got closed
        socket.onclose = _ => {
            return rej(connClosed)
        }

        // due to some error encountered, closing connection with backend
        socket.onerror = _ => {
            socket.close()

            return rej(connError)
        }

        // Handling case when message being received from server
        socket.onmessage = e => {
            // data received from server
            const msg = JSON.parse(e.data)

            // -- Starting to handle subscription/ unsubsciption messages
            if ('code' in msg) {
                self.clients.matchAll({ includeUncontrolled: true }).then(clients => {
                    clients.forEach(client => client.postMessage(JSON.stringify(msg)))
                })

                return
            }
            // -- upto this point

            this.registration.showNotification('Gasz ⚡️', {
                body: `Gas Price for ${msg['txType'].slice(0, 1).toUpperCase() + msg['txType'].slice(1)} transaction just reached ${msg['price']} Gwei`,
                icon: 'gasz.png',
                badge: 'gasz.png',
                tag: msg['topic'],
                requireInteraction: true,
                vibrate: [200, 100, 200]
            })
        }

    })

}

this.addEventListener('activate', _ => {

    // Checking whether already connected via websocket or not
    //
    // if not, new connection will be created
    createWebsocketConnection()
        .then(console.log)
        .catch(console.error)

})

this.addEventListener('message', m => {
    createWebsocketConnection()
        .then(v => {
            console.log(v)

            // Keeping track of which topic this client is subscribed to
            subscriptions[`${m.data['field']} : ${m.data['operator']} ${m.data['threshold']}`] = JSON.parse(m.data)

            socket.send(m.data)
        })
        .catch(console.error)
})

this.addEventListener('notificationclick', e => {

    console.log(e.data)

    if (subscriptions.length > 0) {
        socket.send(JSON.stringify({ ...subscriptions[0], type: 'unsubscription' }))
    }

    e.notification.close()
})
